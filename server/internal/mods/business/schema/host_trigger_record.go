package schema

import (
	"time"

	"monitor-gin-admin/pkg/util"
)

// HostTriggerRecord 触发记录模型
type HostTriggerRecord struct {
	ID            uint      `json:"id" gorm:"primarykey;autoIncrement"`
	RuleID        int       `json:"rule_id" gorm:"column:rule_id;type:int;not null;default:0"`
	AdvertiserID  int64     `json:"advertiser_id" gorm:"column:advertiser_id;type:bigint;not null;default:0"`
	PromotionID   int64     `json:"promotion_id" gorm:"column:promotion_id;type:bigint;not null;default:0"`
	MaterialID    int64     `json:"material_id" gorm:"column:material_id;type:bigint;not null;default:0"`
	Target        string    `json:"target" gorm:"column:target;type:varchar(10);not null;default:''"`
	ExecuteAction string    `json:"execute_action" gorm:"column:execute_action;type:varchar(10);not null;default:''"`
	ExecuteStatus string    `json:"execute_status" gorm:"column:execute_status;type:varchar(10);not null;default:''"`
	ExecuteMsg    string    `json:"execute_msg" gorm:"column:execute_msg;type:varchar(255);not null;default:''"`
	TriggerReason string    `json:"trigger_reason" gorm:"column:trigger_reason;type:varchar(255);not null;default:''"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (HostTriggerRecord) TableName() string {
	return "nb_host_trigger_record"
}

// HostTriggerRecords 列表类型
type HostTriggerRecords []*HostTriggerRecord

// HostTriggerRecordQueryParam 查询参数
type HostTriggerRecordQueryParam struct {
	util.PaginationParam
	RuleID        *int   `form:"ruleId"`
	AdvertiserID  *int64 `form:"advertiserId"`
	Target        string `form:"target"`
	ExecuteStatus string `form:"executeStatus"`
}

// HostTriggerRecordQueryResult 查询结果
type HostTriggerRecordQueryResult struct {
	Data       HostTriggerRecords
	PageResult *util.PaginationResult
}
