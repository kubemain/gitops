package main

import (
	"aiops-user-service/internal/config"
	"aiops-user-service/internal/registry"
	"aiops-user-service/internal/router"
	"aiops-user-service/pkg/database"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
)

func main() {
	printBanner()

	// 1. 加载配置文件
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 2. 初始化 MySQL 数据库
	log.Println("🔄 正在初始化 MySQL...")
	if err := database.InitMySQL(&cfg.MySQL); err != nil {
		log.Fatalf("❌ MySQL 初始化失败: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("⚠️  关闭 MySQL 连接失败: %v", err)
		}
	}()

	// 3. 初始化 Redis
	log.Println("🔄 正在初始化 Redis...")
	if err := database.InitRedis(&cfg.Redis); err != nil {
		log.Fatalf("❌ Redis 初始化失败: %v", err)
	}
	defer func() {
		if err := database.CloseRedis(); err != nil {
			log.Printf("⚠️  关闭 Redis 连接失败: %v", err)
		}
	}()

	// 4. 初始化路由
	log.Println("🔄 正在初始化路由...")
	r := router.SetupRouter()

	// 5. 获取服务地址
	serviceAddr := cfg.Server.Host
	if serviceAddr == "" || serviceAddr == "0.0.0.0" {
		// 自动获取本机 IP
		serviceAddr = getLocalIP()
		if serviceAddr == "" {
			log.Println("⚠️  无法获取本机 IP，使用 127.0.0.1")
			serviceAddr = "127.0.0.1"
		}
	}
	log.Printf("📍 服务地址: %s", serviceAddr)

	// 6. 注册到 Consul
	log.Println("🔄 正在注册到 Consul...")
	consulRegistry, err := registry.NewConsulRegistry(cfg.Consul.Address)
	if err != nil {
		log.Fatalf("❌ Consul 客户端创建失败: %v", err)
	}

	err = consulRegistry.Register(registry.ServiceConfig{
		ID:      fmt.Sprintf("aiops-user-service-%s-%d", serviceAddr, cfg.Server.Port),
		Name:    "aiops-user-service",
		Address: serviceAddr,
		Port:    cfg.Server.Port,
		Tags:    []string{"user", "auth", "role", "permission"},
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", serviceAddr, cfg.Server.Port),
			Interval:                       "10s",
			Timeout:                        "3s",
			DeregisterCriticalServiceAfter: "30s",
		},
	})

	if err != nil {
		log.Fatalf("❌ 服务注册失败: %v", err)
	}

	// 7. 创建 HTTP 服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 8. 在 Goroutine 中启动服务器
	go func() {
		log.Println("========================================")
		log.Printf("✅ User Service 启动成功！")
		log.Printf("📍 监听地址: http://0.0.0.0%s", addr)
		log.Printf("📍 注册地址: http://%s%s", serviceAddr, addr)
		log.Printf("📍 健康检查: http://%s%s/health", serviceAddr, addr)
		log.Printf("📍 API 接口: http://%s%s/api/v1", serviceAddr, addr)
		log.Println("========================================")
		log.Printf("💡 按 Ctrl+C 优雅停止服务")
		log.Println("========================================")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 服务启动失败: %v", err)
		}
	}()

	// 9. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n🛑 正在关闭服务器...")

	// 10. 注销 Consul 服务
	if err := consulRegistry.Deregister(); err != nil {
		log.Printf("⚠️  服务注销失败: %v", err)
	}

	// 11. 优雅停止
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("⚠️  服务器强制关闭: %v", err)
	}

	log.Println("✅ 服务器已安全退出")
}

// getLocalIP 获取本机 IP 地址
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}

func printBanner() {
	banner := `
 ╔═══════════════════════════════════════════════════╗
 ║                                                   ║
 ║     👤 AIOps User Service                        ║
 ║                                                   ║
 ║     Version: 1.0.0                               ║
 ║     Author:  OPS Team                            ║
 ║                                                   ║
 ╚═══════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}
