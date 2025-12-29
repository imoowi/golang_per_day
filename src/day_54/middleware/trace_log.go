package middleware

import (
	"day54/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// TraceLog 返回一个 Gin 中间件，用于为 HTTP 请求添加跟踪日志
// 该中间件会创建 OpenTelemetry 跟踪 span，并记录请求的关键信息
func TraceLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 HTTP 跟踪器
		tracer := otel.Tracer("http")
		// 为当前请求创建一个 span，使用请求的完整路径作为 span 名称
		ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
		// 确保请求结束后关闭 span
		defer span.End()

		// 获取跟踪 ID
		traceID := span.SpanContext().TraceID().String()
		// 将跟踪 ID 存储到 Gin 上下文中，以便后续处理函数使用
		c.Set("trace_id", traceID)
		// 更新请求的上下文，包含跟踪信息
		c.Request = c.Request.WithContext(ctx)

		// 记录请求开始时间
		start := time.Now()
		// 调用后续中间件和处理函数
		c.Next()

		// 记录请求日志，包含跟踪 ID、路径、状态码和处理时间
		logger.Log.Info("http request",
			zap.String("trace_id", traceID),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("cost", time.Since(start)),
		)
	}
}
