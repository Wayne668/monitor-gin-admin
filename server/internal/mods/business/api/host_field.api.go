package api

import (
	"strconv"

	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// HostField 托管字段API
type HostField struct {
	HostFieldBIZ *biz.HostField
}

// Query 查询托管字段列表
func (a *HostField) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.HostFieldQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.HostFieldBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// FindAll 查询全部托管字段（不分页，前端渲染用）
func (a *HostField) FindAll(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := a.HostFieldBIZ.FindAll(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, list)
}

// Get 获取托管字段详情（id 为 "all" 时返回全部字段）
func (a *HostField) Get(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	if idStr == "all" {
		list, err := a.HostFieldBIZ.FindAll(ctx)
		if err != nil {
			util.ResError(c, err)
			return
		}
		util.ResSuccess(c, list)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	item, err := a.HostFieldBIZ.Get(ctx, uint(id))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// Create 新增托管字段
func (a *HostField) Create(c *gin.Context) {
	ctx := c.Request.Context()
	form := new(schema.HostFieldForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.HostFieldBIZ.Create(ctx, form)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

// Update 更新托管字段
func (a *HostField) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	form := new(schema.HostFieldForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err = a.HostFieldBIZ.Update(ctx, uint(id), form)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// Delete 删除托管字段
func (a *HostField) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	err = a.HostFieldBIZ.Delete(ctx, uint(id))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
