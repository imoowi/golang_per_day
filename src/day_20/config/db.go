package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// dsn := "root:123456@tcp(127.0.0.1:3306)/golang_per_day??charset=utf8&parseTime=True&loc=Local&timeout=1000ms"
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		panic("请在环境变量里配置【DATABASE_DSN】")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接接数据库失败：", err)
	}
	log.Println("数据库连接成功:", db)
	sqlDb, err := db.DB()
	if err == nil {
		sqlDb.SetMaxIdleConns(10)
		sqlDb.SetMaxOpenConns(100)
		sqlDb.SetConnMaxLifetime(time.Hour)
	}
	// 开启debug模式
	DB = db.Debug()

}
