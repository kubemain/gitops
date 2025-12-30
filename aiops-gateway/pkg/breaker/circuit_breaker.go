package breaker

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

var (
	breakers = make(map[string]*gobreaker.CircuitBreaker)
	mu       sync.RWMutex
)

type Config struct {
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
}

// GetCircuitBreaker 获取或创建熔断器
func GetCircuitBreaker(name string, config Config) *gobreaker.CircuitBreaker {
	mu.RLock()
	cb, exists := breakers[name]
	mu.RUnlock()

	if exists {
		return cb
	}

	mu.Lock()
	defer mu.Unlock()

	// 双重检查
	if cb, exists := breakers[name]; exists {
		return cb
	}

	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("🔄 熔断器 [%s] 状态变化: %s -> %s\n", name, from, to)
		},
	}

	cb = gobreaker.NewCircuitBreaker(settings)
	breakers[name] = cb
	return cb
}

// Execute 执行带熔断的函数
func Execute(name string, config Config, fn func() (interface{}, error)) (interface{}, error) {
	cb := GetCircuitBreaker(name, config)
	return cb.Execute(fn)
}
