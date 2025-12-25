package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

const (
	srvName = "golang_per_day_50"
	srvHost = "192.168.1.31" //这里写自己的 IP 地址，否则docker里的Consul无法访问到
	srvPort = 8080
	srvID   = "golang_per_day_50-192.168.1.31:8080"
)

var (
	appConfig *Config
	lock      sync.RWMutex
)

type Config struct {
	Consul ConsulConfig `mapstructure:"consul" json:"consul"`
	Redis  RedisConfig  `mapstructure:"redis" json:"redis"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Db       int    `mapstructure:"db" json:"db"`
	Password string `mapstructure:"password" json:"password"`
}
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

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
	// 2. 读取配置中心的配置，这里用Consul
	config := api.DefaultConfig()
	config.Address = "192.168.1.31:8500" // Consul 的地址
	client, _ := api.NewClient(config)
	kvPath := "config/golang_per_day/50"
	loadConfig(client, kvPath)
	go watchConfig(client, kvPath)
	// 3. 注册服务到 Consul
	registerToConsul()
	// 4. 启动 Gin 服务
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", srvPort)); err != nil {
			log.Fatal("服务启动失败: ", err)
		}
	}()

	// 5. 优雅退出：监听信号注销服务
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
	config.Address = fmt.Sprintf("%s:%d", appConfig.Consul.Host, appConfig.Consul.Port) // Consul 的地址
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

// 读取配置中心的 Redis 配置
func loadConfig(client *api.Client, path string) {
	pair, _, err := client.KV().Get(path, nil)
	if err != nil {
		log.Fatal("读取配置失败: ", err)
	}
	v := viper.New()
	v.SetConfigType("yaml")
	err = v.ReadConfig(bytes.NewBuffer(pair.Value))
	if err != nil {
		log.Fatal("读取配置失败: ", err)
	}
	newConfig := Config{}
	err = v.Unmarshal(&newConfig)
	if err != nil {
		log.Fatal("解析配置失败: ", err)
	}
	lock.Lock()
	defer lock.Unlock()
	appConfig = &newConfig
	log.Println("Redis 配置更新: ", appConfig)
}

// 监听 Consul K/V 变化
func watchConfig(client *api.Client, path string) {
	var lastIndex uint64
	for {
		// Consul 的 WaitIndex 机制：如果有变化会立即返回，否则阻塞直到超时
		pair, meta, err := client.KV().Get(path, &api.QueryOptions{
			WaitIndex: lastIndex,
		})
		if err != nil {
			log.Printf("监听出错: %v", err)
			continue
		}

		if pair != nil && meta.LastIndex > lastIndex {
			lastIndex = meta.LastIndex
			loadConfig(client, path)
		}
	}
}
