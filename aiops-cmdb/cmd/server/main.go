package main

import (
	"aiops-cmdb/internal/config"
	"aiops-cmdb/internal/model"
	"aiops-cmdb/internal/router"
	"aiops-cmdb/pkg/database"
	"aiops-cmdb/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
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

	// 3. 连接数据库
	if err := database.InitMySQL(&cfg.Database); err != nil {
		logger.Fatal("❌ 数据库连接失败", zap.Error(err))
	}

	// 4. 自动迁移数据表
	db := database.GetDB()
	if err := db.AutoMigrate(
		&model.Host{},
		&model.HostGroup{},
		&model.HostChange{},
	); err != nil {
		logger.Fatal("❌ 数据表迁移失败", zap.Error(err))
	}
	logger.Info("✅ 数据表迁移成功")

	// 5. 注册到 Consul
	consulClient, err := api.NewClient(&api.Config{
		Address: cfg.Consul.Address,
		Scheme:  cfg.Consul.Scheme,
	})
	if err != nil {
		logger.Fatal("❌ Consul 连接失败", zap.Error(err))
	}

	registration := &api.AgentServiceRegistration{
		ID:      cfg.Consul.ServiceID,
		Name:    cfg.Consul.ServiceName,
		Address: cfg.Consul.ServiceAddress,
		Port:    cfg.Consul.ServicePort,
		Tags:    []string{"cmdb", "aiops"},
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", cfg.Consul.ServiceAddress, cfg.Consul.ServicePort),
			Interval:                       cfg.Consul.HealthCheckInterval,
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: cfg.Consul.DeregisterCriticalServiceAfter,
		},
	}

	if err := consulClient.Agent().ServiceRegister(registration); err != nil {
		logger.Fatal("❌ 服务注册失败", zap.Error(err))
	}
	logger.Info("✅ 服务注册成功", zap.String("service_id", cfg.Consul.ServiceID))

	// 6. 初始化路由
	r := router.SetupRouter(cfg)

	// 7. 启动 HTTP 服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info("========================================")
		logger.Info("✅ CMDB 服务启动成功！")
		logger.Info(fmt.Sprintf("📍 监听地址: http://0.0.0.0%s", addr))
		logger.Info(fmt.Sprintf("📍 健康检查: http://localhost%s/health", addr))
		logger.Info(fmt.Sprintf("📍 API 文档: http://localhost%s/api/v1/hosts", addr))
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

	logger.Info("🛑 收到停止信号，正在关闭服务...")

	// 注销服务
	if err := consulClient.Agent().ServiceDeregister(cfg.Consul.ServiceID); err != nil {
		logger.Error("服务注销失败", zap.Error(err))
	} else {
		logger.Info("✅ 服务注销成功")
	}

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
 ║     📦 AIOps CMDB Service                        ║
 ║                                                   ║
 ║     Version: 1.0.0                               ║
 ║     Author:  OPS Team                            ║
 ║     Go:      1.21+                               ║
 ║                                                   ║
 ╚═══════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}
