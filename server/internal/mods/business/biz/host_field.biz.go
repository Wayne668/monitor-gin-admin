package biz

import (
	"context"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"
)

// HostField 托管字段业务逻辑
type HostField struct {
	HostFieldDAL *dal.HostField
}

// Query 查询托管字段列表
func (a *HostField) Query(ctx context.Context, params schema.HostFieldQueryParam) (*schema.HostFieldQueryResult, error) {
	params.Pagination = true
	result, err := a.HostFieldDAL.Query(ctx, params, schema.HostFieldQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.HostFieldsOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindAll 查询全部托管字段（不分页）
func (a *HostField) FindAll(ctx context.Context) (schema.HostFields, error) {
	return a.HostFieldDAL.FindAll(ctx)
}

// Get 获取指定托管字段详情
func (a *HostField) Get(ctx context.Context, id uint) (*schema.HostField, error) {
	item, err := a.HostFieldDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.NotFound("", "托管字段不存在")
	}
	return item, nil
}

// Create 新增托管字段
func (a *HostField) Create(ctx context.Context, formItem *schema.HostFieldForm) (*schema.HostField, error) {
	item := &schema.HostField{}
	formItem.FillTo(item)
	if err := a.HostFieldDAL.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

// Update 更新托管字段
func (a *HostField) Update(ctx context.Context, id uint, formItem *schema.HostFieldForm) error {
	item, err := a.HostFieldDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.NotFound("", "托管字段不存在")
	}
	formItem.FillTo(item)
	return a.HostFieldDAL.Update(ctx, item)
}

// Delete 删除托管字段
func (a *HostField) Delete(ctx context.Context, id uint) error {
	item, err := a.HostFieldDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.NotFound("", "托管字段不存在")
	}
	return a.HostFieldDAL.Delete(ctx, id)
}
