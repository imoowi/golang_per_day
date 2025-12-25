package routers

import (
	"golang_per_day_24/internal/controllers"
	"golang_per_day_24/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// init()会在程序启动的时候自动运行
func init() {
	RegisterRoute(func(e *gin.Engine) {
		// 路由分组
		api := e.Group("/api/v1")
		// 路由分组
		auth := api.Group("/auth")

		auth.POST("/login", controllers.Login)
		auth.POST("/reg", controllers.Register)
		// 路由/api/v1/auth开头的都加上认证
		auth.Use(middlewares.AuthMiddleware())
		{
			auth.GET("/me", controllers.Me)
			//其他路由
		}
	})
}
