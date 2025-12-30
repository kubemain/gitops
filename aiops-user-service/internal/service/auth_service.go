// internal/service/auth_service.go
package service

import (
	"aiops-user-service/internal/repository"
	"aiops-user-service/pkg/database"
	"aiops-user-service/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type AuthService struct {
	userRepo *repository.UserRepository
	permRepo *repository.PermissionRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
		permRepo: repository.NewPermissionRepository(),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token       string    `json:"token"`
	User        *UserInfo `json:"user"`
	Permissions []string  `json:"permissions"`
}

type UserInfo struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest, ip string) (*LoginResponse, error) {
	// 1. 查询用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 2. 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 3. 验证密码
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 4. 提取角色ID
	roleIDs := make([]uint, len(user.Roles))
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleIDs[i] = role.ID
		roleNames[i] = role.Name
	}

	// 5. 生成 Token
	token, err := utils.GenerateToken(user.ID, user.Username, roleIDs)
	if err != nil {
		return nil, err
	}

	// 6. 查询用户权限
	permissions, err := s.permRepo.FindByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	permCodes := make([]string, len(permissions))
	for i, perm := range permissions {
		permCodes[i] = perm.Code
	}

	// 7. 缓存用户权限到 Redis（30分钟）
	cacheKey := fmt.Sprintf("user:perms:%d", user.ID)
	permJSON, _ := json.Marshal(permCodes)
	ctx := context.Background()
	database.RDB.Set(ctx, cacheKey, permJSON, 30*time.Minute)

	// 8. 更新最后登录信息
	go s.userRepo.UpdateLastLogin(user.ID, ip)

	// 9. 返回响应
	return &LoginResponse{
		Token: token,
		User: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Avatar:   user.Avatar,
			Roles:    roleNames,
		},
		Permissions: permCodes,
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(userID uint, token string) error {
	ctx := context.Background()

	// 删除权限缓存
	cacheKey := fmt.Sprintf("user:perms:%d", userID)
	database.RDB.Del(ctx, cacheKey)

	// 将 Token 加入黑名单（过期时间与 Token 一致）
	blacklistKey := fmt.Sprintf("token:blacklist:%s", token)
	database.RDB.Set(ctx, blacklistKey, "1", 24*time.Hour)

	return nil
}

// GetUserPermissions 获取用户权限（带缓存）
func (s *AuthService) GetUserPermissions(userID uint) ([]string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:perms:%d", userID)

	// 先从缓存读取
	val, err := database.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var perms []string
		if err := json.Unmarshal([]byte(val), &perms); err == nil {
			return perms, nil
		}
	}

	// 缓存未命中，查询数据库
	permissions, err := s.permRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	permCodes := make([]string, len(permissions))
	for i, perm := range permissions {
		permCodes[i] = perm.Code
	}

	// 写入缓存
	permJSON, _ := json.Marshal(permCodes)
	database.RDB.Set(ctx, cacheKey, permJSON, 30*time.Minute)

	return permCodes, nil
}

// CheckPermission 检查用户是否有某权限
func (s *AuthService) CheckPermission(userID uint, permCode string) (bool, error) {
	perms, err := s.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p == permCode || p == "*:*" { // *:* 表示超级管理员
			return true, nil
		}
	}

	return false, nil
}
