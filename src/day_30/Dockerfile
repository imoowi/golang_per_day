# 第一阶段：构建阶段
FROM golang:1.23-alpine AS builder
#LABEL 指令用来给镜像添加一些元数据（metadata），以键值对的形式
LABEL maintainer="Codee君"

#设置容器Go语言环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

#为 RUN、CMD、ENTRYPOINT、COPY 和 ADD 设置工作目录，就是切换目录
WORKDIR /go/release

#COPY 拷贝文件或目录到容器中，跟ADD类似，但不具备自动下载或解压的功能
COPY . .

#RUN 构建镜像时运行的指令
#sed 命令用于文本处理，这里是用来替换 Alpine Linux 的软件源镜像为中科大的镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add tzdata
#安装swago，并自动生成API文档
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.10 
RUN swag init  --parseDependency=true

#go build 编译程序
RUN CGO_ENABLED=0 GOOS=linux go build -p 1 -ldflags="-w -s" -a -installsuffix cgo -o golang_per_day_30 .


# 第二阶段：运行阶段
# FROM alpine:latest
# # 把第一阶段构建的二进制文件和依赖拷贝到新的镜像中
# COPY --from=builder /go/release/golang_per_day_30 /

# COPY --from=builder /go/release/configs /configs

# COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# #EXPOSE 声明容器的服务端口（仅仅是声明）
# EXPOSE 8080

# #CMD 运行容器时执行的shell命令
# CMD ["/golang_per_day_30","server","-c", "/configs/config.yaml"]

#极限瘦身，只保留二进制文件和配置文件
FROM scratch
COPY --from=builder /go/release/golang_per_day_30 /
CMD ["/golang_per_day_30","server","-c", "/configs/config.docker.yaml"]
