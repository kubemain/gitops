package proxy

import (
	"math/rand"
	"sync"
	"time"
)

type LoadBalancer interface {
	Next(services []string) string
}

// RandomLoadBalancer 随机负载均衡
type RandomLoadBalancer struct {
	rand *rand.Rand
	mu   sync.Mutex
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (lb *RandomLoadBalancer) Next(services []string) string {
	if len(services) == 0 {
		return ""
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	return services[lb.rand.Intn(len(services))]
}

// RoundRobinLoadBalancer 轮询负载均衡
type RoundRobinLoadBalancer struct {
	current int
	mu      sync.Mutex
}

func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{}
}

func (lb *RoundRobinLoadBalancer) Next(services []string) string {
	if len(services) == 0 {
		return ""
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	service := services[lb.current%len(services)]
	lb.current++
	return service
}
