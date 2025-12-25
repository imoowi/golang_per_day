package models

import (
	"golang_per_day_30/internal/interfaces"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 查询参数
type UserFilter struct {
	interfaces.Filter
}

func (f *UserFilter) BuildPageListFilter(c *gin.Context, db *gorm.DB) *gorm.DB {
	if f.GetSearchKey() != `` {
		db = db.Where(`name LIKE ?`, `%`+f.GetSearchKey()+`%`)
	}
	return f.Filter.BuildPageListFilter(c, db)
}
