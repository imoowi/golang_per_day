// 包声明，定义包名为 main
package main

// 导入依赖包
import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 主函数，程序入口
func main() {
	fmt.Println("Hello, Codee君!")
	// 创建默认的 Gin 路由引擎
	r := gin.Default()
	// 定义 GET 路由 /ping，处理函数返回 JSON 格式的响应
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 启动 HTTP 服务器，监听 8080 端口
	r.Run(":8080")
}
