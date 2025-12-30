package handler

import (
	"aiops-cmdb/internal/model"
	"aiops-cmdb/internal/service"
	"aiops-cmdb/pkg/logger"
	"aiops-cmdb/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HostGroupHandler struct {
	service *service.HostGroupService
}

func NewHostGroupHandler() *HostGroupHandler {
	return &HostGroupHandler{
		service: service.NewHostGroupService(),
	}
}

func (h *HostGroupHandler) Create(c *gin.Context) {
	var group model.HostGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Create(&group); err != nil {
		logger.Error("创建分组失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	logger.Info("创建分组成功", zap.String("name", group.Name))
	response.SuccessWithMessage(c, "创建成功", group)
}

func (h *HostGroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var group model.HostGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Update(uint(id), &group); err != nil {
		logger.Error("更新分组失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	logger.Info("更新分组成功", zap.Uint("id", uint(id)))
	response.SuccessWithMessage(c, "更新成功", nil)
}

func (h *HostGroupHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		logger.Error("删除分组失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	logger.Info("删除分组成功", zap.Uint("id", uint(id)))
	response.SuccessWithMessage(c, "删除成功", nil)
}

func (h *HostGroupHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	group, err := h.service.GetByID(uint(id))
	if err != nil {
		logger.Error("获取分组失败", zap.Error(err))
		response.NotFound(c, "分组不存在")
		return
	}

	response.Success(c, group)
}

func (h *HostGroupHandler) List(c *gin.Context) {
	var params model.QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	groups, total, err := h.service.List(&params)
	if err != nil {
		logger.Error("获取分组列表失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	response.PageSuccess(c, groups, total, params.GetPage(), params.GetPageSize())
}

func (h *HostGroupHandler) GetAll(c *gin.Context) {
	groups, err := h.service.GetAll()
	if err != nil {
		logger.Error("获取所有分组失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, groups)
}
