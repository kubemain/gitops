package proxy

import (
	"aiops-gateway/internal/config"
	"aiops-gateway/internal/registry"
	"aiops-gateway/pkg/breaker"
	"aiops-gateway/pkg/logger"
	"aiops-gateway/pkg/response"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReverseProxy struct {
	discovery      *registry.ServiceDiscovery
	client         *http.Client
	breakerConfig  breaker.Config
	breakerEnabled bool
}

func NewReverseProxy(discovery *registry.ServiceDiscovery, cfg *config.Config) *ReverseProxy {
	// 优化 HTTP 客户端配置
	transport := &http.Transport{
		MaxIdleConns:        100,              // 最大空闲连接数
		MaxIdleConnsPerHost: 10,               // 每个 host 最大空闲连接
		IdleConnTimeout:     90 * time.Second, // 空闲连接超时
		DisableKeepAlives:   false,            // 启用 Keep-Alive
	}

	return &ReverseProxy{
		discovery: discovery,
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		breakerConfig: breaker.Config{
			MaxRequests: cfg.CircuitBreaker.MaxRequests,
			Interval:    time.Duration(cfg.CircuitBreaker.Interval) * time.Second,
			Timeout:     time.Duration(cfg.CircuitBreaker.Timeout) * time.Second,
		},
		breakerEnabled: cfg.CircuitBreaker.Enabled,
	}
}

// ProxyRequest 代理请求
func (p *ReverseProxy) ProxyRequest(c *gin.Context, route config.RouteConfig) {
	// 1. 获取目标服务地址
	targetURL, err := p.discovery.GetServiceAddress(route.ServiceName)
	if err != nil {
		logger.Error("服务发现失败",
			zap.String("service", route.ServiceName),
			zap.Error(err),
		)
		response.ServiceUnavailable(c, fmt.Sprintf("服务 %s 不可用", route.ServiceName))
		return
	}

	// 2. 构建目标URL
	path := c.Request.URL.Path

	// 处理路径参数
	if c.Param("path") != "" {
		path = c.Request.URL.Path
	}

	// 处理前缀剥离
	if route.StripPrefix {
		path = strings.TrimPrefix(path, route.Prefix)
	}

	targetURL = fmt.Sprintf("%s%s", targetURL, path)
	if c.Request.URL.RawQuery != "" {
		targetURL = fmt.Sprintf("%s?%s", targetURL, c.Request.URL.RawQuery)
	}

	// 3. 读取请求体
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("读取请求体失败", zap.Error(err))
			response.InternalServerError(c, "读取请求失败")
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// 4. 执行代理（带熔断）
	if p.breakerEnabled {
		p.proxyWithCircuitBreaker(c, targetURL, bodyBytes, route)
	} else {
		p.doProxy(c, targetURL, bodyBytes)
	}
}

// proxyWithCircuitBreaker 带熔断的代理
func (p *ReverseProxy) proxyWithCircuitBreaker(c *gin.Context, targetURL string, bodyBytes []byte, route config.RouteConfig) {
	result, err := breaker.Execute(route.ServiceName, p.breakerConfig, func() (interface{}, error) {
		return p.executeRequest(c, targetURL, bodyBytes)
	})

	if err != nil {
		logger.Error("熔断器执行失败",
			zap.String("service", route.ServiceName),
			zap.String("url", targetURL),
			zap.Error(err),
		)
		response.ServiceUnavailable(c, "服务暂时不可用，请稍后重试")
		return
	}

	resp := result.(*http.Response)
	p.writeResponse(c, resp)
}

// doProxy 直接代理（不带熔断）
func (p *ReverseProxy) doProxy(c *gin.Context, targetURL string, bodyBytes []byte) {
	resp, err := p.executeRequest(c, targetURL, bodyBytes)
	if err != nil {
		logger.Error("代理请求失败",
			zap.String("url", targetURL),
			zap.Error(err),
		)
		response.InternalServerError(c, "代理请求失败")
		return
	}

	p.writeResponse(c, resp)
}

// executeRequest 执行HTTP请求
func (p *ReverseProxy) executeRequest(c *gin.Context, targetURL string, bodyBytes []byte) (*http.Response, error) {
	// 创建请求
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 复制请求头（跳过不必要的头）
	for key, values := range c.Request.Header {
		// 跳过这些头，让 HTTP 客户端自动处理
		if key == "Content-Length" || key == "Host" || key == "Connection" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 添加追踪ID
	if traceID := c.GetString("trace_id"); traceID != "" {
		req.Header.Set("X-Trace-ID", traceID)
	}

	// 添加真实 IP
	req.Header.Set("X-Real-IP", c.ClientIP())
	req.Header.Set("X-Forwarded-For", c.ClientIP())

	// 添加用户信息（如果存在）
	if userID, exists := c.Get("user_id"); exists {
		req.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
	}
	if username, exists := c.Get("username"); exists {
		req.Header.Set("X-Username", fmt.Sprintf("%v", username))
	}

	// 执行请求
	start := time.Now()
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	duration := time.Since(start)

	logger.Info("代理请求",
		zap.String("method", c.Request.Method),
		zap.String("url", targetURL),
		zap.Int("status", resp.StatusCode),
		zap.Duration("duration", duration),
		zap.String("trace_id", c.GetString("trace_id")),
	)

	return resp, nil
}

// writeResponse 写入响应
func (p *ReverseProxy) writeResponse(c *gin.Context, resp *http.Response) {
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// ✅ 强制设置 Content-Type 为 UTF-8
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "charset") {
		if strings.Contains(contentType, "application/json") {
			c.Header("Content-Type", "application/json; charset=utf-8")
		} else if strings.Contains(contentType, "text/") {
			c.Header("Content-Type", contentType+"; charset=utf-8")
		}
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取响应体失败", zap.Error(err))
		response.InternalServerError(c, "读取响应失败")
		return
	}

	// 写入响应
	c.Data(resp.StatusCode, c.GetHeader("Content-Type"), body)
}
