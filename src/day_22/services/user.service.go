package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gindemo2/config"
	"gindemo2/models"
	"gindemo2/repos"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	User *UserService
	Ctx  context.Context = context.Background()
)

type UserService struct {
}

// 初始化用户服务
func init() {
	User = &UserService{}
}
func (s *UserService) GetById(c *gin.Context) (user models.User, err error) {
	// 获取上下文传过来的用户id
	uid := c.GetUint(`user_id`)
	// 组合redis的key
	key := fmt.Sprintf("user:%d", uid)
	// 从redis里取数据
	data, err := config.Redis.Get(Ctx, key).Result()
	if err == nil {
		// 取到了，反序列化到user里
		json.Unmarshal([]byte(data), &user)
		return
	}
	// 没取到，从repo里取数据
	user, err = repos.User.GetById(uid)
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
	config.Redis.Set(Ctx, key, btdata, time.Hour)
	return
}
