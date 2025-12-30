# 🌐 AIOps API Gateway

## 功能特性

✅ **服务路由** - 统一入口，智能转发  
✅ **服务发现** - Consul集成，自动感知服务变化  
✅ **负载均衡** - 随机/轮询算法  
✅ **熔断降级** - Circuit Breaker模式  
✅ **限流保护** - Token Bucket算法  
✅ **认证鉴权** - JWT Token验证  
✅ **链路追踪** - TraceID传递  
✅ **结构化日志** - Zap日志库

## 快速开始

### 1. 安装依赖
```bash
go mod download


2. 启动 Consul
bash
docker run -d --name consul \
  -p 8500:8500 \
  consul:latest agent -dev -ui -client=0.0.0.0
3. 修改配置
编辑 configs/config.yaml

4. 启动服务
bash
go run cmd/server/main.go
5. 测试
bash
# 健康检查
curl http://localhost:8000/health

# 查看服务列表
curl http://localhost:8000/services

# 登录（转发到 user-service）
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
配置说明
路由配置
yaml
routes:
  - name: "user-service"
    prefix: "/api/v1/users"
    service_name: "aiops-user-service"
    strip_prefix: false
    auth_required: true
name: 路由名称
prefix: 路径前缀
service_name: 目标服务在Consul中的名称
strip_prefix: 是否去除前缀
auth_required: 是否需要JWT认证
架构图
code
Client → API Gateway → Consul → Backend Services
          ↓
       [JWT验证]
       [限流]
       [熔断]
       [负载均衡]
License
MIT

code

---

## 🎯 使用说明

### 1. 启动 Consul
```bash
docker run -d --name consul \
  -p 8500:8500 \
  consul:latest agent -dev -ui -client=0.0.0.0
2. 启动 Gateway
bash
cd aiops-gateway
go mod tidy
go run cmd/server/main.go
