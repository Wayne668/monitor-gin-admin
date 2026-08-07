package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetDeleteUnauditedMaterialDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.DeleteUnauditedMaterial))
}

// DeleteUnauditedMaterial 素材删除记录数据访问
type DeleteUnauditedMaterial struct {
	DB *gorm.DB
}

// Query 查询素材删除记录列表
func (a *DeleteUnauditedMaterial) Query(ctx context.Context, params schema.DeleteUnauditedMaterialQueryParam, opts ...schema.DeleteUnauditedMaterialQueryOptions) (*schema.DeleteUnauditedMaterialQueryResult, error) {
	var opt schema.DeleteUnauditedMaterialQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetDeleteUnauditedMaterialDB(ctx, a.DB)

	if v := params.MaterialID; v != nil {
		db = db.Where("material_id = ?", *v)
	}
	if v := params.AdvertiserID; v != nil {
		db = db.Where("advertiser_id = ?", *v)
	}
	if v := params.IsDeleted; len(v) > 0 {
		db = db.Where("is_deleted = ?", v)
	}
	if v := params.MaterialName; len(v) > 0 {
		db = db.Where("material_name LIKE ?", "%"+v+"%")
	}

	var list schema.DeleteUnauditedMaterials
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &schema.DeleteUnauditedMaterialQueryResult{
		PageResult: pageResult,
		Data:       list,
	}, nil
}

// Get 获取指定素材删除记录
func (a *DeleteUnauditedMaterial) Get(ctx context.Context, id uint) (*schema.DeleteUnauditedMaterial, error) {
	item := new(schema.DeleteUnauditedMaterial)
	ok, err := util.FindOne(ctx, GetDeleteUnauditedMaterialDB(ctx, a.DB).Where("id=?", id), util.QueryOptions{}, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Create 新增素材删除记录
func (a *DeleteUnauditedMaterial) Create(ctx context.Context, item *schema.DeleteUnauditedMaterial) error {
	result := GetDeleteUnauditedMaterialDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update 更新素材删除记录
func (a *DeleteUnauditedMaterial) Update(ctx context.Context, item *schema.DeleteUnauditedMaterial) error {
	result := GetDeleteUnauditedMaterialDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete 删除素材删除记录
func (a *DeleteUnauditedMaterial) Delete(ctx context.Context, id uint) error {
	result := GetDeleteUnauditedMaterialDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.DeleteUnauditedMaterial))
	return errors.WithStack(result.Error)
}
