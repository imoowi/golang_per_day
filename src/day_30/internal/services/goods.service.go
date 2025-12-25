package services

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_per_day_30/internal/components"
	"golang_per_day_30/internal/interfaces"
	"golang_per_day_30/internal/models"
	"golang_per_day_30/internal/repos"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Goods *GoodsService

type GoodsService struct {
	redis *redis.Client
	interfaces.Service[*models.Goods]
}

func NewGoodsService(r *repos.GoodsRepo, redisClient *redis.Client) *GoodsService {
	return &GoodsService{
		redis:   components.Redis,
		Service: *interfaces.NewService(repos.Goods),
	}
}

func init() {
	RegisterServices(func() {
		Goods = NewGoodsService(repos.Goods, components.Redis)
	})
}
func (s *GoodsService) GetById(c *gin.Context) (goods *models.Goods, err error) {
	// 获取上下文传过来的用户id
	uid := c.GetUint(`goods_id`)
	// 组合redis的key
	key := fmt.Sprintf("goods:%d", uid)
	// 从redis里取数据
	data, err := s.redis.Get(context.Background(), key).Result()
	if err == nil {
		// 取到了，反序列化到goods里
		json.Unmarshal([]byte(data), &goods)
		return
	}
	// 没取到，从repo里取数据
	goods, err = s.One(c, uid)
	if err != nil {
		fields := []zap.Field{
			zap.Uint("goods_id", uid),
			zap.String("trace_id", c.GetString(`trace_id`)),
		}
		zap.L().Warn("repo find goods failed", fields...)
		err = fmt.Errorf("service GetById:%w", err)
		return
	}
	// 取到了，序列化数据
	btdata, _ := json.Marshal(goods)
	// 放入redis里
	s.redis.Set(context.Background(), key, btdata, time.Hour)
	return
}
