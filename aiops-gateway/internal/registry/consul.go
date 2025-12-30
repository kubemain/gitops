package registry

import (
	"aiops-gateway/pkg/logger"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type ConsulClient struct {
	client *api.Client
}

func NewConsulClient(address, scheme string) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = address
	config.Scheme = scheme

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建 Consul 客户端失败: %w", err)
	}

	return &ConsulClient{client: client}, nil
}

// DiscoverService 服务发现
func (c *ConsulClient) DiscoverService(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("服务发现失败: %w", err)
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("服务 %s 不存在或不健康", serviceName)
	}

	return services, nil
}

// GetServiceAddress 获取服务地址
func (c *ConsulClient) GetServiceAddress(serviceName string) (string, error) {
	services, err := c.DiscoverService(serviceName)
	if err != nil {
		return "", err
	}

	service := services[0]
	address := fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port)
	return address, nil
}

// WatchService 监听服务变化
func (c *ConsulClient) WatchService(serviceName string, callback func([]*api.ServiceEntry)) {
	var lastIndex uint64
	retryCount := 0
	const maxRetries = 5
	const retryInterval = 5 * time.Second

	for {
		services, meta, err := c.client.Health().Service(serviceName, "", true, &api.QueryOptions{
			WaitIndex: lastIndex,
			WaitTime:  10 * time.Second,
		})

		if err != nil {
			retryCount++
			if retryCount >= maxRetries {
				logger.Error("服务监听失败次数过多，停止监听",
					zap.String("service", serviceName),
					zap.Int("retry_count", retryCount),
					zap.Error(err),
				)
				return
			}

			logger.Warn("监听服务失败，将重试",
				zap.String("service", serviceName),
				zap.Int("retry_count", retryCount),
				zap.Duration("retry_after", retryInterval),
				zap.Error(err),
			)
			time.Sleep(retryInterval)
			continue
		}

		// 重置重试计数
		retryCount = 0

		if meta.LastIndex != lastIndex {
			lastIndex = meta.LastIndex
			callback(services)
		}
	}
}
