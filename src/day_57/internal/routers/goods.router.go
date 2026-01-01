package routers

import (
	"codee_jun/internal/controllers"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterRoute(func(e *gin.Engine) {
		// 路由分组
		v1 := e.Group("/api/v1")
		// v1.Use(middlewares.AuthMiddleware())
		goods := v1.Group("/goods")
		{
			goods.GET("", controllers.GoodsPageList)   //分页
			goods.GET("/:id", controllers.GoodsOne)    //一个
			goods.POST("", controllers.GoodsAdd)       //新增
			goods.PUT("/:id", controllers.GoodsUpdate) //更新
			goods.DELETE("/:id", controllers.GoodsDel) //删除
		}
	})
}
