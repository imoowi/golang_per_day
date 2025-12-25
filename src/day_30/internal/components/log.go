package components

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitLog() {
	var logger *zap.Logger
	var err error
	if gin.Mode() == gin.DebugMode {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic("初始化Zap日志失败：" + err.Error())
	}
	zap.ReplaceGlobals(logger)
}
