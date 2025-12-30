package controllers

import (
	"codee_jun/internal/interfaces"
	"codee_jun/internal/models"
	"codee_jun/internal/services"
	"codee_jun/internal/utils/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// @Summary	分页列表(pagelist)
// @Tags	goods
// @Accept	application/json
// @Produce	application/json
// @Param	Authorization	header		string				true	"Bearer 用户令牌"
// @Param   {object}             query                 models.GoodsFilter  false "query参数"
// @Success 200                        {object} response.PageListT[models.Goods] "成功"
// @Failure 400                        "请求错误"
// @Failure 401                        "token验证失败"
// @Failure 500                         "内部错误"
// @Router		/api/v1/goods [get]
func GoodsPageList(c *gin.Context) {
	var filter interfaces.IFilter = &models.GoodsFilter{}
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}

	searchKey := c.DefaultQuery(`searchKey`, ``)
	filter.SetSearchKey(searchKey)
	page := c.DefaultQuery(`page`, `1`)
	filter.SetPage(cast.ToInt64(page))
	pageSize := c.DefaultQuery(`pageSize`, `20`)
	filter.SetPageSize(cast.ToInt64(pageSize))
	if 0 >= filter.GetPage() { //如果不传Page，默认为1
		filter.SetPage(1)
	}
	if 0 >= filter.GetPageSize() { //如果不传PageSize，默认取20条
		filter.SetPageSize(20)
	}
	if filter.GetPageSize() > 1000 {
		response.Error(c, `每一页不能超过1000条记录`, http.StatusBadRequest)
		return
	}
	// nType := c.DefaultQuery(`type`, `0`)
	// filter.SetType(nType)
	result, err := services.Goods.PageList(c, &filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, err.Error(), http.StatusNotFound)
			return
		}
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.OK(c, result)
}

// @Summary	详情(one)
// @Tags		goods
// @Accept		application/json
// @Produce	application/json
// @Param		Authorization	header	string	true	"Bearer 用户令牌"
// @Param		id				path	int		true	"id"
// @Success	200 						{object} models.Goods "成功"
// @Failure 400                        "请求错误"
// @Failure 401                        "token验证失败"
// @Failure 500                         "内部错误"
// @Router		/api/v1/goods/{id} [get]
func GoodsOne(c *gin.Context) {
	//参数校验，防止穿透
	id := cast.ToUint(c.Param(`id`))
	if id <= 0 {
		response.Error(c, `pls input id`, http.StatusBadRequest)
		return
	}

	one, err := services.Goods.GetById(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, err.Error(), http.StatusNotFound)
			return
		}
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.OK(c, one)
}

// @Summary	新增(add)
// @Tags		goods
// @Accept		application/json
// @Produce	application/json
// @Param		Authorization	header	string	true	"Bearer 用户令牌"
// @Param 		{object} body  models.Goods true "body"
// @Success	200
// @Failure 400                        "请求错误"
// @Failure 401                        "token验证失败"
// @Failure 500                         "内部错误"
// @Router		/api/v1/goods [post]
func GoodsAdd(c *gin.Context) {
	model := &models.Goods{}
	err := c.ShouldBindBodyWith(&model, binding.JSON)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	newId, err := services.Goods.AddOne(c, model)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.OK(c, newId)
}

// @Summary	更新(update)
// @Tags		goods
// @Accept		application/json
// @Produce	application/json
// @Param		Authorization	header	string	true	"Bearer 用户令牌"
// @Param		id				path	int		true	"id"
// @Param 		{object} body  models.Goods true "body"
// @Success	200
// @Failure 400                        "请求错误"
// @Failure 401                        "token验证失败"
// @Failure 500                         "内部错误"
// @Router		/api/v1/goods/{id} [put]
func GoodsUpdate(c *gin.Context) {
	//参数校验，防止穿透
	id := cast.ToUint(c.Param(`id`))
	if id <= 0 {
		response.Error(c, `pls input id`, http.StatusBadRequest)
		return
	}
	model := make(map[string]any)
	err := c.ShouldBindBodyWith(&model, binding.JSON)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	delete(model, `updated_at`)
	err = services.Goods.UpdateById(c, id, model)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, err.Error(), http.StatusNotFound)
			return
		}
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.OK(c, true)
}

// @Summary	删除(delete)
// @Tags		goods
// @Accept		application/json
// @Produce	application/json
// @Param		Authorization	header	string	true	"Bearer 用户令牌"
// @Param		id				path	int		true	"id"
// @Success	200
// @Failure 400                        "请求错误"
// @Failure 401                        "token验证失败"
// @Failure 500                         "内部错误"
// @Router		/api/v1/goods/{id} [delete]
func GoodsDel(c *gin.Context) {
	id := c.Param(`id`)
	if id == `` {
		response.Error(c, `pls input id`, http.StatusBadRequest)
		return
	}
	deleted, err := services.Goods.Delete(c, cast.ToUint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, err.Error(), http.StatusNotFound)
			return
		}
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.OK(c, deleted)
}
