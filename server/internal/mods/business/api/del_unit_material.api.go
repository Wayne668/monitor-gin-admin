package api

import (
	"strconv"
	"strings"

	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// DeleteUnauditedMaterial 素材删除记录 API
type DeleteUnauditedMaterial struct {
	DeleteUnauditedMaterialBIZ *biz.DeleteUnauditedMaterial
}

// Query 查询素材删除记录列表
func (a *DeleteUnauditedMaterial) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.DeleteUnauditedMaterialQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.DeleteUnauditedMaterialBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// Get 获取素材删除记录详情
func (a *DeleteUnauditedMaterial) Get(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid id"))
		return
	}

	item, err := a.DeleteUnauditedMaterialBIZ.Get(ctx, uint(id))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// GetUnAudititedMaterial 获取未审核素材列表
func (a *DeleteUnauditedMaterial) GetUnAudititedMaterial(c *gin.Context) {
	ctx := c.Request.Context()
	accountIDsStr := c.Query("accountIds")
	if accountIDsStr == "" {
		util.ResError(c, errors.BadRequest("", "accountIds is required"))
		return
	}

	parts := strings.Split(accountIDsStr, ",")
	accountIDs := make([]int64, 0, len(parts))
	for _, p := range parts {
		id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			util.ResError(c, errors.BadRequest("", "invalid accountIds"))
			return
		}
		accountIDs = append(accountIDs, id)
	}

	items, err := a.DeleteUnauditedMaterialBIZ.GetUnAudititedMaterialWithFallback(ctx, 0, accountIDs)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, items)
}

// DeleteUnAudititedMaterial 删除未审核素材
func (a *DeleteUnauditedMaterial) DeleteUnAudititedMaterial(c *gin.Context) {
	ctx := c.Request.Context()
	req := new(schema.UnAudititedMaterialReq)
	if err := util.ParseJSON(c, req); err != nil {
		util.ResError(c, err)
		return
	}

	if len(req.Materials) == 0 {
		util.ResError(c, errors.BadRequest("", "materials is required"))
		return
	}

	failedRecords, err := a.DeleteUnauditedMaterialBIZ.DeleteAndSave(ctx, req)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, map[string]interface{}{
		"failed": failedRecords,
	})
}

// RetryFailedDelete 重试删除失败的记录
func (a *DeleteUnauditedMaterial) RetryFailedDelete(c *gin.Context) {
	ctx := c.Request.Context()
	accountIDStr := c.Query("accountId")
	if accountIDStr == "" {
		util.ResError(c, errors.BadRequest("", "accountId is required"))
		return
	}
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid accountId"))
		return
	}

	err = a.DeleteUnauditedMaterialBIZ.RetryFailedDelete(ctx, accountID)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
