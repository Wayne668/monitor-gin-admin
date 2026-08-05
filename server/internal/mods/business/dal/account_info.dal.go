package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"time"

	"gorm.io/gorm"
)

// AccountInfo management
type AccountInfo struct {
	DB *gorm.DB
}

// FilterExistingAdvertiserIDs 过滤已存在的广告主ID，返回不存在的新ID列表
func (a *AccountInfo) FilterExistingAdvertiserIDs(advertiserIDs []int64) ([]int64, error) {
	if len(advertiserIDs) == 0 {
		return nil, nil
	}

	batchSize := 1000
	existingSet := make(map[int64]bool)

	for i := 0; i < len(advertiserIDs); i += batchSize {
		end := i + batchSize
		if end > len(advertiserIDs) {
			end = len(advertiserIDs)
		}
		chunk := advertiserIDs[i:end]

		var batchIDs []int64
		err := a.DB.Model(&schema.AccountInfo{}).
			Where("advertiser_id IN ?", chunk).
			Pluck("advertiser_id", &batchIDs).Error
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, id := range batchIDs {
			existingSet[id] = true
		}
	}

	var newIDs []int64
	for _, id := range advertiserIDs {
		if !existingSet[id] {
			newIDs = append(newIDs, id)
		}
	}

	return newIDs, nil
}

// FindEnabledAdvertisers 查询状态为 STATUS_ENABLE 的账户列表
func (a *AccountInfo) FindEnabledAdvertisers(ctx context.Context) ([]schema.AccountInfo, error) {
	var list []schema.AccountInfo
	err := a.DB.Where("advertiser_status = ?", "STATUS_ENABLE").
		Order("advertiser_id ASC").
		Find(&list).Error
	return list, errors.WithStack(err)
}

// SaveToTable 批量保存账户详情到数据库
func (a *AccountInfo) SaveToTable(details []schema.AccountDetail) error {
	if len(details) == 0 {
		return nil
	}

	records := make([]schema.AccountInfo, 0, len(details))
	now := time.Now()
	for _, d := range details {
		records = append(records, schema.AccountInfo{
			AdvertiserId:       &d.AdvertiserID,
			AdvertiserName:     &d.AdvertiserName,
			AdvCompanyId:       &d.AdvCompanyID,
			AdvCompanyName:     &d.AdvCompanyName,
			FirstIndustryName:  &d.FirstIndustryName,
			SecondIndustryName: &d.SecondIndustryName,
			AdvertiserStatus:   &d.AdvertiserStatus,
			UpdatedAt:          &now,
		})
	}

	err := a.DB.Create(&records).Error
	return errors.WithStack(err)
}
