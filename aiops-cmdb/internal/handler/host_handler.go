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

type HostHandler struct {
	service *service.HostService
}

func NewHostHandler() *HostHandler {
	return &HostHandler{
		service: service.NewHostService(),
	}
}

// Create 创建主机
func (h *HostHandler) Create(c *gin.Context) {
	var host model.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		logger.Error("参数绑定失败", zap.Error(err))
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Create(&host); err != nil {
		logger.Error("创建主机失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	logger.Info("创建主机成功", zap.String("hostname", host.Hostname))
	response.SuccessWithMessage(c, "创建成功", host)
}

// Update 更新主机
func (h *HostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var host model.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Update(uint(id), &host); err != nil {
		logger.Error("更新主机失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	logger.Info("更新主机成功", zap.Uint("id", uint(id)))
	response.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除主机
func (h *HostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		logger.Error("删除主机失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	logger.Info("删除主机成功", zap.Uint("id", uint(id)))
	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetByID 获取主机详情
func (h *HostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	host, err := h.service.GetByID(uint(id))
	if err != nil {
		logger.Error("获取主机失败", zap.Error(err))
		response.NotFound(c, "主机不存在")
		return
	}

	response.Success(c, host)
}

// List 获取主机列表
func (h *HostHandler) List(c *gin.Context) {
	var params model.HostQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	hosts, total, err := h.service.List(&params)
	if err != nil {
		logger.Error("获取主机列表失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	response.PageSuccess(c, hosts, total, params.GetPage(), params.GetPageSize())
}

// GetStats 获取统计数据
func (h *HostHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		logger.Error("获取统计数据失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, stats)
}

// BatchDelete 批量删除
func (h *HostHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.service.BatchDelete(req.IDs); err != nil {
		logger.Error("批量删除失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	logger.Info("批量删除成功", zap.Int("count", len(req.IDs)))
	response.SuccessWithMessage(c, "删除成功", nil)
}

// UpdateStatus 更新状态
func (h *HostHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		logger.Error("更新状态失败", zap.Error(err))
		response.Error(c, 500, err.Error())
		return
	}

	logger.Info("更新状态成功", zap.Uint("id", uint(id)), zap.String("status", req.Status))
	response.SuccessWithMessage(c, "更新成功", nil)
}
