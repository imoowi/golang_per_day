package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 使用原子操作保证线程安全
	var isReady int32 = 0

	// 模拟异步初始化任务（如加载大型配置、预热缓存、建立 DB 连接）
	go func() {
		log.Println("正在初始化服务组件...")
		time.Sleep(10 * time.Second) // 模拟耗时操作
		atomic.StoreInt32(&isReady, 1)
		log.Println("服务初始化完成，可以接收流量！")
	}()

	// 健康检查路由组
	health := router.Group("/health")
	{
		// 1. Startup: 只要 Gin 跑起来了就返回 200
		health.GET("/startup", func(c *gin.Context) {
			c.String(http.StatusOK, "started")
		})

		// 2. Liveness: 检查进程是否还在正常循环，是否死锁
		health.GET("/live", func(c *gin.Context) {
			c.String(http.StatusOK, "alive")
		})

		// 3. Readiness: 核心逻辑！检查依赖是否就绪
		health.GET("/ready", func(c *gin.Context) {
			if atomic.LoadInt32(&isReady) == 1 {
				c.String(http.StatusOK, "ready")
			} else {
				// 关键：返回 503，K8s 会将此 Pod 从 Service 摘除
				c.String(http.StatusServiceUnavailable, "initializing")
			}
		})
	}

	// --- 优雅关闭 (Graceful Shutdown) ---
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	// 技巧：收到关闭信号后，可以立即将 isReady 设为 0
	// 这样在 Pod 还没被彻底杀掉前，Readiness 探针会先失败，确保新流量不再进来
	atomic.StoreInt32(&isReady, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}
	log.Println("服务器已退出")
}
