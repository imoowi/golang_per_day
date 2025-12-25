package interfaces

import (
	"errors"
	"golang_per_day_30/internal/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 数据资源接口
type IRepo[T IModel] interface {
	// 分页查询
	PageList(c *gin.Context, query *IFilter) (res *response.PageListT[T], err error)
	// 分页查询
	PageListWithSelectOption(c *gin.Context, query *IFilter, selectOpt []string) (res *response.PageListT[T], err error)
	// 查询一个
	One(c *gin.Context, id uint) (res T, err error)
	// 查询一个
	OneWithSelectOption(c *gin.Context, id uint, selectOpt []string) (res T, err error)
	// 根据名称查询
	OneByName(c *gin.Context, name string) (res T, err error)
	// 根据名称查询
	OneByNameWithSelectOption(c *gin.Context, name string, selectOpt []string) (res T, err error)
	// 添加
	Add(c *gin.Context, model T) (newId uint, err error)
	// 更新,传什么就更新什么
	Update(c *gin.Context, updateFields map[string]any, id uint) (updated bool, err error)
	// 删除
	Delete(c *gin.Context, id uint) (deleted bool, err error)
}

// 数据资源接口实现
type Repo[T IModel] struct {
	DB *gorm.DB
}

// 新建一个数据资源
func NewRepo[T IModel](db *gorm.DB) *Repo[T] {
	return &Repo[T]{
		DB: db,
	}
}

// 分页查询数据
func (r *Repo[T]) PageList(c *gin.Context, f *IFilter) (res *response.PageListT[T], err error) {
	db := r.DB
	db = (*f).BuildPageListFilter(c, db)
	offset := ((*f).GetPage() - 1) * (*f).GetPageSize()
	db = db.Model(new(T)).Offset(int(offset)).Limit(int((*f).GetPageSize()))
	objs := make([]T, 0)
	err = db.Find(&objs).Error
	var count int64
	db.Offset(-1).Limit(-1).Select("count(id)").Count(&count)

	res = &response.PageListT[T]{
		List:  objs,
		Pages: response.MakePages(count, (*f).GetPage(), (*f).GetPageSize()),
	}

	return
}

// 分页查询数据
func (r *Repo[T]) PageListWithSelectOption(c *gin.Context, f *IFilter, selectOpt []string) (res *response.PageListT[T], err error) {
	db := r.DB
	db = (*f).BuildPageListFilter(c, db)
	offset := ((*f).GetPage() - 1) * (*f).GetPageSize()
	db = db.Model(new(T)).Offset(int(offset)).Limit(int((*f).GetPageSize()))
	if len(selectOpt) > 0 {
		db = db.Select(selectOpt)
	}
	objs := make([]T, 0)
	err = db.Find(&objs).Error
	var count int64
	db.Offset(-1).Limit(-1).Select("count(id)").Count(&count)

	res = &response.PageListT[T]{
		List:  objs,
		Pages: response.MakePages(count, (*f).GetPage(), (*f).GetPageSize()),
	}

	return
}

// 根据id查询一条记录
func (r *Repo[T]) One(c *gin.Context, id uint) (res T, err error) {
	db := r.DB
	err = db.Model(new(T)).Where(`id=?`, id).First(&res).Error
	return
}

// 根据id查询一条记录
func (r *Repo[T]) OneWithSelectOption(c *gin.Context, id uint, selectOpt []string) (res T, err error) {
	db := r.DB
	db = db.Model(new(T)).Where(`id=?`, id)
	if len(selectOpt) > 0 {
		db = db.Select(selectOpt)
	}
	err = db.First(&res).Error
	return
}

// 根据名字查询一条记录
func (r *Repo[T]) OneByName(c *gin.Context, name string) (res T, err error) {
	db := r.DB
	err = db.Model(new(T)).Where(`name=?`, name).First(&res).Error
	return
}

// 根据名字查询一条记录
func (r *Repo[T]) OneByNameWithSelectOption(c *gin.Context, name string, selectOpt []string) (res T, err error) {
	db := r.DB
	db = db.Model(new(T)).Where(`name=?`, name)
	if len(selectOpt) > 0 {
		db = db.Select(selectOpt)
	}
	err = db.First(&res).Error
	return
}

// 新建资源
func (r *Repo[T]) Add(c *gin.Context, model T) (newId uint, err error) {
	db := r.DB
	err = db.Create(model).Error
	newId = model.GetId()
	return
}

// 通过id更新资源，只更新updateFields里有的字段
func (r *Repo[T]) Update(c *gin.Context, updateFields map[string]any, id uint) (updated bool, err error) {
	if id <= 0 {
		updated = false
		err = errors.New(`pls input id`)
		return
	}
	_, err = r.One(c, id)
	if err != nil {
		return
	}
	db := r.DB
	err = db.Model(new(T)).Omit(`created_at`).Where(`id=?`, id).Updates(updateFields).Error
	if err == nil {
		updated = true
	}
	return
}

// 根据id删除资源
func (r *Repo[T]) Delete(c *gin.Context, id uint) (deleted bool, err error) {
	if id <= 0 {
		deleted = false
		err = errors.New(`pls input id`)
		return
	}
	db := r.DB
	model, err := r.One(c, id)
	if err != nil {
		return
	}
	err = db.Model(new(T)).Where(`id=?`, id).Delete(&model).Error
	if err == nil {
		deleted = true
	}
	return
}
