// internal/router/router.go
package router

import (
	"aiops-user-service/internal/handler"
	"aiops-user-service/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())      // 恢复panic
	r.Use(middleware.Logger()) // 日志
	r.Use(middleware.CORS())   // 跨域

	// 初始化处理器
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	roleHandler := handler.NewRoleHandler()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// API v1 版本组
	v1 := r.Group("/api/v1")
	{
		// ==================== 公开接口（不需要认证）====================
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login) // 登录
		}

		// ==================== 需要认证的接口 ====================
		authRequired := v1.Group("")
		authRequired.Use(middleware.JWT()) // 应用 JWT 认证中间件
		{
			// 认证相关
			auth := authRequired.Group("/auth")
			{
				auth.POST("/logout", authHandler.Logout)       // 登出
				auth.GET("/userinfo", authHandler.GetUserInfo) // 获取当前用户信息
			}

			// ==================== 用户管理 ====================
			users := authRequired.Group("/users")
			{
				// 查看用户列表（需要权限）
				users.GET("",
					middleware.RequirePermission("system:user:view"),
					userHandler.List)

				// 查看用户详情
				users.GET("/:id",
					middleware.RequirePermission("system:user:view"),
					userHandler.GetByID)

				// 创建用户
				users.POST("",
					middleware.RequirePermission("system:user:create"),
					userHandler.Create)

				// 更新用户
				users.PUT("/:id",
					middleware.RequirePermission("system:user:edit"),
					userHandler.Update)

				// 删除用户
				users.DELETE("/:id",
					middleware.RequirePermission("system:user:delete"),
					userHandler.Delete)

				// 重置密码
				users.PUT("/:id/password",
					middleware.RequirePermission("system:user:edit"),
					userHandler.ResetPassword)
			}

			// ==================== 角色管理 ====================
			roles := authRequired.Group("/roles")
			{
				// 角色列表
				roles.GET("",
					middleware.RequirePermission("system:role:view"),
					roleHandler.List)

				// 角色详情
				roles.GET("/:id",
					middleware.RequirePermission("system:role:view"),
					roleHandler.GetByID)

				// 分配权限
				roles.POST("/:id/permissions",
					middleware.RequirePermission("system:role:assign"),
					roleHandler.AssignPermissions)
			}

			// ==================== 权限管理 ====================
			permissions := authRequired.Group("/permissions")
			{
				// 获取所有权限（树形结构）
				permissions.GET("",
					middleware.RequirePermission("system:role:view"),
					roleHandler.GetAllPermissions)
			}
		}
	}

	return r
}
