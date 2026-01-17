package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"log/slog" // 使用 Go 官方的 slog

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

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
func main() {
	// 初始化 OpenTelemetry Provider，并设置清理函数
	cleanup := initProvider()
	defer cleanup(context.Background())
	// 注册指标
	prometheus.MustRegister(httpTotal)
	// 创建 Gin 路由器，默认包含 Logger 和 Recovery 中间件
	r := gin.Default()
	// 使用中间件，所有路由都会被监控
	r.Use(prometheusMiddleware())
	// 使用 otelgin 中间件，这会自动为每个请求生成 Trace
	r.Use(otelgin.Middleware("my-server"))
	// 暴露 metrics 接口
	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	// 定义 GET /user/:id 路由
	r.GET("/user/:id", func(c *gin.Context) {
		// 1. 获取当前请求的 Span（由 otelgin 中间件创建）
		span := trace.SpanFromContext(c.Request.Context())
		traceID := span.SpanContext().TraceID().String()

		// 2. Metrics 埋点示例：使用 OpenTelemetry Metrics API
		meter := otel.Meter("gin-server")
		counter, _ := meter.Int64Counter("api_requests_total")
		// 增加计数器，并添加 path 属性
		counter.Add(c.Request.Context(), 1, metric.WithAttributes(attribute.String("path", "/user")))

		// 3. Logs 实战：手动注入 TraceID
		// 在实际项目中，建议写一个 slog 的自定义 Handler 自动提取
		logger := slog.With("trace_id", traceID)
		logger.Info("开始处理用户请求", "user_id", c.Param("id"))

		// 模拟业务逻辑处理
		time.Sleep(100 * time.Millisecond)

		// 如果用户 ID 为 "error"，返回 404 错误
		if c.Param("id") == "error" {
			// 设置 span 属性，标记为错误
			span.SetAttributes(attribute.Bool("error", true))
			logger.Error("用户不存在")
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// 返回成功响应，包含 trace_id
		c.JSON(http.StatusOK, gin.H{"status": "ok", "trace_id": traceID})
	})

	// 启动 HTTP 服务器，监听 8080 端口
	r.Run(":8080")
}
