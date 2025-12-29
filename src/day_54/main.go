package main

import (
	"day54/logger"
	"day54/middleware"
	"day54/tracing"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// main 函数是应用程序的入口点
func main() {
	// 初始化日志系统
	logger.Init()
	// 初始化跟踪系统，服务名称为 "log-tracer"
	tracing.Init("log-tracer")

	// 创建 Gin 引擎实例（不使用默认中间件）
	r := gin.New()
	// 添加恢复中间件，防止 panic 导致服务崩溃
	r.Use(gin.Recovery())
	// 添加跟踪日志中间件，为每个请求添加跟踪信息
	r.Use(middleware.TraceLog())

	// 定义 /ping 路由的 GET 请求处理函数
	r.GET("/ping", func(c *gin.Context) {
		// 从 Gin 上下文中获取跟踪 ID
		traceID, _ := c.Get("trace_id")

		// 记录 ping 请求日志
		logger.Log.Info("ping called",
			zap.String("trace_id", traceID.(string)),
		)

		// 返回 JSON 响应，包含 "pong" 消息和跟踪 ID
		c.JSON(http.StatusOK, gin.H{
			"msg":      "pong",
			"trace_id": traceID,
		})
	})

	// 启动 HTTP 服务器，监听 8080 端口
	r.Run(":8080")
}
