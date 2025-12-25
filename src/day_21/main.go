package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"gindemo2/config"
	"gindemo2/models"
	"gindemo2/routers"
)

func main() {
	// 加载env变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using environment variables")
	}

	// 连接数据库
	config.ConnectDatabase()
	// 连接Redis
	config.ConnectRedis()
	// 自动迁移
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	// 初始化路由
	r := routers.InitRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
