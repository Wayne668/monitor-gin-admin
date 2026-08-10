package api

import (
	"strconv"

	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// HostAccount 托管账户 API
type HostAccount struct {
	HostAccountBIZ *biz.HostAccount
}

// Query 查询托管账户列表
func (a *HostAccount) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.HostAccountQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.HostAccountBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// Get 获取托管账户详情
func (a *HostAccount) Get(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	item, err := a.HostAccountBIZ.Get(ctx, uint(id))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// Create 新增托管账户
func (a *HostAccount) Create(c *gin.Context) {
	ctx := c.Request.Context()
	form := new(schema.HostAccountForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.HostAccountBIZ.Create(ctx, form)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

// Update 更新托管账户
func (a *HostAccount) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	form := new(schema.HostAccountForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err = a.HostAccountBIZ.Update(ctx, uint(id), form)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// Delete 软删除托管账户
func (a *HostAccount) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	err = a.HostAccountBIZ.Delete(ctx, uint(id))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// ListByAgent 根据代理商ID获取账户列表
func (a *HostAccount) ListByAgent(c *gin.Context) {
	ctx := c.Request.Context()
	agentIDStr := c.Query("agentId")
	if agentIDStr == "" {
		util.ResError(c, errors.BadRequest("", "agentId is required"))
		return
	}
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid agentId"))
		return
	}

	result, err := a.HostAccountBIZ.Query(ctx, schema.HostAccountQueryParam{AgentID: agentID})
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result.Data)
}