package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetHostFieldDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.HostField))
}

// HostField 托管字段数据访问
type HostField struct {
	DB *gorm.DB
}

// Query 查询托管字段列表
func (a *HostField) Query(ctx context.Context, params schema.HostFieldQueryParam, opts ...schema.HostFieldQueryOptions) (*schema.HostFieldQueryResult, error) {
	var opt schema.HostFieldQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetHostFieldDB(ctx, a.DB).Where("status = ?", 1)

	if v := params.Cate; len(v) > 0 {
		db = db.Where("cate = ?", v)
	}

	var list schema.HostFields
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &schema.HostFieldQueryResult{
		PageResult: pageResult,
		Data:       list,
	}, nil
}

// FindAll 查询全部托管字段（不分页，用于前端下拉渲染）
func (a *HostField) FindAll(ctx context.Context) (schema.HostFields, error) {
	var list schema.HostFields
	err := GetHostFieldDB(ctx, a.DB).Where("status = ?", 1).Order("id ASC").Find(&list).Error
	return list, errors.WithStack(err)
}

// Get 获取指定托管字段
func (a *HostField) Get(ctx context.Context, id uint) (*schema.HostField, error) {
	item := new(schema.HostField)
	ok, err := util.FindOne(ctx, GetHostFieldDB(ctx, a.DB).Where("id=?", id), util.QueryOptions{}, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Create 新增托管字段
func (a *HostField) Create(ctx context.Context, item *schema.HostField) error {
	result := GetHostFieldDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update 更新托管字段
func (a *HostField) Update(ctx context.Context, item *schema.HostField) error {
	result := GetHostFieldDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete 删除托管字段
func (a *HostField) Delete(ctx context.Context, id uint) error {
	result := GetHostFieldDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.HostField))
	return errors.WithStack(result.Error)
}

// FindByField 根据 field 名称查找托管字段
func (a *HostField) FindByField(ctx context.Context, field string) (*schema.HostField, error) {
	item := new(schema.HostField)
	ok, err := util.FindOne(ctx, GetHostFieldDB(ctx, a.DB).Where("field = ? AND status = 1", field), util.QueryOptions{}, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}
