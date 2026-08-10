package schema

import (
	"monitor-gin-admin/pkg/util"
)

var (
	HostAccountOrderParams = []util.OrderByParam{
		{Field: "id", Direction: util.DESC},
	}
)

// HostAccount 托管账户
type HostAccount struct {
	ID             uint   `json:"id" gorm:"primarykey;autoIncrement;column:id"`
	AgentID        int64  `json:"agentId" gorm:"column:agent_id;NOT NULL;DEFAULT:0"`
	AdvertiserID   int64  `json:"advertiserId" gorm:"column:advertiser_id;NOT NULL;DEFAULT:0"`
	AdvertiserName string `json:"advertiserName" gorm:"column:advertiser_name;size:100;NOT NULL;DEFAULT:''"`
	Status         int8   `json:"status" gorm:"column:status;NOT NULL;DEFAULT:0"`
}

func (HostAccount) TableName() string {
	return "nb_host_account"
}

// HostAccountQueryParam 查询参数
type HostAccountQueryParam struct {
	util.PaginationParam
	AgentID        int64  `form:"agentId"`
	AdvertiserID   int64  `form:"advertiserId"`
	AdvertiserName string `form:"advertiserName"`
}

// HostAccountQueryOptions 查询选项
type HostAccountQueryOptions struct {
	util.QueryOptions
}

// HostAccountQueryResult 查询结果
type HostAccountQueryResult struct {
	Data       HostAccounts
	PageResult *util.PaginationResult
}

// HostAccounts 列表类型
type HostAccounts []*HostAccount

// HostAccountForm 新增/编辑表单
type HostAccountForm struct {
	AgentID        int64  `json:"agentId" binding:"required"`
	AdvertiserID   int64  `json:"advertiserId" binding:"required"`
	AdvertiserName string `json:"advertiserName" binding:"required,max=100"`
}

func (a *HostAccountForm) Validate() error {
	return nil
}

func (a *HostAccountForm) FillTo(item *HostAccount) {
	item.AgentID = a.AgentID
	item.AdvertiserID = a.AdvertiserID
	item.AdvertiserName = a.AdvertiserName
}
