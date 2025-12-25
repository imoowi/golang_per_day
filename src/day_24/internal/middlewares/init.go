package middlewares

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func InitMiddlewares(r *gin.Engine) {
	// 全局中间件
	r.Use(RateLimitMiddleware())
	// 为每个请求生成唯一的请求ID
	r.Use(requestid.New())
	// 追踪ID中间件
	r.Use(TraceIdMiddleware())
	// 日志中间件
	r.Use(LoggerMiddleware())
}
