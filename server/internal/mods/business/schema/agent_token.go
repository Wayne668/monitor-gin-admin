package schema

import (
	"monitor-gin-admin/pkg/util"
)

var (
	AgentTokensOrderParams = []util.OrderByParam{
		{Field: "id", Direction: util.DESC},
	}
)

// AgentToken 代理商账号授权&accesstoken
type AgentToken struct {
	ID           uint   `json:"id" gorm:"primarykey;autoIncrement;column:id"`
	AccountName  string `json:"accountName" gorm:"column:account_name;type:varchar(255);comment:账号名"`
	AccountID    string `json:"accountId" gorm:"column:account_id;type:varchar(50);comment:账号id"`
	AuthStatus   string `json:"authStatus" gorm:"column:authstatus;type:varchar(255);comment:是否授权成功"`
	AccessToken  string `json:"accessToken" gorm:"column:accesstoken;type:varchar(255)"`
	RefreshToken string `json:"refreshToken" gorm:"column:refreshtoken;type:varchar(255);default:''"`
	TokenTime    int64  `json:"tokenTime" gorm:"column:tokentime;comment:token更新时间"`
	Remarks      string `json:"remarks" gorm:"column:remarks;type:varchar(255);default:'';comment:备注"`
	AppID        string `json:"appId" gorm:"column:app_id;type:varchar(20);not null;default:''"`
	AppSecret    string `json:"appSecret" gorm:"column:app_secret;type:varchar(50);not null;default:''"`
	AppName      string `json:"appName" gorm:"column:app_name;type:varchar(255);not null;default:''"`
	CreatedAt    string `json:"createdAt" gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt    string `json:"updatedAt" gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP"`
}

func (AgentToken) TableName() string {
	return "nb_agent_token"
}

// AgentTokenQueryParam 查询参数
type AgentTokenQueryParam struct {
	util.PaginationParam
	AccountName string `form:"accountName"`
	AccountID   string `form:"accountId"`
}

// AgentTokenQueryOptions 查询选项
type AgentTokenQueryOptions struct {
	util.QueryOptions
}

// AgentTokenQueryResult 查询结果
type AgentTokenQueryResult struct {
	Data       AgentTokens
	PageResult *util.PaginationResult
}

// AgentTokens 列表类型
type AgentTokens []*AgentToken

// AgentTokenForm 新增/编辑表单
type AgentTokenForm struct {
	AccountName  string `json:"accountName" binding:"max=255"`
	AccountID    string `json:"accountId" binding:"required,max=50"`
	AuthStatus   string `json:"authStatus" binding:"max=255"`
	AccessToken  string `json:"accessToken" binding:"max=255"`
	RefreshToken string `json:"refreshToken" binding:"max=255"`
	TokenTime    int64  `json:"tokenTime"`
	Remarks      string `json:"remarks" binding:"max=255"`
	AppID        string `json:"appId" binding:"max=20"`
	AppSecret    string `json:"appSecret" binding:"max=50"`
	AppName      string `json:"appName" binding:"max=255"`
}

func (a *AgentTokenForm) Validate() error {
	return nil
}

func (a *AgentTokenForm) FillTo(item *AgentToken) {
	item.AccountName = a.AccountName
	item.AccountID = a.AccountID
	item.AuthStatus = a.AuthStatus
	item.AccessToken = a.AccessToken
	item.RefreshToken = a.RefreshToken
	item.TokenTime = a.TokenTime
	item.Remarks = a.Remarks
	item.AppID = a.AppID
	item.AppSecret = a.AppSecret
	item.AppName = a.AppName
}
