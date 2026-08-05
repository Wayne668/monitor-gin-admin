package dal

import (
	"context"
	"fmt"
	"time"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

// MaterialVideo 视频素材数据访问
type MaterialVideo struct {
	DB *gorm.DB
}

// FindByAccountIDs 根据账户ID列表查询视频素材列表，fields 指定查询字段
func (a *MaterialVideo) FindByAccountIDs(ctx context.Context, accountIDs []int64, fields []string) ([]schema.MaterialVideo, error) {
	var list []schema.MaterialVideo
	db := util.GetDB(ctx, a.DB).
		Where("advertiser_id IN ?", accountIDs)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	err := db.Find(&list).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

// SaveBatch 批量保存视频素材数据
func (a *MaterialVideo) SaveBatch(ctx context.Context, items []schema.MaterialVideo) error {
	if len(items) == 0 {
		return nil
	}
	err := util.GetDB(ctx, a.DB).Create(&items).Error
	return errors.WithStack(err)
}

// GetLastUpdatedAtByAccountIDs 查询每个账户的最新 updated_at
func (a *MaterialVideo) GetLastUpdatedAtByAccountIDs(ctx context.Context, accountIDs []int64) (map[int64]time.Time, error) {
	type result struct {
		AdvertiserID int64     `gorm:"column:advertiser_id"`
		MaxUpdatedAt time.Time `gorm:"column:max_updated_at"`
	}
	var results []result
	err := util.GetDB(ctx, a.DB).
		Table("nb_material_video").
		Select("advertiser_id, MAX(updated_at) as max_updated_at").
		Where("advertiser_id IN ?", accountIDs).
		Group("advertiser_id").
		Find(&results).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	m := make(map[int64]time.Time, len(results))
	for _, r := range results {
		m[r.AdvertiserID] = r.MaxUpdatedAt
	}
	return m, nil
}

// FindExistingKeys 查询已存在的 (advertiser_id, material_id) 组合
func (a *MaterialVideo) FindExistingKeys(ctx context.Context, items []schema.MaterialVideo) (map[string]bool, error) {
	if len(items) == 0 {
		return nil, nil
	}
	// 收集涉及的所有 advertiser_id
	advertiserIDSet := make(map[int64]bool)
	for _, item := range items {
		advertiserIDSet[item.AdvertiserID] = true
	}
	advertiserIDs := make([]int64, 0, len(advertiserIDSet))
	for id := range advertiserIDSet {
		advertiserIDs = append(advertiserIDs, id)
	}

	// 查询这些账户下已存在的所有 (advertiser_id, material_id)
	var existing []schema.MaterialVideo
	err := util.GetDB(ctx, a.DB).
		Select("advertiser_id, material_id").
		Where("advertiser_id IN ?", advertiserIDs).
		Find(&existing).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	m := make(map[string]bool, len(existing))
	for _, e := range existing {
		m[fmt.Sprintf("%d_%d", e.AdvertiserID, e.MaterialID)] = true
	}
	return m, nil
}
