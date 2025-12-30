// internal/middleware/permission.go
package middleware

import (
	"aiops-user-service/internal/service"
	"aiops-user-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RequirePermission 权限校验中间件
func RequirePermission(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 2. 检查权限
		authService := service.NewAuthService()
		hasPermission, err := authService.CheckPermission(userID.(uint), permCode)
		if err != nil {
			utils.InternalServerError(c, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}
