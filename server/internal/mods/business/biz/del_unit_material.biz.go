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
	PromotionMaterialDAL       *dal.PromotionMaterial
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
func (a *DeleteUnauditedMaterial) FetchUnauditedMaterial(ctx context.Context, agentID int64, accountIDs []int64) ([]*schema.DeleteUnauditedMaterial, error) {
	filtering := map[string]interface{}{
		"status_first": schema.PromotionStatusEnable,
	}
	fields := []string{"promotion_id", "promotion_materials", "advertiser_id", "promotion_name"}

	result := make([]*schema.DeleteUnauditedMaterial, 0)
	for _, advertiserID := range accountIDs {
		items, err := a.Oceanengine.GetRefPromotionData(ctx, agentID, advertiserID, filtering, fields)
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

// GetUnAudititedMaterial 直接查询 nb_promotion_material 获取未审核素材
func (a *DeleteUnauditedMaterial) GetUnAudititedMaterial(ctx context.Context, accountIDs []int64) ([]schema.PromotionMaterial, error) {
	return a.PromotionMaterialDAL.FindMaterialsByAccountIDs(ctx, accountIDs, schema.MaterialStatusOfflineAudit)
}

// GetUnAudititedMaterialWithFallback 先查本地表，若为空则回退到媒体接口拉取
func (a *DeleteUnauditedMaterial) GetUnAudititedMaterialWithFallback(ctx context.Context, agentID int64, accountIDs []int64) ([]schema.PromotionMaterial, error) {
	items, err := a.GetUnAudititedMaterial(ctx, accountIDs)
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items, nil
	}

	// 本地表无数据，回退到媒体接口
	fetched, err := a.FetchUnauditedMaterial(ctx, agentID, accountIDs)
	if err != nil {
		return nil, err
	}
	result := make([]schema.PromotionMaterial, 0, len(fetched))
	for _, m := range fetched {
		result = append(result, schema.PromotionMaterial{
			AdvertiserID: m.AdvertiserID,
			PromotionID:  m.PromotionID,
			MaterialID:   m.MaterialID,
			FileName:     m.MaterialName,
		})
	}
	return result, nil
}

// DeleteAndSave 删除素材并保存记录，返回失败的记录列表
func (a *DeleteUnauditedMaterial) DeleteAndSave(ctx context.Context, req *schema.UnAudititedMaterialReq) ([]*schema.DeleteUnauditedMaterial, error) {
	accountID := util.ParseInt64(req.AccountID)
	var failedRecords []*schema.DeleteUnauditedMaterial
	for _, m := range req.Materials {
		materialID := util.ParseInt64(m.MaterialID)
		promotionID := util.ParseInt64(m.PromotionID)
		advertiserID := util.ParseInt64(m.AdvertiserID)

		record := &schema.DeleteUnauditedMaterial{
			MaterialID:   materialID,
			AdvertiserID: advertiserID,
			PromotionID:  promotionID,
			MaterialName: m.MaterialName,
			IsDeleted:    "pending",
		}

		_, err := a.Oceanengine.DeleteMaterialUnderPromotion(ctx, accountID, materialID, promotionID, advertiserID)
		if err != nil {
			record.IsDeleted = "failed"
			record.ErrorMsg = err.Error()
			failedRecords = append(failedRecords, record)
		} else {
			record.IsDeleted = "deleted"
			// 更新 nb_promotion_material 素材状态为已删除
			if updateErr := a.PromotionMaterialDAL.UpdateMaterialStatus(ctx, advertiserID, promotionID, materialID, schema.MaterialStatusDelete); updateErr != nil {
				record.ErrorMsg = fmt.Sprintf("更新素材状态失败: %v", updateErr)
				failedRecords = append(failedRecords, record)
			}
		}

		if err := a.DeleteUnauditedMaterialDAL.Create(ctx, record); err != nil {
			return nil, fmt.Errorf("保存删除记录失败: %w", err)
		}
	}
	return failedRecords, nil
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
