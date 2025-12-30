package main

import (
	"aiops-gateway/internal/config"
	"aiops-gateway/internal/registry"
	"aiops-gateway/internal/router"
	"aiops-gateway/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	printBanner()

	// 1. 加载配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 2. 初始化日志
	if err := logger.InitLogger(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output); err != nil {
		log.Fatalf("❌ 日志初始化失败: %v", err)
	}
	logger.Info("✅ 日志初始化成功")

	// 3. 连接 Consul
	consulClient, err := registry.NewConsulClient(cfg.Consul.Address, cfg.Consul.Scheme)
	if err != nil {
		logger.Fatal("❌ Consul 连接失败", zap.Error(err))
	}
	logger.Info("✅ Consul 连接成功", zap.String("address", cfg.Consul.Address))

	// 4. 启动服务发现
	discovery := registry.NewServiceDiscovery(consulClient)
	serviceNames := extractServiceNames(cfg.Routes)
	discovery.Start(serviceNames)

	// 5. 初始化路由
	r := router.SetupRouter(cfg, discovery)

	// 6. 创建 HTTP 服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 7. 启动服务器
	go func() {
		logger.Info("========================================")
		logger.Info("✅ API Gateway 启动成功！")
		logger.Info(fmt.Sprintf("📍 监听地址: http://0.0.0.0%s", addr))
		logger.Info(fmt.Sprintf("📍 健康检查: http://localhost%s/health", addr))
		logger.Info(fmt.Sprintf("📍 服务列表: http://localhost%s/services", addr))
		logger.Info("========================================")
		logger.Info("💡 按 Ctrl+C 优雅停止服务")
		logger.Info("========================================")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("❌ 服务启动失败", zap.Error(err))
		}
	}()

	// 8. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 收到停止信号，正在关闭服务器...")

	// 停止服务发现
	discovery.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器强制关闭", zap.Error(err))
	} else {
		logger.Info("✅ 服务器已安全退出")
	}
}

func printBanner() {
	banner := `
 ╔═══════════════════════════════════════════════════╗
 ║                                                   ║
 ║     🌐 AIOps API Gateway                         ║
 ║                                                   ║
 ║     Version: 1.0.0                               ║
 ║     Author:  OPS Team                            ║
 ║     Go:      1.21+                               ║
 ║                                                   ║
 ╚═══════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

func extractServiceNames(routes []config.RouteConfig) []string {
	serviceMap := make(map[string]bool)
	for _, route := range routes {
		serviceMap[route.ServiceName] = true
	}

	serviceNames := make([]string, 0, len(serviceMap))
	for name := range serviceMap {
		serviceNames = append(serviceNames, name)
	}
	return serviceNames
}
