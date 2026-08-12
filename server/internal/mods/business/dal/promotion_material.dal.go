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
		Select("promotion_id, ANY_VALUE(promotion_name) as promotion_name, ANY_VALUE(advertiser_id) as advertiser_id").
		Where("advertiser_id IN ? AND status_first = ?", accountIDs, schema.PromotionStatusEnable).
		Group("promotion_id")
	err := db.Find(&list).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

// FindMaterialsByAccountIDs 根据账户ID列表查询素材状态正常的素材（JOIN 获取 file_name）
func (a *PromotionMaterial) FindMaterialsByAccountIDs(ctx context.Context, accountIDs []int64, materialStatus string) ([]schema.PromotionMaterial, error) {
	var list []schema.PromotionMaterial
	db := util.GetDB(ctx, a.DB).
		Table("nb_promotion_material pm").
		Select("pm.material_id, ANY_VALUE(pm.advertiser_id) AS advertiser_id, ANY_VALUE(mv.file_name) AS file_name, ANY_VALUE(pm.promotion_id) AS promotion_id").
		Joins("LEFT JOIN nb_material_video mv ON pm.material_id = mv.material_id").
		Where("pm.advertiser_id IN ? AND pm.material_status = ?", accountIDs, schema.MaterialStatusOK).
		Group("pm.material_id")
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

// UpsertBatch 批量 upsert（唯一键：advertiser_id + promotion_id + material_id）
func (a *PromotionMaterial) UpsertBatch(ctx context.Context, items []schema.PromotionMaterial) error {
	if len(items) == 0 {
		return nil
	}
	db := util.GetDB(ctx, a.DB)
	for _, item := range items {
		var existing schema.PromotionMaterial
		err := db.Where("advertiser_id = ? AND promotion_id = ? AND material_id = ?",
			item.AdvertiserID, item.PromotionID, item.MaterialID).First(&existing).Error
		if err == nil {
			// 存在则更新
			updateErr := db.Model(&existing).Updates(map[string]interface{}{
				"status_first":    item.StatusFirst,
				"material_status": item.MaterialStatus,
				"opt_status":      item.OptStatus,
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

// FindExistingTargetIDs 返回已存在的 target_id 集合
func (a *PromotionMaterial) FindExistingTargetIDs(ctx context.Context, targetIDs []string, advertiserID int64, target string) ([]int64, error) {
	if len(targetIDs) == 0 {
		return nil, nil
	}
	var list []int64
	var err error
	if target == "promotion" {
		err = util.GetDB(ctx, a.DB).
			Select("promotion_id").
			Where("advertiser_id = ?", advertiserID).
			Where("promotion_id IN ?", targetIDs).
			Find(&list).Error
	} else {
		err = util.GetDB(ctx, a.DB).
			Select("material_id").
			Where("advertiser_id = ?", advertiserID).
			Where("material_id IN ?", targetIDs).
			Find(&list).Error
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list, nil
}
