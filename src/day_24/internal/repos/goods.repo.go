package repos

import (
	"golang_per_day_24/internal/components"
	"golang_per_day_24/internal/interfaces"
	"golang_per_day_24/internal/models"
)

var Goods *GoodsRepo

type GoodsRepo struct {
	interfaces.Repo[*models.Goods]
}

func NewGoodsRepo() {
	Goods = &GoodsRepo{
		Repo: *interfaces.NewRepo[*models.Goods](components.DB),
	}
}

func init() {
	RegisterRepos(NewGoodsRepo)
}

func (r *GoodsRepo) GetById(id uint) (res models.Goods, err error) {
	err = r.Repo.DB.Where(`id=?`, id).First(&res).Error
	if err != nil {
		return
	}
	return
}
