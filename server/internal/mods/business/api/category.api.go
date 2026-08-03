package api

import (
	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/util"
	"github.com/gin-gonic/gin"
)

// Category management for business
type Category struct {
	CategoryBIZ *biz.Category
}

// @Tags CategoryAPI
// @Security ApiKeyAuth
// @Summary Query category list
// @Param name query string false "Name of category"
// @Param status query string false "Status of category"
// @Param current query int false "Current page"
// @Param pageSize query int false "Page size"
// @Success 200 {object} util.ResponseResult{data=[]schema.Category}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/categories [get]
func (a *Category) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.CategoryQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.CategoryBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// @Tags CategoryAPI
// @Security ApiKeyAuth
// @Summary Get category record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult{data=schema.Category}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/categories/{id} [get]
func (a *Category) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.CategoryBIZ.Get(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// @Tags CategoryAPI
// @Security ApiKeyAuth
// @Summary Create category record
// @Param body body schema.CategoryForm true "Request body"
// @Success 200 {object} util.ResponseResult{data=schema.Category}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/categories [post]
func (a *Category) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.CategoryForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.CategoryBIZ.Create(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

// @Tags CategoryAPI
// @Security ApiKeyAuth
// @Summary Update category record by ID
// @Param id path string true "unique id"
// @Param body body schema.CategoryForm true "Request body"
// @Success 200 {object} util.ResponseResult
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/categories/{id} [put]
func (a *Category) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.CategoryForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.CategoryBIZ.Update(ctx, c.Param("id"), item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// @Tags CategoryAPI
// @Security ApiKeyAuth
// @Summary Delete category record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/categories/{id} [delete]
func (a *Category) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.CategoryBIZ.Delete(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
