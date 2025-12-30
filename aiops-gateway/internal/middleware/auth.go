package middleware

import (
	"aiops-gateway/internal/config"
	"aiops-gateway/pkg/logger"
	"aiops-gateway/pkg/response"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleIDs  []uint `json:"role_ids"`
	jwt.RegisteredClaims
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// OPTIONS 请求直接放行
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("请求缺少 Authorization 头",
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 解析 Token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logger.Warn("Token 格式错误",
				zap.String("auth_header", authHeader),
			)
			response.Unauthorized(c, "Token 格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 Token
		claims, err := parseToken(tokenString)
		if err != nil {
			logger.Warn("Token 验证失败",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
			)
			response.Unauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_ids", claims.RoleIDs)
		c.Set("token", tokenString)

		logger.Debug("JWT 验证成功",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
		)

		c.Next()
	}
}

func parseToken(tokenString string) (*Claims, error) {
	cfg := config.GlobalConfig
	if cfg == nil {
		return nil, errors.New("配置未初始化")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
