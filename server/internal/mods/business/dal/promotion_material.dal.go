package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

// PromotionMaterial 广告素材数据访问
type PromotionMaterial struct {
	DB *gorm.DB
}

// FindPromotionsByAccountIDs 根据账户ID列表查询启用状态的广告（promotion_id 去重）
func (a *PromotionMaterial) FindPromotionsByAccountIDs(ctx context.Context, accountIDs []int64) ([]schema.PromotionMaterial, error) {
	var list []schema.PromotionMaterial
	db := util.GetDB(ctx, a.DB).
		Where("advertiser_id IN ? AND status_first = ?", accountIDs, schema.PromotionStatusEnable).
		Group("promotion_id")
	err := db.Find(&list).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

// FindMaterialsByAccountIDs 根据账户ID列表查询素材状态正常的素材
func (a *PromotionMaterial) FindMaterialsByAccountIDs(ctx context.Context, accountIDs []int64) ([]schema.PromotionMaterial, error) {
	var list []schema.PromotionMaterial
	db := util.GetDB(ctx, a.DB).
		Where("advertiser_id IN ? AND material_status = ?", accountIDs, schema.MaterialStatusOK)
	err := db.Find(&list).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

// FindExistingMaterialIDs 返回已存在的 material_id 集合
func (a *PromotionMaterial) FindExistingMaterialIDs(ctx context.Context, materialIDs []int64) (map[int64]bool, error) {
	if len(materialIDs) == 0 {
		return nil, nil
	}
	var list []schema.PromotionMaterial
	err := util.GetDB(ctx, a.DB).
		Select("material_id").
		Where("material_id IN ?", materialIDs).
		Find(&list).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	m := make(map[int64]bool, len(list))
	for _, item := range list {
		m[item.MaterialID] = true
	}
	return m, nil
}

// UpsertBatch 批量 upsert（存在则更新，不存在则插入）
func (a *PromotionMaterial) UpsertBatch(ctx context.Context, items []schema.PromotionMaterial) error {
	if len(items) == 0 {
		return nil
	}
	db := util.GetDB(ctx, a.DB)
	for _, item := range items {
		var existing schema.PromotionMaterial
		err := db.Where("material_id = ?", item.MaterialID).First(&existing).Error
		if err == nil {
			// 存在则更新
			updateErr := db.Model(&existing).Updates(map[string]interface{}{
				"status_first":    item.StatusFirst,
				"status_second":   item.StatusSecond,
				"material_status": item.MaterialStatus,
			}).Error
			if updateErr != nil {
				return errors.WithStack(updateErr)
			}
		} else {
			// 不存在则插入
			createErr := db.Create(&item).Error
			if createErr != nil {
				return errors.WithStack(createErr)
			}
		}
	}
	return nil
}
