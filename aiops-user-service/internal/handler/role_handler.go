// internal/handler/role_handler.go
package handler

import (
	"aiops-user-service/internal/service"
	"aiops-user-service/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: service.NewRoleService(),
	}
}

// List 角色列表
// @Summary 角色列表
// @Tags 角色管理
// @Success 200 {object} utils.Response
// @Router /api/v1/roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roleService.GetAll()
	if err != nil {
		utils.InternalServerError(c, "查询失败")
		return
	}

	utils.Success(c, roles)
}

// GetByID 查询角色详情
// @Summary 查询角色详情
// @Tags 角色管理
// @Param id path int true "角色ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	role, err := h.roleService.GetByID(uint(id))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}

	utils.Success(c, role)
}

// AssignPermissions 分配权限
// @Summary 分配权限
// @Tags 角色管理
// @Param id path int true "角色ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID格式错误")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.roleService.AssignPermissions(uint(id), req.PermissionIDs); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "权限分配成功", nil)
}

// GetAllPermissions 获取所有权限
// @Summary 获取所有权限
// @Tags 角色管理
// @Success 200 {object} utils.Response
// @Router /api/v1/permissions [get]
func (h *RoleHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.roleService.GetAllPermissions()
	if err != nil {
		utils.InternalServerError(c, "查询失败")
		return
	}

	utils.Success(c, permissions)
}
