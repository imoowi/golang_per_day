package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// main 函数是程序的入口点
func main() {
	// 创建Gin引擎实例
	r := gin.Default()
	// Go 1.20+ 无需手动设置随机数种子，默认已自动初始化

	// 定义获取用户信息的API端点
	r.GET("/user/:id", func(c *gin.Context) {
		// 获取用户ID参数
		id := c.Param("id")
		// 模拟随机延迟和失败
		delay := rand.Intn(3) // 生成0~2秒的随机延迟
		time.Sleep(time.Duration(delay) * time.Second)

		// 30%概率模拟服务失败
		if rand.Float32() < 0.3 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user service failed"})
			return
		}

		// 返回成功响应，包含用户ID和用户名
		c.JSON(http.StatusOK, gin.H{"id": id, "name": fmt.Sprintf("User%s", id)})
	})

	// 启动用户服务，监听8081端口
	fmt.Println("User Service running at :8081")
	r.Run(":8081")
}
