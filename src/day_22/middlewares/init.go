package middlewares

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func InitMiddlewares(r *gin.Engine) {
	r.Use(requestid.New())
	r.Use(TraceIdMiddleware())
	r.Use(LoggerMiddleware())
}
