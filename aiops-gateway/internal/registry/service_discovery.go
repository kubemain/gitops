package registry

import (
	"aiops-gateway/pkg/logger"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type ServiceDiscovery struct {
	consul   *ConsulClient
	services map[string][]*api.ServiceEntry
	mu       sync.RWMutex
	stopChan chan struct{}
}

func NewServiceDiscovery(consul *ConsulClient) *ServiceDiscovery {
	return &ServiceDiscovery{
		consul:   consul,
		services: make(map[string][]*api.ServiceEntry),
		stopChan: make(chan struct{}),
	}
}

// Start 启动服务发现
func (sd *ServiceDiscovery) Start(serviceNames []string) {
	for _, name := range serviceNames {
		// 初始化服务列表
		services, err := sd.consul.DiscoverService(name)
		if err != nil {
			logger.Warn("服务暂未注册，将在后台监听",
				zap.String("service", name),
				zap.Error(err),
			)
		} else {
			sd.updateServices(name, services)
			logger.Info("服务发现成功",
				zap.String("service", name),
				zap.Int("instances", len(services)),
			)
		}

		// 后台监听服务变化
		go sd.watchService(name)
	}
}

// Stop 停止服务发现
func (sd *ServiceDiscovery) Stop() {
	logger.Info("停止服务发现")
	close(sd.stopChan)

	// 清理资源
	sd.mu.Lock()
	sd.services = make(map[string][]*api.ServiceEntry)
	sd.mu.Unlock()
}

// watchService 监听单个服务
func (sd *ServiceDiscovery) watchService(serviceName string) {
	sd.consul.WatchService(serviceName, func(services []*api.ServiceEntry) {
		select {
		case <-sd.stopChan:
			logger.Info("停止监听服务", zap.String("service", serviceName))
			return
		default:
			sd.updateServices(serviceName, services)
			logger.Info("服务实例更新",
				zap.String("service", serviceName),
				zap.Int("instances", len(services)),
			)
		}
	})
}

// GetServiceAddress 获取服务地址（负载均衡）
func (sd *ServiceDiscovery) GetServiceAddress(serviceName string) (string, error) {
	sd.mu.RLock()
	services, exists := sd.services[serviceName]
	sd.mu.RUnlock()

	if !exists || len(services) == 0 {
		return "", fmt.Errorf("服务 %s 不可用", serviceName)
	}

	// 随机负载均衡
	rand.Seed(time.Now().UnixNano())
	service := services[rand.Intn(len(services))]
	address := fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port)

	return address, nil
}

// GetInstances 获取服务实例列表
func (sd *ServiceDiscovery) GetInstances(serviceName string) []*api.ServiceEntry {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	return sd.services[serviceName]
}

// updateServices 更新服务列表
func (sd *ServiceDiscovery) updateServices(serviceName string, services []*api.ServiceEntry) {
	sd.mu.Lock()
	defer sd.mu.Unlock()
	sd.services[serviceName] = services
}

// GetAllServices 获取所有服务
func (sd *ServiceDiscovery) GetAllServices() map[string]interface{} {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	result := make(map[string]interface{})
	for name, instances := range sd.services {
		serviceInfo := make([]map[string]interface{}, 0, len(instances))
		for _, instance := range instances {
			serviceInfo = append(serviceInfo, map[string]interface{}{
				"id":      instance.Service.ID,
				"address": instance.Service.Address,
				"port":    instance.Service.Port,
				"status":  instance.Checks.AggregatedStatus(),
			})
		}
		result[name] = serviceInfo
	}

	return result
}
