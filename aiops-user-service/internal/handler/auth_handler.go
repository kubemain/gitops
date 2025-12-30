// internal/handler/auth_handler.go
package handler

import (
	"aiops-user-service/internal/service"
	"aiops-user-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取客户端IP
	ip := c.ClientIP()

	result, err := h.authService.Login(&req, ip)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	utils.Success(c, result)
}

// Logout 用户登出
// @Summary 用户登出
// @Tags 认证
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get("user_id")
	token, _ := c.Get("token")

	if err := h.authService.Logout(userID.(uint), token.(string)); err != nil {
		utils.InternalServerError(c, "登出失败")
		return
	}

	utils.SuccessWithMessage(c, "登出成功", nil)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/userinfo [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	// 获取用户权限
	permissions, err := h.authService.GetUserPermissions(userID.(uint))
	if err != nil {
		utils.InternalServerError(c, "获取用户信息失败")
		return
	}

	utils.Success(c, gin.H{
		"user_id":     userID,
		"username":    username,
		"permissions": permissions,
	})
}
