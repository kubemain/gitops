// internal/middleware/auth.go
package middleware

import (
	"aiops-user-service/pkg/database"
	"aiops-user-service/pkg/utils"
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT 认证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 2. 解析 Token（格式：Bearer <token>）
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "Token 格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 3. 检查 Token 是否在黑名单
		ctx := context.Background()
		blacklistKey := fmt.Sprintf("token:blacklist:%s", token)
		if database.RDB.Exists(ctx, blacklistKey).Val() > 0 {
			utils.Unauthorized(c, "Token 已失效")
			c.Abort()
			return
		}

		// 4. 解析 Token
		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.Unauthorized(c, "Token 无效")
			c.Abort()
			return
		}

		// 5. 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_ids", claims.RoleIDs)
		c.Set("token", token)

		c.Next()
	}
}
