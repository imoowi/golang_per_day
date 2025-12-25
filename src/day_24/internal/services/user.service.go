package services

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_per_day_24/internal/components"
	"golang_per_day_24/internal/interfaces"
	"golang_per_day_24/internal/models"
	"golang_per_day_24/internal/repos"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var User *UserService

type UserService struct {
	redis *redis.Client
	interfaces.Service[*models.User]
}

func NewUserService(r *repos.UserRepo, redisClient *redis.Client) *UserService {
	return &UserService{
		redis:   redisClient,
		Service: *interfaces.NewService(repos.User),
	}
}

// 初始化用户服务

func init() {
	RegisterServices(func() {
		User = NewUserService(repos.User, components.Redis)
	})
}
func (s *UserService) GetById(c *gin.Context) (user *models.User, err error) {
	// 获取上下文传过来的用户id
	uid := c.GetUint(`user_id`)
	// 组合redis的key
	key := fmt.Sprintf("user:%d", uid)
	// 从redis里取数据
	data, err := components.Redis.Get(context.Background(), key).Result()
	if err == nil {
		// 取到了，反序列化到user里
		json.Unmarshal([]byte(data), &user)
		return
	}
	// 没取到，从repo里取数据
	user, err = repos.User.One(c, uid)
	if err != nil {
		fields := []zap.Field{
			zap.Uint("user_id", uid),
			zap.String("trace_id", c.GetString(`trace_id`)),
		}
		zap.L().Warn("repo find user failed", fields...)
		err = fmt.Errorf("service GetById:%w", err)
		return
	}
	// 取到了，序列化数据
	btdata, _ := json.Marshal(user)
	// 放入redis里
	components.Redis.Set(context.Background(), key, btdata, time.Hour)
	return
}
