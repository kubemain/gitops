package middleware

import (
	"aiops-gateway/pkg/response"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func RateLimit(requestsPerSecond int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 基于IP限流
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			response.Error(c, 429, "请求过于频繁，请稍后重试")
			c.Abort()
			return
		}

		c.Next()
	}
}
