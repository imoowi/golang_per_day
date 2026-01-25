# 第一阶段：编译阶段
# 使用 golang:1.25.5-alpine 作为基础镜像，体积小且包含 Go 编译环境
FROM golang:1.25.5-alpine AS builder

# 设置环境变量
# GO111MODULE=on：启用 Go 模块
# GOPROXY：设置 Go 模块代理，国内环境使用 goproxy.cn 加快下载速度
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制依赖文件并下载（利用 Docker 缓存）
# 先复制 go.mod 和 go.sum，只有当依赖文件变化时才会重新下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并编译
COPY . .
# CGO_ENABLED=0：禁用 CGO，确保静态链接，能在 alpine 下运行
# GOOS=linux GOARCH=amd64：设置目标操作系统和架构
# go build -o main .：编译生成名为 main 的可执行文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 第二阶段：运行阶段
# 使用 alpine:latest 作为基础镜像，体积非常小
FROM alpine:latest

# 安装基础工具（可选）
# ca-certificates：包含 SSL 证书，用于 HTTPS 连接
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从编译阶段复制二进制文件
# 使用 --from=builder 从之前的构建阶段复制文件
COPY --from=builder /app/main .

# 暴露应用监听的端口
# 8080：Gin 服务器监听的端口
EXPOSE 8080

# 运行应用
# 启动编译生成的 main 可执行文件
CMD ["./main"]