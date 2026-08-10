package biz

import (
	"context"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"
)

// HostAccount 托管账户业务逻辑
type HostAccount struct {
	HostAccountDAL *dal.HostAccount
}

// Query 查询托管账户列表
func (a *HostAccount) Query(ctx context.Context, params schema.HostAccountQueryParam) (*schema.HostAccountQueryResult, error) {
	params.Pagination = true
	result, err := a.HostAccountDAL.Query(ctx, params, schema.HostAccountQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.HostAccountOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get 获取指定托管账户详情
func (a *HostAccount) Get(ctx context.Context, id uint) (*schema.HostAccount, error) {
	item, err := a.HostAccountDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.NotFound("", "托管账户不存在")
	}
	return item, nil
}

// Create 新增托管账户
func (a *HostAccount) Create(ctx context.Context, formItem *schema.HostAccountForm) (*schema.HostAccount, error) {
	item := &schema.HostAccount{Status: 1}
	formItem.FillTo(item)
	if err := a.HostAccountDAL.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

// Update 更新托管账户
func (a *HostAccount) Update(ctx context.Context, id uint, formItem *schema.HostAccountForm) error {
	item, err := a.HostAccountDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.NotFound("", "托管账户不存在")
	}
	formItem.FillTo(item)
	return a.HostAccountDAL.Update(ctx, item)
}

// Delete 软删除托管账户
func (a *HostAccount) Delete(ctx context.Context, id uint) error {
	item, err := a.HostAccountDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.NotFound("", "托管账户不存在")
	}
	return a.HostAccountDAL.Delete(ctx, id)
}