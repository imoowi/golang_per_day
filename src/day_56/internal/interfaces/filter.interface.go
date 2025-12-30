package interfaces

import (
	"codee_jun/internal/utils/request"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 过滤器接口
type IFilter interface {
	GetPage() int64
	SetPage(page int64)
	GetPageSize() int64
	SetPageSize(pageSize int64)
	GetSearchKey() string
	SetSearchKey(searchKey string)
	BuildPageListFilter(c *gin.Context, db *gorm.DB) *gorm.DB
}

// 过滤器接口实现
type Filter struct {
	request.PageList
}

func (f *Filter) GetPage() int64 {
	return f.Page
}

func (f *Filter) SetPage(page int64) {
	f.Page = page
}

func (f *Filter) GetPageSize() int64 {
	return f.PageSize
}
func (f *Filter) SetPageSize(pageSize int64) {
	f.PageSize = pageSize
}

func (f *Filter) GetSearchKey() string {
	return f.SearchKey
}
func (f *Filter) SetSearchKey(searchKey string) {
	f.SearchKey = searchKey
}

// 分页查询过滤器构建方法
func (f *Filter) BuildPageListFilter(c *gin.Context, db *gorm.DB) *gorm.DB {
	return db
}
