package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig         `mapstructure:"server"`
	Consul         ConsulConfig         `mapstructure:"consul"`
	Redis          RedisConfig          `mapstructure:"redis"`
	JWT            JWTConfig            `mapstructure:"jwt"`
	RateLimit      RateLimitConfig      `mapstructure:"ratelimit"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
	Log            LogConfig            `mapstructure:"log"`
	Trace          TraceConfig          `mapstructure:"trace"`
	Routes         []RouteConfig        `mapstructure:"routes"`
}

type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type ConsulConfig struct {
	Address string `mapstructure:"address"`
	Scheme  string `mapstructure:"scheme"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerSecond int  `mapstructure:"requests_per_second"`
	Burst             int  `mapstructure:"burst"`
}

type CircuitBreakerConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	MaxRequests uint32 `mapstructure:"max_requests"`
	Interval    int    `mapstructure:"interval"`
	Timeout     int    `mapstructure:"timeout"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type TraceConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	JaegerEndpoint string `mapstructure:"jaeger_endpoint"`
	ServiceName    string `mapstructure:"service_name"`
}

type RouteConfig struct {
	Name         string `mapstructure:"name"`
	Prefix       string `mapstructure:"prefix"`
	ServiceName  string `mapstructure:"service_name"`
	StripPrefix  bool   `mapstructure:"strip_prefix"`
	AuthRequired bool   `mapstructure:"auth_required"`
}

var GlobalConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 配置验证
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	GlobalConfig = &config
	log.Printf("✅ 配置加载成功: %s", configPath)
	return &config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("服务器端口无效: %d", c.Server.Port)
	}

	if c.Server.Mode != "debug" && c.Server.Mode != "release" {
		return fmt.Errorf("服务器模式无效: %s (应为 debug 或 release)", c.Server.Mode)
	}

	// 验证 Consul 配置
	if c.Consul.Address == "" {
		return fmt.Errorf("Consul 地址不能为空")
	}

	// 验证 JWT 配置
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT Secret 不能为空")
	}
	if c.JWT.Secret == "ops-platform-secret-key-change-in-production" {
		log.Println("⚠️  警告：正在使用默认 JWT Secret，生产环境请务必修改！")
	}
	if c.JWT.ExpireHours <= 0 {
		return fmt.Errorf("JWT 过期时间必须大于 0")
	}

	// 验证路由配置
	if len(c.Routes) == 0 {
		return fmt.Errorf("至少需要配置一个路由")
	}

	for i, route := range c.Routes {
		if route.Name == "" {
			return fmt.Errorf("路由 [%d] 名称不能为空", i)
		}
		if route.Prefix == "" {
			return fmt.Errorf("路由 [%s] 前缀不能为空", route.Name)
		}
		if route.ServiceName == "" {
			return fmt.Errorf("路由 [%s] 服务名不能为空", route.Name)
		}
	}

	return nil
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
