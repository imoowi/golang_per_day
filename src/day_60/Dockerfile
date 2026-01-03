#-------------构建阶段--------------
# 使用 Go 1.25.5 Alpine 版本作为构建镜像
FROM golang:1.25.5-alpine AS builder
# 设置工作目录为 /app
WORKDIR /app
# 复制当前目录下的所有文件到容器中
COPY . .
# 编译 Go 程序，生成可执行文件 main
RUN go build -o main .

#-------------运行阶段--------------
# 使用最新的 Alpine 镜像作为运行时镜像
FROM alpine:latest
# 设置工作目录为 /app
WORKDIR /app
# 暴露 8080 端口
EXPOSE 8080
# 从构建阶段复制编译好的可执行文件到运行时镜像
COPY --from=builder /app/main /app/golang_per_day
# 运行程序
CMD ["/app/golang_per_day"]
