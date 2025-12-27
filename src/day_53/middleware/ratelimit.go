package middleware

import (
	"day53/limiter"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RateLimit 创建一个基于令牌桶的Gin限流中间件
// tb: 令牌桶限流器实例
func RateLimit(tb *limiter.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否允许请求通过限流器
		if !tb.Allow() {
			// 如果请求被限流，返回429 Too Many Requests状态码
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			// 终止请求处理链
			c.Abort()
			return
		}
		// 如果允许请求通过，继续处理下一个中间件
		c.Next()
	}
}
