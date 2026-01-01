package models

import (
	"codee_jun/internal/interfaces"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 查询参数
type OrderFilter struct {
	interfaces.Filter
}

func (f *OrderFilter) BuildPageListFilter(c *gin.Context, db *gorm.DB) *gorm.DB {
	if f.GetSearchKey() != `` {
		db = db.Where(`name LIKE ?`, `%`+f.GetSearchKey()+`%`)
	}
	return f.Filter.BuildPageListFilter(c, db)
}
