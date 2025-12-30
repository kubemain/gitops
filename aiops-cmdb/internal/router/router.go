package router

import (
	"aiops-cmdb/internal/config"
	"aiops-cmdb/internal/handler"
	"aiops-cmdb/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 健康检查（根路径，不在 API 路由下）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "aiops-cmdb",
		})
	})

	// ✅ API 路由 - 添加完整路径 /api/v1/cmdb
	api := r.Group("/api/v1/cmdb")

	// 主机分组路由
	groupHandler := handler.NewHostGroupHandler()
	groups := api.Group("/host-groups")
	{
		groups.GET("", groupHandler.List)
		groups.GET("/all", groupHandler.GetAll)
		groups.GET("/:id", groupHandler.GetByID)
		groups.POST("", groupHandler.Create)
		groups.PUT("/:id", groupHandler.Update)
		groups.DELETE("/:id", groupHandler.Delete)
	}

	// 主机路由
	hostHandler := handler.NewHostHandler()
	hosts := api.Group("/hosts")
	{
		hosts.GET("", hostHandler.List)
		hosts.GET("/stats", hostHandler.GetStats)
		hosts.GET("/:id", hostHandler.GetByID)
		hosts.POST("", hostHandler.Create)
		hosts.PUT("/:id", hostHandler.Update)
		hosts.DELETE("/:id", hostHandler.Delete)
		hosts.POST("/batch-delete", hostHandler.BatchDelete)
		hosts.PATCH("/:id/status", hostHandler.UpdateStatus)
	}

	return r
}
