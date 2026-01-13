// 包声明，定义包名为 main
package main

// 导入依赖包
import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
)

// Prometheus中间件
func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先执行后续逻辑

		// 请求完成后记录指标
		status := strconv.Itoa(c.Writer.Status())
		httpTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
	}
}

// 主函数，程序入口
func main() {
	fmt.Println("Hello, Codee君!")
	// 注册指标
	prometheus.MustRegister(httpTotal)
	// 创建默认的 Gin 路由引擎
	r := gin.Default()
	// 使用中间件，所有路由都会被监控
	r.Use(prometheusMiddleware())
	// 定义 GET 路由 /ping，处理函数返回 JSON 格式的响应
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 暴露 metrics 接口
	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/error", func(c *gin.Context) {
		c.JSON(500, gin.H{"error": "boom"})
	})
	// 启动 HTTP 服务器，监听 8080 端口
	r.Run(":8080")
}
