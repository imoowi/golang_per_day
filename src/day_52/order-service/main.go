package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

// client 是用于发送 HTTP 请求的客户端
var client = &http.Client{}

// cb 是用于调用用户服务的断路器实例
var cb *gobreaker.CircuitBreaker

// init 函数在程序启动时初始化断路器
func init() {
	// 配置断路器参数
	settings := gobreaker.Settings{
		Name:        "UserService",   // 断路器名称
		MaxRequests: 3,               // 半开状态下允许的最大请求数
		Interval:    5 * time.Second, // 统计周期
		Timeout:     5 * time.Second, // 断路器从打开到半开的超时时间
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// 当连续失败次数达到3次时，断路器打开
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			// 记录断路器状态变化
			log.Printf("Circuit breaker state changed: %s -> %s\n", from.String(), to.String())
		},
	}
	// 创建断路器实例
	cb = gobreaker.NewCircuitBreaker(settings)
}

// callUserService 通过断路器调用用户服务获取用户信息
// id: 用户ID
// 返回用户信息的字节数组和可能的错误
func callUserService(id string) ([]byte, error) {
	// 创建带有2秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 构建用户服务的请求URL
	url := fmt.Sprintf("http://localhost:8081/user/%s", id)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	// 使用断路器执行请求
	body, err := cb.Execute(func() (interface{}, error) {
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	})

	if err != nil {
		return nil, err
	}
	return body.([]byte), nil
}

// main 函数是程序的入口点
func main() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	// 创建Gin引擎实例
	r := gin.Default()

	// 定义获取订单的API端点
	r.GET("/order/:id", func(c *gin.Context) {
		// 获取订单ID参数
		id := c.Param("id")
		// 调用用户服务获取用户信息
		data, err := callUserService(id)
		if err != nil {
			// 如果调用失败，返回服务不可用状态
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}

		// 解析用户信息JSON
		var user map[string]interface{}
		json.Unmarshal(data, &user)

		// 构造订单响应
		response := gin.H{
			"order_id": fmt.Sprintf("order-%d", rand.Intn(1000)), // 生成随机订单ID
			"user":     user,                                     // 包含用户信息
		}

		// 返回成功响应
		c.JSON(http.StatusOK, response)
	})

	// 启动订单服务，监听8080端口
	fmt.Println("Order Service running at :8080")
	r.Run(":8080")
}
