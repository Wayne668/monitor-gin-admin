package api

import (
	"strconv"

	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// AgentToken 账户token API
type AgentToken struct {
	AgentTokenBIZ *biz.AgentToken
}

// Query 查询账户token列表
func (a *AgentToken) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.AgentTokenQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.AgentTokenBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// Get 获取账户token详情
func (a *AgentToken) Get(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	item, err := a.AgentTokenBIZ.Get(ctx, uint(id))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// Create 新增账户token
func (a *AgentToken) Create(c *gin.Context) {
	ctx := c.Request.Context()
	form := new(schema.AgentTokenForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.AgentTokenBIZ.Create(ctx, form)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

// Update 更新账户token
func (a *AgentToken) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	form := new(schema.AgentTokenForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err = a.AgentTokenBIZ.Update(ctx, uint(id), form)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// Delete 删除账户token
func (a *AgentToken) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	err = a.AgentTokenBIZ.Delete(ctx, uint(id))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
