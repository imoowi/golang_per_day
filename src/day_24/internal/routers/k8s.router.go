package routers

import (
	"golang_per_day_24/internal/services"

	"github.com/gin-gonic/gin"
)

// 支持k8s的健康检查接口
func init() {
	RegisterRoute(func(e *gin.Engine) {
		e.GET("/healthz", func(ctx *gin.Context) {
			if services.IsReady() {
				ctx.JSON(200, gin.H{"status": "ok"})
			} else {
				ctx.JSON(503, gin.H{"status": "not ready"})
			}
		})
	})
}
