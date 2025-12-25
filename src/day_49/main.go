package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

const (
	srvName = "golang_per_day_49"
	srvHost = "192.168.1.31" //这里写自己的 IP 地址，否则docker里的Consul无法访问到
	srvPort = 8080
	srvID   = "golang_per_day_49-192.168.1.31:8080"
)

func main() {
	r := gin.Default()
	// 定义一个 ping 接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 1. 定义健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	// 2. 注册服务到 Consul
	registerToConsul()

	// 3. 启动 Gin 服务
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", srvPort)); err != nil {
			log.Fatal("服务启动失败: ", err)
		}
	}()

	// 4. 优雅退出：监听信号注销服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 退出前注销服务
	deregisterFromConsul()
	log.Println("服务已退出并从 Consul 注销")
}

// 注册逻辑
func registerToConsul() {
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500" // Consul 的地址
	client, _ := api.NewClient(config)

	registration := &api.AgentServiceRegistration{
		ID:      srvID,
		Name:    srvName,
		Address: srvHost,
		Port:    srvPort,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", srvHost, srvPort),
			Interval: "5s", // 每 5 秒检查一次
			Timeout:  "3s", // 超时时间
		},
	}

	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("注册失败: ", err)
	}
	log.Println("服务注册成功!")
}

// 注销逻辑
func deregisterFromConsul() {
	config := api.DefaultConfig()
	client, _ := api.NewClient(config)
	_ = client.Agent().ServiceDeregister(srvID)
}
