package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

// Promotion 广告数据访问
type Promotion struct {
	DB *gorm.DB
}

// FindByAccountIDs 根据账户ID列表查询启用状态的广告列表，fields 指定查询字段
func (a *Promotion) FindByAccountIDs(ctx context.Context, accountIDs []int64, fields []string) ([]schema.Promotion, error) {
	var list []schema.Promotion
	db := util.GetDB(ctx, a.DB).
		// Where("advertiser_id IN ? AND status_first = ?", accountIDs, schema.PromotionStatusEnable)
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

// SaveBatch 批量保存广告数据
func (a *Promotion) SaveBatch(ctx context.Context, items []schema.Promotion) error {
	if len(items) == 0 {
		return nil
	}
	err := util.GetDB(ctx, a.DB).Create(&items).Error
	return errors.WithStack(err)
}
