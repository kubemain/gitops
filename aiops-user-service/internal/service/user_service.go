// internal/service/user_service.go
package service

import (
	"aiops-user-service/internal/domain/entity"
	"aiops-user-service/internal/repository"
	"aiops-user-service/pkg/utils"
	"errors"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6"`
	RoleIDs  []uint `json:"role_ids"`
	Status   int8   `json:"status"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	RoleIDs  []uint `json:"role_ids"`
	Status   int8   `json:"status"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Username string `form:"username"`
	Status   int    `form:"status"`
}

// Create 创建用户
func (s *UserService) Create(req *CreateUserRequest) (*entity.User, error) {
	// 1. 检查用户名是否存在
	existUser, _ := s.userRepo.FindByUsername(req.Username)
	if existUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 2. 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. 创建用户
	user := &entity.User{
		Username:     req.Username,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Status:       req.Status,
	}

	if user.Status == 0 {
		user.Status = 1 // 默认启用
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 4. 分配角色
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.AssignRoles(user.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	// 5. 重新查询（包含角色信息）
	return s.userRepo.FindByID(user.ID)
}

// GetByID 根据ID查询用户
func (s *UserService) GetByID(id uint) (*entity.User, error) {
	return s.userRepo.FindByID(id)
}

// Update 更新用户
func (s *UserService) Update(id uint, req *UpdateUserRequest) (*entity.User, error) {
	// 1. 查询用户
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 2. 更新字段
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Avatar = req.Avatar
	user.Status = req.Status

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// 3. 更新角色
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.AssignRoles(id, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	// 4. 重新查询
	return s.userRepo.FindByID(id)
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	// 不能删除ID为1的超级管理员
	if id == 1 {
		return errors.New("不能删除超级管理员")
	}
	return s.userRepo.Delete(id)
}

// List 用户列表
func (s *UserService) List(req *UserListRequest) ([]*entity.User, int64, error) {
	filters := make(map[string]interface{})
	if req.Username != "" {
		filters["username"] = req.Username
	}
	if req.Status > 0 {
		filters["status"] = req.Status
	}

	return s.userRepo.List(req.Page, req.PageSize, filters)
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(id uint, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}
