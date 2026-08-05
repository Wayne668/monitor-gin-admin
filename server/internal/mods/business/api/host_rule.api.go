package api

import (
	"strconv"

	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// HostRule 托管规则API
type HostRule struct {
	HostRuleBIZ *biz.HostRule
}

// Query 查询托管规则列表
func (a *HostRule) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.HostRuleQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.HostRuleBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// Get 获取托管规则详情
func (a *HostRule) Get(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	item, err := a.HostRuleBIZ.Get(ctx, uint(id))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// UpdateStatus 修改托管规则状态
func (a *HostRule) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	form := new(schema.HostRuleUpdateStatusForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err = a.HostRuleBIZ.UpdateStatus(ctx, uint(id), form)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// GetTargetsByAccount 根据账户获取目标列表
func (a *HostRule) GetTargetsByAccount(c *gin.Context) {
	ctx := c.Request.Context()
	req := new(schema.TargetByAccountReq)
	if err := util.ParseJSON(c, req); err != nil {
		util.ResError(c, err)
		return
	} else if err := req.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	items, err := a.HostRuleBIZ.GetTargetsByAccount(ctx, req)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, items)
}
