package router

import (
	"aiops-gateway/internal/config"
	"aiops-gateway/internal/middleware"
	"aiops-gateway/internal/proxy"
	"aiops-gateway/internal/registry"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, discovery *registry.ServiceDiscovery) *gin.Engine {
	// 设置运行模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件（顺序很重要）
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())   // CORS 必须在最前面
	r.Use(middleware.Trace())  // 追踪
	r.Use(middleware.Logger()) // 日志

	// 限流中间件
	if cfg.RateLimit.Enabled {
		r.Use(middleware.RateLimit(cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API Gateway 运行正常",
		})
	})

	// 服务状态
	r.GET("/services", func(c *gin.Context) {
		services := discovery.GetAllServices()
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    services,
		})
	})

	// 创建反向代理
	reverseProxy := proxy.NewReverseProxy(discovery, cfg)

	// 配置路由
	for _, route := range cfg.Routes {
		setupRoute(r, route, reverseProxy)
	}

	return r
}

func setupRoute(r *gin.Engine, route config.RouteConfig, reverseProxy *proxy.ReverseProxy) {
	handler := func(c *gin.Context) {
		reverseProxy.ProxyRequest(c, route)
	}

	// 同时注册根路径和通配符路径
	if route.AuthRequired {
		// 需要认证的路由
		r.Any(route.Prefix, middleware.JWT(), handler)
		r.Any(route.Prefix+"/*path", middleware.JWT(), handler)
	} else {
		// 不需要认证的路由
		r.Any(route.Prefix, handler)
		r.Any(route.Prefix+"/*path", handler)
	}
}
