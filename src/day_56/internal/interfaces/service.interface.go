package interfaces

import (
	"codee_jun/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// 服务层接口
type IService[T IModel] interface {

	// 分页查询
	PageList(c *gin.Context, filter *IFilter) (res *response.PageListT[T], err error)
	// 分页查询
	PageListWithSelectOption(c *gin.Context, filter *IFilter, selectOpt []string) (res *response.PageListT[T], err error)
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

// 服务层接口实现
type Service[T IModel] struct {
	Repo *IRepo[T]
}

// 新建服务
func NewService[T IModel](r IRepo[T]) *Service[T] {
	return &Service[T]{
		Repo: &r,
	}
}

// 分页查询
func (s *Service[T]) PageList(c *gin.Context, filter *IFilter) (res *response.PageListT[T], err error) {
	repo := *s.Repo
	return repo.PageList(c, filter)
}

// 分页查询
func (s *Service[T]) PageListWithSelectOption(c *gin.Context, filter *IFilter, selectOpt []string) (res *response.PageListT[T], err error) {
	repo := *s.Repo
	return repo.PageListWithSelectOption(c, filter, selectOpt)
}

// 查一条，根据id
func (s *Service[T]) One(c *gin.Context, id uint) (res T, err error) {
	repo := *s.Repo
	return repo.One(c, id)
}

// 查一条，根据id
func (s *Service[T]) OneWithSelectOption(c *gin.Context, id uint, selectOpt []string) (res T, err error) {
	repo := *s.Repo
	return repo.OneWithSelectOption(c, id, selectOpt)
}

// 查一条，根据名字
func (s *Service[T]) OneByName(c *gin.Context, name string) (res T, err error) {
	repo := *s.Repo
	return repo.OneByName(c, name)
}

// 查一条，根据名字
func (s *Service[T]) OneByNameWithSelectOption(c *gin.Context, name string, selectOpt []string) (res T, err error) {
	repo := *s.Repo
	return repo.OneByNameWithSelectOption(c, name, selectOpt)
}

// 新建资源
func (s *Service[T]) Add(c *gin.Context, model T) (newId uint, err error) {
	repo := *s.Repo
	return repo.Add(c, model)
}

// 更新资源
func (s *Service[T]) Update(c *gin.Context, updateFields map[string]any, id uint) (updated bool, err error) {
	delete(updateFields, `created_at`)
	delete(updateFields, `updated_at`)
	delete(updateFields, `deleted_at`)
	repo := *s.Repo
	return repo.Update(c, updateFields, id)
}

// 删除资源，根据id
func (s *Service[T]) Delete(c *gin.Context, id uint) (deleted bool, err error) {
	repo := *s.Repo
	return repo.Delete(c, id)
}
