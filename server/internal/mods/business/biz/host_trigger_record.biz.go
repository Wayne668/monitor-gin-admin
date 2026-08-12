package biz

import (
	"context"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
)

// HostTriggerRecord 触发记录业务逻辑
type HostTriggerRecord struct {
	HostTriggerRecordDAL *dal.HostTriggerRecord
}

// Query 查询触发记录列表
func (a *HostTriggerRecord) Query(ctx context.Context, params schema.HostTriggerRecordQueryParam) (*schema.HostTriggerRecordQueryResult, error) {
	return a.HostTriggerRecordDAL.Query(ctx, params)
}