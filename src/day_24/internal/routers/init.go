package routers

import (
	"golang_per_day_24/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// 定义路由函数
type Router func(*gin.Engine)

// 这里放所有的路由
var routers = []Router{}

// 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	middlewares.InitMiddlewares(r)
	for _, route := range routers {
		route(r)
	}
	return r
}

// 注册路由通用函数
func RegisterRoute(r ...Router) {
	routers = append(routers, r...)
}
