// internal/handler/user_handler.go
package handler

import (
	"aiops-user-service/internal/service"
	"aiops-user-service/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Create 创建用户
// @Summary 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.CreateUserRequest true "用户信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Create(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", user)
}

// GetByID 查询用户详情
// @Summary 查询用户详情
// @Tags 用户管理
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}

	utils.Success(c, user)
}

// Update 更新用户
// @Summary 更新用户
// @Tags 用户管理
// @Param id path int true "用户ID"
// @Param request body service.UpdateUserRequest true "用户信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Update(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", user)
}

// Delete 删除用户
// @Summary 删除用户
// @Tags 用户管理
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// List 用户列表
// @Summary 用户列表
// @Tags 用户管理
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param username query string false "用户名"
// @Success 200 {object} utils.Response
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	var req service.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	users, total, err := h.userService.List(&req)
	if err != nil {
		utils.InternalServerError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, users, total, req.Page, req.PageSize)
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Tags 用户管理
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/users/{id}/password [put]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.ResetPassword(uint(id), req.NewPassword); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "密码重置成功", nil)
}
