package api

import (
	"strconv"

	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// AdvertiserBudget 广告主预算 API
type AdvertiserBudget struct {
	OceanengineBIZ         *biz.Oceanengine
	UpdateAccountBudgetDAL *dal.UpdateAccountBudget
}

// UpdateImmediate 立即生效：调用 Oceanengine API 更新预算，并记录到 nb_update_account_budget
func (a *AdvertiserBudget) UpdateImmediate(c *gin.Context) {
	ctx := c.Request.Context()
	form := new(schema.AdvertiserBudgetUpdateForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	requestID, err := a.OceanengineBIZ.UpdateAdvertiserBudget(ctx, form.AccountID, form.AdvertiserID, form.BudgetMode, form.Budget)
	if err != nil {
		util.ResError(c, err)
		return
	}

	// 记录到 nb_update_account_budget
	record := &schema.UpdateAccountBudget{
		AdvertiserID: form.AdvertiserID,
		Budget:       form.Budget,
		BudgetMode:    "immediate",
		IsSet:        1,
	}
	if err := a.UpdateAccountBudgetDAL.Create(ctx, record); err != nil {
		util.ResError(c, err)
		return
	}

	util.ResSuccess(c, gin.H{"requestId": requestID})
}

// ScheduleNextDay 次日生效：保存预算记录到 nb_update_account_budget 表
func (a *AdvertiserBudget) ScheduleNextDay(c *gin.Context) {
	ctx := c.Request.Context()
	form := new(schema.UpdateAccountBudgetForm)
	if err := util.ParseJSON(c, form); err != nil {
		util.ResError(c, err)
		return
	} else if err := form.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	item := &schema.UpdateAccountBudget{}
	form.FillTo(item)
	if item.BudgetMode == "" {
		item.BudgetMode = "nextDay"
	}

	if err := a.UpdateAccountBudgetDAL.Create(ctx, item); err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// QueryByAdvertiser 根据广告主ID查询预算修改记录
func (a *AdvertiserBudget) QueryByAdvertiser(c *gin.Context) {
	ctx := c.Request.Context()
	advertiserIDStr := c.Query("advertiserId")
	if advertiserIDStr == "" {
		util.ResError(c, errors.BadRequest("", "advertiserId is required"))
		return
	}
	advertiserID, err := strconv.ParseInt(advertiserIDStr, 10, 64)
	if err != nil {
		util.ResError(c, errors.BadRequest("", "invalid advertiserId"))
		return
	}

	list, err := a.UpdateAccountBudgetDAL.QueryByAdvertiserID(ctx, advertiserID)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, list)
}
