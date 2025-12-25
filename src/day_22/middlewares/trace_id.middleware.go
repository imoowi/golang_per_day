package middlewares

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func TraceIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(`trace_id`, requestid.Get(c))
		c.Next()
	}
}
