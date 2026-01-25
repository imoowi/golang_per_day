package main

import (
	"github.com/gin-gonic/gin"
)

// main 函数是应用程序的入口点
// 创建一个 Gin 路由器并定义 API 端点
func main() {
	// 创建 Gin 路由器，默认包含 Logger 和 Recovery 中间件
	r := gin.Default()

	// 定义 GET /user/:id 路由
	// 该路由处理用户信息查询请求
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param(`id`)
		// 返回 JSON 响应，包含应用名称信息
		c.JSON(200, gin.H{
			"message": "golang-per-day-80-user:" + id,
		})
	})

	// 启动 HTTP 服务器，监听 8080 端口
	r.Run(":8080")
}
