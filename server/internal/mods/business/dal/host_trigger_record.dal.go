package dal

import (
	"context"

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

// UpdateResult 更新执行结果
func (a *HostTriggerRecord) UpdateResult(ctx context.Context, id uint, executeStatus string, executeMsg string) error {
	result := a.getDB(ctx).Where("id = ?", id).Updates(map[string]interface{}{
		"execute_status": executeStatus,
		"execute_msg":    executeMsg,
	})
	return errors.WithStack(result.Error)
}
