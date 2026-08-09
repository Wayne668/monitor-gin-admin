package biz

import (
	"context"
	"fmt"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
)

// Promotion 广告业务逻辑
type Promotion struct {
	PromotionDAL *dal.Promotion
	Oceanengine  *Oceanengine
}

// FindByAccountIDs 根据账户ID列表查询启用状态的广告列表，fields 指定查询字段
// DB无数据时调用Oceanengine API拉取并写入数据库后返回
func (a *Promotion) FindByAccountIDs(ctx context.Context, agentID int64, accountIDs []int64, fields []string) ([]schema.Promotion, error) {
	list, err := a.PromotionDAL.FindByAccountIDs(ctx, accountIDs, fields)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list, nil
	}

	// DB无数据，调用Oceanengine API拉取
	if agentID == 0 {
		// 未配置 accountID，直接返回空列表，避免调用API失败导致500
		return nil, nil
	}

	apiItems, err := a.FetchFromOceanengine(ctx, accountIDs, agentID)
	if err != nil {
		return nil, fmt.Errorf("从Oceanengine拉取广告数据失败: %w", err)
	}
	if len(apiItems) == 0 {
		return nil, nil
	}

	if err := a.PromotionDAL.SaveBatch(ctx, apiItems); err != nil {
		return nil, fmt.Errorf("写入广告数据失败: %w", err)
	}

	// 重新查询DB获取统一字段
	return a.PromotionDAL.FindByAccountIDs(ctx, accountIDs, fields)
}

// FetchFromOceanengine 调用Oceanengine API拉取广告数据
func (a *Promotion) FetchFromOceanengine(ctx context.Context, accountIDs []int64, accountID int64) ([]schema.Promotion, error) {
	filtering := map[string]interface{}{
		"status_first": schema.PromotionStatusEnable,
	}
	fields := []string{"promotion_id", "promotion_name", "status_first", "status_second", "opt_status"}

	result := make([]schema.Promotion, 0)
	for _, advertiserID := range accountIDs {
		items, err := a.Oceanengine.GetRefPromotionData(ctx, accountID, advertiserID, filtering, fields)
		if err != nil {
			return nil, fmt.Errorf("拉取账户 %d 广告失败: %w", advertiserID, err)
		}
		for _, item := range items {
			result = append(result, schema.Promotion{
				AdvertiserID:  advertiserID,
				PromotionID:   item.PromotionId,
				PromotionName: item.PromotionName,
				StatusFirst:   item.StatusFirst,
				StatusSecond:  item.StatusSecond,
				OptStatus:     item.OptStatus,
			})
		}
	}
	return result, nil
}
