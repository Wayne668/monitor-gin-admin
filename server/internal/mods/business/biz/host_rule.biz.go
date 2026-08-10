package biz

import (
	"context"
	"strconv"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"
)

// HostRule 托管规则业务逻辑
type HostRule struct {
	HostRuleDAL      *dal.HostRule
	PromotionBIZ     *Promotion
	MaterialVideoBIZ *MaterialVideo
}

// Query 查询托管规则列表
func (a *HostRule) Query(ctx context.Context, params schema.HostRuleQueryParam) (*schema.HostRuleQueryResult, error) {
	result, err := a.HostRuleDAL.Query(ctx, params, schema.HostRuleQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.HostRulesOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get 获取指定托管规则详情
func (a *HostRule) Get(ctx context.Context, id uint) (*schema.HostRule, error) {
	item, err := a.HostRuleDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.NotFound("", "托管规则不存在")
	}
	return item, nil
}

// UpdateStatus 修改托管规则状态
func (a *HostRule) UpdateStatus(ctx context.Context, id uint, form *schema.HostRuleUpdateStatusForm) error {
	item, err := a.HostRuleDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.NotFound("", "托管规则不存在")
	}
	return a.HostRuleDAL.UpdateStatus(ctx, id, form.Status)
}

// GetTargetsByAccount 根据账户ID列表和目标类型查询目标列表
func (a *HostRule) GetTargetsByAccount(ctx context.Context, req *schema.TargetByAccountReq) ([]schema.TargetItem, error) {
	// 转换为 int64 切片（advertiser_id 为 bigint）
	ids := make([]int64, 0, len(req.AccountIDs))
	for _, id := range req.AccountIDs {
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, errors.BadRequest("", "无效的账户ID: "+id)
		}
		ids = append(ids, n)
	}

	agentID, err := strconv.ParseInt(req.AgentID, 10, 64)
	if err != nil {
		return nil, errors.BadRequest("", "无效的代理商ID: "+req.AgentID)
	}

	items := make([]schema.TargetItem, 0)
	switch req.Target {
	case "promotion":
		promotions, err := a.PromotionBIZ.FindByAccountIDs(ctx, agentID, ids, []string{"id", "promotion_id", "promotion_name", "advertiser_id"})
		if err != nil {
			return nil, err
		}
		for _, p := range promotions {
			items = append(items, schema.TargetItem{
				ID:   strconv.FormatInt(p.PromotionID, 10),
				Name: p.PromotionName,
			})
		}
	case "creative":
		materials, err := a.MaterialVideoBIZ.FindByAccountIDs(ctx, agentID, ids, []string{"id", "material_id", "file_name", "advertiser_id"})
		if err != nil {
			return nil, err
		}
		for _, m := range materials {
			items = append(items, schema.TargetItem{
				ID:   strconv.FormatInt(m.MaterialID, 10),
				Name: m.FileName,
			})
		}
	}

	return items, nil
}

// Save 新增/编辑托管规则
func (a *HostRule) Save(ctx context.Context, formItem *schema.HostRuleForm) (*schema.HostRule, error) {
	if formItem.ID > 0 {
		// 编辑
		item, err := a.HostRuleDAL.Get(ctx, formItem.ID)
		if err != nil {
			return nil, err
		} else if item == nil {
			return nil, errors.NotFound("", "托管规则不存在")
		}
		formItem.FillTo(item)
		if err := a.HostRuleDAL.Update(ctx, item); err != nil {
			return nil, err
		}
		return item, nil
	}

	// 新增
	item := &schema.HostRule{
		Status: schema.HostRuleStatusEnabled,
	}
	formItem.FillTo(item)
	if err := a.HostRuleDAL.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
