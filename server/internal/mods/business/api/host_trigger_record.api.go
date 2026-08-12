package api

import (
	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// HostTriggerRecord 触发记录API
type HostTriggerRecord struct {
	HostTriggerRecordBIZ *biz.HostTriggerRecord
}

// Query 查询触发记录列表
func (a *HostTriggerRecord) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.HostTriggerRecordQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.HostTriggerRecordBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}