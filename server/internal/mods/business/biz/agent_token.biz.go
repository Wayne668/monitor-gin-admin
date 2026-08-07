package biz

import (
	"context"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"
)

// AgentToken 代理商账号授权业务逻辑
type AgentToken struct {
	AgentTokenDAL *dal.AgentToken
}

// Query 查询代理商账号授权列表
func (a *AgentToken) Query(ctx context.Context, params schema.AgentTokenQueryParam) (*schema.AgentTokenQueryResult, error) {
	params.Pagination = true
	result, err := a.AgentTokenDAL.Query(ctx, params, schema.AgentTokenQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.AgentTokensOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get 获取指定代理商账号授权详情
func (a *AgentToken) Get(ctx context.Context, id uint) (*schema.AgentToken, error) {
	item, err := a.AgentTokenDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.NotFound("", "账户token不存在")
	}
	return item, nil
}

// Create 新增代理商账号授权
func (a *AgentToken) Create(ctx context.Context, formItem *schema.AgentTokenForm) (*schema.AgentToken, error) {
	item := &schema.AgentToken{}
	formItem.FillTo(item)
	if err := a.AgentTokenDAL.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

// Update 更新代理商账号授权
func (a *AgentToken) Update(ctx context.Context, id uint, formItem *schema.AgentTokenForm) error {
	item, err := a.AgentTokenDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.NotFound("", "账户token不存在")
	}
	formItem.FillTo(item)
	return a.AgentTokenDAL.Update(ctx, item)
}

// Delete 删除代理商账号授权
func (a *AgentToken) Delete(ctx context.Context, id uint) error {
	item, err := a.AgentTokenDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.NotFound("", "账户token不存在")
	}
	return a.AgentTokenDAL.Delete(ctx, id)
}
