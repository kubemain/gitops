package registry

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
	client     *api.Client
	serviceID  string
	registered bool
}

type ServiceConfig struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Check   *api.AgentServiceCheck
}

func NewConsulRegistry(consulAddr string) (*ConsulRegistry, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建 Consul 客户端失败: %w", err)
	}

	return &ConsulRegistry{
		client: client,
	}, nil
}

// Register 注册服务到 Consul
func (c *ConsulRegistry) Register(cfg ServiceConfig) error {
	registration := &api.AgentServiceRegistration{
		ID:      cfg.ID,
		Name:    cfg.Name,
		Address: cfg.Address,
		Port:    cfg.Port,
		Tags:    cfg.Tags,
		Check:   cfg.Check,
	}

	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("服务注册失败: %w", err)
	}

	c.serviceID = cfg.ID
	c.registered = true

	log.Printf("✅ 服务注册成功: %s (ID: %s, Address: %s:%d)", cfg.Name, cfg.ID, cfg.Address, cfg.Port)

	// 监听退出信号，自动注销服务
	go c.handleShutdown()

	return nil
}

// Deregister 注销服务
func (c *ConsulRegistry) Deregister() error {
	if !c.registered {
		return nil
	}

	if err := c.client.Agent().ServiceDeregister(c.serviceID); err != nil {
		return fmt.Errorf("服务注销失败: %w", err)
	}

	log.Printf("✅ 服务注销成功: %s", c.serviceID)
	c.registered = false
	return nil
}

// handleShutdown 监听退出信号
func (c *ConsulRegistry) handleShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 收到退出信号，正在注销服务...")
	if err := c.Deregister(); err != nil {
		log.Printf("❌ 服务注销失败: %v", err)
	}
}
