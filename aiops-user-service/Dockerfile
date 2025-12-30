# ==================== 第一阶段：构建 ====================
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装编译依赖
RUN apk add --no-cache git gcc musl-dev

# 设置 Go 环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译（-ldflags="-w -s" 减小体积）
RUN go build -ldflags="-w -s" -o aiops-api cmd/server/main.go

# ==================== 第二阶段：运行 ====================
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /build/aiops-api /app/aiops-api

# 复制配置文件
COPY --chown=app:app configs /app/configs

# 创建日志目录
RUN mkdir -p /app/logs && chown -R app:app /app

# 切换到非 root 用户
USER app

# 暴露端口
EXPOSE 8000

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8000/health || exit 1

# 启动服务
ENTRYPOINT ["/app/aiops-api"]
