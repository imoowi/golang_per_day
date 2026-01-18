# 第一阶段：编译阶段
FROM golang:1.25.5-alpine AS builder

# 设置环境变量（国内环境建议开启代理加快下载速度）
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制依赖文件并下载（利用 Docker 缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并编译
COPY . .
# CGO_ENABLED=0 确保静态链接，能在 alpine 下运行
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 第二阶段：运行阶段
FROM alpine:latest

# 安装基础工具（可选）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从编译阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露 Gin 和 Prometheus 监听的端口
EXPOSE 8080

# 运行应用
CMD ["./main"]