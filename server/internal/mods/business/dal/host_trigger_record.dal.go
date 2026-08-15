package dal

import (
	"context"
	"time"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

// HostTriggerRecord 触发记录数据访问
type HostTriggerRecord struct {
	DB *gorm.DB
}

func (a *HostTriggerRecord) getDB(ctx context.Context) *gorm.DB {
	return util.GetDB(ctx, a.DB).Model(new(schema.HostTriggerRecord))
}

// Create 新增触发记录
func (a *HostTriggerRecord) Create(ctx context.Context, item *schema.HostTriggerRecord) error {
	result := a.getDB(ctx).Create(item)
	return errors.WithStack(result.Error)
}

// Query 查询触发记录列表
func (a *HostTriggerRecord) Query(ctx context.Context, params schema.HostTriggerRecordQueryParam) (*schema.HostTriggerRecordQueryResult, error) {
	db := a.getDB(ctx)

	if v := params.RuleID; v != nil {
		db = db.Where("rule_id = ?", *v)
	}
	if v := params.AdvertiserID; v != nil {
		db = db.Where("advertiser_id = ?", *v)
	}
	if v := params.Target; v != "" {
		db = db.Where("target = ?", v)
	}
	if v := params.ExecuteStatus; v != "" {
		db = db.Where("execute_status = ?", v)
	}

	var list schema.HostTriggerRecords
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, util.QueryOptions{}, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &schema.HostTriggerRecordQueryResult{
		PageResult: pageResult,
		Data:       list,
	}, nil
}

// QueryPending 查询待执行的触发记录
func (a *HostTriggerRecord) QueryPending(ctx context.Context) ([]*schema.HostTriggerRecord, error) {
	var list []*schema.HostTriggerRecord
	result := a.getDB(ctx).Where("execute_status = ?", "pending").Find(&list)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return list, nil
}

// ExistsPending 判断是否已存在待执行的触发记录（同规则+同账户+同目标）
func (a *HostTriggerRecord) ExistsPending(ctx context.Context, ruleID int, advertiserID int64, target string, targetID int64) (bool, error) {
	var count int64
	db := a.getDB(ctx).Where("execute_status = ? AND rule_id = ? AND advertiser_id = ?", "pending", ruleID, advertiserID)
	switch target {
	case "creative":
		db = db.Where("material_id = ?", targetID)
	default:
		db = db.Where("promotion_id = ?", targetID)
	}
	result := db.Count(&count)
	if result.Error != nil {
		return false, errors.WithStack(result.Error)
	}
	return count > 0, nil
}

// ClaimPending 原子抢占待执行记录：pending -> processing，防止并发重复执行
func (a *HostTriggerRecord) ClaimPending(ctx context.Context, id uint) (bool, error) {
	result := a.getDB(ctx).
		Where("id = ? AND execute_status = ?", id, "pending").
		Update("execute_status", "processing")
	if result.Error != nil {
		return false, errors.WithStack(result.Error)
	}
	return result.RowsAffected > 0, nil
}

// ResetStaleProcessing 重置超时的processing记录为pending（防止进程崩溃导致记录卡死）
func (a *HostTriggerRecord) ResetStaleProcessing(ctx context.Context, before time.Time) (int64, error) {
	result := a.getDB(ctx).
		Where("execute_status = ? AND updated_at < ?", "processing", before).
		Update("execute_status", "pending")
	if result.Error != nil {
		return 0, errors.WithStack(result.Error)
	}
	return result.RowsAffected, nil
}

// UpdateResult 更新执行结果
func (a *HostTriggerRecord) UpdateResult(ctx context.Context, id uint, executeStatus string, executeMsg string) error {
	result := a.getDB(ctx).Where("id = ?", id).Updates(map[string]interface{}{
		"execute_status": executeStatus,
		"execute_msg":    executeMsg,
	})
	return errors.WithStack(result.Error)
}
