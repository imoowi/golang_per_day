package server

import (
	"context"
	"fmt"
	"golang_per_day_24/internal/components"
	"golang_per_day_24/internal/migrates"
	"golang_per_day_24/internal/routers"
	"golang_per_day_24/internal/services"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func StartHTTPServer(host string, port int) error {
	zap.L().Info("正在初始化服务器...")
	mode := viper.GetString("server.mode")
	switch mode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	components.Init()         // 初始化各个组件
	migrates.DoMigrate()      // 数据库迁移
	services.InitServices()   // 初始化服务层
	r := routers.InitRouter() // 初始化路由
	s := &http.Server{
		Addr:           fmt.Sprintf(`%s:%d`, host, port),
		Handler:        r,
		ReadTimeout:    time.Duration(viper.GetInt("server.readtimeout")) * time.Second,
		WriteTimeout:   time.Duration(viper.GetInt("server.writertimeout")) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// 启动HTTP服务
	go func() {
		zap.L().Info("HTTP服务器启动了", zap.String("addr", s.Addr))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("HTTP服务器启动失败", zap.Error(err))
		}

	}()

	//告诉k8s或者docker我准备好了
	services.SetReady()

	// 捕获系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	zap.L().Warn("收到关闭信号，正在优雅退出...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.SetKeepAlivesEnabled(false) // 避免心得连接进来
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("服务器被强制关闭了:", err)
	}
	zap.L().Info("服务器优雅退出了")
	return nil
}
