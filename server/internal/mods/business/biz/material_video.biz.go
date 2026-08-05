package biz

import (
	"context"
	"fmt"
	"time"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
)

// MaterialVideo 视频素材业务逻辑
type MaterialVideo struct {
	MaterialVideoDAL *dal.MaterialVideo
	Oceanengine      *Oceanengine
}

// FindByAccountIDs 根据账户ID列表查询视频素材列表，fields 指定查询字段
// 按账户维度同步：有记录则增量同步（按 advertiser_id+material_id 去重），无记录则全量同步
func (a *MaterialVideo) FindByAccountIDs(ctx context.Context, accountIDs []int64, fields []string) ([]schema.MaterialVideo, error) {
	// 查询每个账户的最新 updated_at
	lastSyncMap, err := a.MaterialVideoDAL.GetLastUpdatedAtByAccountIDs(ctx, accountIDs)
	if err != nil {
		return nil, fmt.Errorf("查询素材同步时间失败: %w", err)
	}

	endDate := time.Now().Format("2006-01-02")
	var allNewItems []schema.MaterialVideo

	accessToken := "" // TODO: 从配置/上下文中获取 accessToken

	for _, advertiserID := range accountIDs {
		lastUpdatedAt, hasRecord := lastSyncMap[advertiserID]

		var startDate, apiEndDate string
		if hasRecord {
			// 有记录：增量同步 startDate=updated_at, endDate=今天
			startDate = lastUpdatedAt.Format("2006-01-02")
			apiEndDate = endDate
		}
		// 无记录：startDate="" endDate="" 全量拉取
		items, err := a.Oceanengine.GetVideoMaterial(accessToken, advertiserID, startDate, apiEndDate)
		if err != nil {
			return nil, fmt.Errorf("拉取账户 %d 视频素材失败: %w", advertiserID, err)
		}
		if len(items) == 0 {
			continue
		}

		if hasRecord {
			// 增量：过滤已存在的 (advertiser_id, material_id)
			existingKeys, err := a.MaterialVideoDAL.FindExistingKeys(ctx, items)
			if err != nil {
				return nil, fmt.Errorf("查询已存在素材失败: %w", err)
			}
			for _, item := range items {
				key := fmt.Sprintf("%d_%d", item.AdvertiserID, item.MaterialID)
				if !existingKeys[key] {
					allNewItems = append(allNewItems, item)
				}
			}
		} else {
			// 全量：直接插入
			allNewItems = append(allNewItems, items...)
		}
	}

	// 插入新数据
	if len(allNewItems) > 0 {
		if err := a.MaterialVideoDAL.SaveBatch(ctx, allNewItems); err != nil {
			return nil, fmt.Errorf("写入视频素材数据失败: %w", err)
		}
	}

	// 重新查询 DB 返回统一字段
	return a.MaterialVideoDAL.FindByAccountIDs(ctx, accountIDs, fields)
}
