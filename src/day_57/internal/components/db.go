package components

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	dsn := viper.GetString("mysql.dsn")
	if dsn == "" {
		panic("请在配置文件里配置【mysql.dsn】")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("连接接数据库失败：", zap.Error(err))
		return
	}
	zap.L().Info("数据库连接成功:", zap.String("dsn", dsn))
	sqlDb, err := db.DB()
	if err == nil {
		sqlDb.SetMaxIdleConns(10)
		sqlDb.SetMaxOpenConns(100)
		sqlDb.SetConnMaxLifetime(time.Hour)
	}
	// 开启debug模式
	DB = db.Debug()
}