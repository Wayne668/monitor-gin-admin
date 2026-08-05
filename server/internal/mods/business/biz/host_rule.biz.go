package biz

import (
	"context"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"
)

// HostRule 托管规则业务逻辑
type HostRule struct {
	HostRuleDAL *dal.HostRule
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
