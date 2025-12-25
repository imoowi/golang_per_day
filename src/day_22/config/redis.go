package config

import (
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

var Redis *redis.Client

func ConnectRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		panic("请在环境变量里配置【REDIS_ADDR】")
	}
	pass := os.Getenv("REDIS_PASS")
	if pass == "" {
		panic("请在环境变量里配置【REDIS_PASS】")
	}
	db := os.Getenv("REDIS_DB")
	if db == "" {
		panic("请在环境变量里配置【REDIS_DB】")
	}
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       cast.ToInt(db),
		PoolSize: 10,
	})
}
