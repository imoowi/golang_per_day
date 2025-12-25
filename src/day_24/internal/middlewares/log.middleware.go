package middlewares

import (
	"math"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next() // 调用该请求的剩余处理程序
		stopTime := time.Since(startTime)
		spendTime := int(math.Ceil(float64(stopTime.Nanoseconds() / 1000000)))

		hostName, err := os.Hostname()
		if err != nil {
			hostName = "Unknown"
		}

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		url := c.Request.RequestURI
		raw := c.Request.URL.RawQuery
		// 构建 Zap 字段
		fields := []zap.Field{
			zap.String("hostName", hostName),
			zap.Int("spendTime", spendTime),
			zap.String("path", url),
			zap.String("query", raw),
			zap.String("method", method),
			zap.Int("status", statusCode),
			zap.String("clientIp", clientIP),
			zap.Int("dataSize", dataSize),
			zap.String("userAgent", userAgent),
			zap.String("trace_id", c.GetString(`trace_id`)),
		}

		// 添加错误信息（如果有）
		if len(c.Errors) > 0 {
			fields = append(fields, zap.Strings("errors", c.Errors.Errors()))
		}

		// 根据状态码级别记录日志
		switch {
		case statusCode >= 500:
			// L returns the global Logger, which can be reconfigured with ReplaceGlobals. It's safe for concurrent use.
			// zap.L()返回zap.ReplaceGlobals(logger)设置的全局logger，它是并发安全的
			zap.L().Error("HTTP Server Error", fields...)
		case statusCode >= 400:
			zap.L().Warn("HTTP Client Error", fields...)
		default:
			zap.L().Info("HTTP Request", fields...)
		}
	}
}
