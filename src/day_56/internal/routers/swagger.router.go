package routers

import (
	_ "codee_jun/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	RegisterRoute(SwaggerRouters)
}

func SwaggerRouters(e *gin.Engine) {
	// swagger
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
