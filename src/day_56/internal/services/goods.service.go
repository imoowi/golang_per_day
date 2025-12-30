package services

import (
	"codee_jun/internal/components"
	"codee_jun/internal/interfaces"
	"codee_jun/internal/models"
	"codee_jun/internal/repos"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
func (s *GoodsService) GetById(c *gin.Context, id uint) (*models.Goods, error) {
	return repos.Goods.GetById(c, id)
}
func (s *GoodsService) UpdateById(c *gin.Context, id uint, data map[string]interface{}) (err error) {
	return repos.Goods.UpdateById(c, id, data)
}

// AddOne 添加一个商品
func (s *GoodsService) AddOne(c *gin.Context, model *models.Goods) (newId uint, err error) {
	return repos.Goods.AddOne(c, model)
}
