package biz

import (
	"context"
	"fmt"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"
)

// DeleteUnauditedMaterial 素材删除记录业务逻辑
type DeleteUnauditedMaterial struct {
	DeleteUnauditedMaterialDAL *dal.DeleteUnauditedMaterial
	Oceanengine                *Oceanengine
}

// Query 查询素材删除记录列表
func (a *DeleteUnauditedMaterial) Query(ctx context.Context, params schema.DeleteUnauditedMaterialQueryParam) (*schema.DeleteUnauditedMaterialQueryResult, error) {
	params.Pagination = true
	result, err := a.DeleteUnauditedMaterialDAL.Query(ctx, params, schema.DeleteUnauditedMaterialQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.DeleteUnauditedMaterialOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get 获取指定素材删除记录详情
func (a *DeleteUnauditedMaterial) Get(ctx context.Context, id uint) (*schema.DeleteUnauditedMaterial, error) {
	item, err := a.DeleteUnauditedMaterialDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.NotFound("", "记录不存在")
	}
	return item, nil
}

// FetchUnauditedMaterial 拉取未审核素材
func (a *DeleteUnauditedMaterial) FetchUnauditedMaterial(ctx context.Context, accountID int64, accountIDs []int64) ([]*schema.DeleteUnauditedMaterial, error) {
	filtering := map[string]interface{}{
		"status_first": schema.PromotionStatusEnable,
	}
	fields := []string{"promotion_id", "promotion_materials", "advertiser_id", "promotion_name"}

	result := make([]*schema.DeleteUnauditedMaterial, 0)
	for _, advertiserID := range accountIDs {
		items, err := a.Oceanengine.GetRefPromotionData(ctx, accountID, advertiserID, filtering, fields)
		if err != nil {
			return nil, fmt.Errorf("拉取账户 %d 广告失败: %w", advertiserID, err)
		}
		for _, item := range items {
			if item.PromotionMaterials == nil {
				continue
			}
			for _, videoMaterial := range item.PromotionMaterials.VideoMaterialList {
				if videoMaterial.MaterialStatus != "MATERIAL_STATUS_OFFLINE_AUDIT" {
					continue
				}
				result = append(result, &schema.DeleteUnauditedMaterial{
					AdvertiserID: advertiserID,
					PromotionID:  item.PromotionId,
					MaterialID:   videoMaterial.MaterialID,
				})
			}
		}
	}
	return result, nil
}

// DeleteAndSave 删除素材并保存记录
func (a *DeleteUnauditedMaterial) DeleteAndSave(ctx context.Context, req *schema.UnAudititedMaterialReq) error {
	for _, m := range req.Materials {
		record := &schema.DeleteUnauditedMaterial{
			MaterialID:   m.MaterialID,
			AdvertiserID: m.AdvertiserID,
			PromotionID:  m.PromotionID,
			IsDeleted:    "pending",
		}

		_, err := a.Oceanengine.DeleteMaterialUnderPromotion(ctx, req.AccountID, m.MaterialID, m.PromotionID, m.AdvertiserID)
		if err != nil {
			record.IsDeleted = "failed"
			record.ErrorMsg = err.Error()
		} else {
			record.IsDeleted = "deleted"
		}

		if err := a.DeleteUnauditedMaterialDAL.Create(ctx, record); err != nil {
			return fmt.Errorf("保存删除记录失败: %w", err)
		}
	}
	return nil
}

// RetryFailedDelete 重试删除失败的记录（is_deleted=failed 且 retry_times < 3）
func (a *DeleteUnauditedMaterial) RetryFailedDelete(ctx context.Context, accountID int64) error {
	records, err := a.DeleteUnauditedMaterialDAL.QueryFailedForRetry(ctx)
	if err != nil {
		return fmt.Errorf("查询失败记录失败: %w", err)
	}

	if len(records) == 0 {
		return nil
	}

	for _, record := range records {
		_, err := a.Oceanengine.DeleteMaterialUnderPromotion(ctx, accountID, record.MaterialID, record.PromotionID, record.AdvertiserID)
		if err != nil {
			record.IsDeleted = "failed"
			record.ErrorMsg = err.Error()
		} else {
			record.IsDeleted = "deleted"
			record.ErrorMsg = ""
		}
		record.RetryTimes++

		if err := a.DeleteUnauditedMaterialDAL.Update(ctx, record); err != nil {
			return fmt.Errorf("更新删除记录失败(id=%d): %w", record.ID, err)
		}
	}

	return nil
}
