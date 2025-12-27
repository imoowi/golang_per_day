package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"day53/limiter"
	"day53/middleware"

	"github.com/gin-gonic/gin"
)

// main 函数是程序的入口点
func main() {
	// 打印欢迎信息
	fmt.Println("Hello, Codee君!\nWelcome to golang_per_day")
	r := gin.Default()

	// 创建令牌桶限流器：容量为10，每秒生成5个令牌
	tb := limiter.NewTokenBucket(10, 5)
	// 将限流中间件应用到所有路由
	r.Use(middleware.RateLimit(tb))

	r.GET("/profile/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 模拟下游服务返回 503 触发降级
		if rand.Float32() < 0.4 {
			// 触发降级
			c.JSON(http.StatusOK, gin.H{
				"user_id": id,
				"name":    "anonymous",
				"remark":  "degraded response",
			})
			return
		}

		// 正常完整响应
		time.Sleep(100 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{
			"user_id": id,
			"name":    "Tom",
			"age":     28,
			"email":   "tom@example.com",
		})
	})

	r.Run(":8080")
}
