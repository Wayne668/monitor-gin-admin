package schema

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"monitor-gin-admin/pkg/util"
)

const (
	HostRuleStatusEnabled = 1
	HostRuleStatusPaused  = 2
	HostRuleStatusDeleted = 3
)

var (
	HostRulesOrderParams = []util.OrderByParam{
		{Field: "id", Direction: util.DESC},
	}
)

// HostRule 托管规则模型
type HostRule struct {
	ID                 uint      `json:"id" gorm:"primarykey;autoIncrement"`
	Target             string    `json:"target" gorm:"column:target;type:varchar(20);not null;default:''"`
	TargetAccounts     string    `json:"target_accounts" gorm:"column:target_accounts;type:text;not null"`
	TargetPromotion    string    `json:"target_promotion" gorm:"column:target_promotion;type:text;not null"`
	TargetProjects     string    `json:"target_projects" gorm:"column:target_projects;type:text;not null"`
	TargetMaterial     string    `json:"target_material" gorm:"column:target_material;type:text;not null"`
	AgentID            int64     `json:"agent_id" gorm:"column:agent_id;type:bigint;not null;default:0"`
	TriggerCondition   string    `json:"trigger_condition" gorm:"column:trigger_condition;type:text;not null"`
	ExecuteAction      string    `json:"execute_action" gorm:"column:execute_action;type:varchar(255);not null;default:''"`
	OperateMethod      int8      `json:"operate_method" gorm:"column:operate_method;type:tinyint(1);not null;default:0"`
	TriggerStartDate   time.Time `json:"trigger_start_date" gorm:"column:trigger_start_date;type:datetime;not null"`
	TriggerEndDate     time.Time `json:"trigger_end_date" gorm:"column:trigger_end_date;type:datetime;not null"`
	TriggerFrequency   int       `json:"trigger_frequency" gorm:"column:trigger_frequency;type:int;not null;default:0"`
	NotifyFrequency    int       `json:"notify_frequency" gorm:"column:notify_frequency;type:int;not null;default:0"`
	RuleName           string    `json:"rule_name" gorm:"column:rule_name;type:varchar(50);not null;default:''"`
	UserID             int       `json:"userid" gorm:"column:userid;type:int;not null;default:0"`
	UserName           string    `json:"user_name" gorm:"->;-:migration"`
	DingtalkWebhookUrl string    `json:"dingtalk_webhook_url" gorm:"column:dingtalk_webhook_url;type:varchar(500);not null;default:''"`
	Status             int8      `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (HostRule) TableName() string {
	return "nb_host_rule"
}

// HostRuleQueryParam 查询参数
type HostRuleQueryParam struct {
	util.PaginationParam
	LikeRuleName string `form:"ruleName"`
	Status       *int8  `form:"status"`
}

// HostRuleQueryOptions 查询选项
type HostRuleQueryOptions struct {
	util.QueryOptions
}

// HostRuleQueryResult 查询结果
type HostRuleQueryResult struct {
	Data       HostRules
	PageResult *util.PaginationResult
}

// HostRules 列表类型
type HostRules []*HostRule

// HostRuleUpdateStatusForm 修改状态表单
type HostRuleUpdateStatusForm struct {
	Status int8 `json:"status" binding:"required,oneof=1 2 3"`
}

func (a *HostRuleUpdateStatusForm) Validate() error {
	return nil
}

// TargetByAccountReq 根据账户获取目标列表请求
type TargetByAccountReq struct {
	Target     string   `json:"target" binding:"required,oneof=promotion creative"`
	AccountIDs []string `json:"accountIds" binding:"required,min=1"`
	AgentID    string   `json:"agentId" binding:"required"`
}

func (a *TargetByAccountReq) Validate() error {
	return nil
}

// TargetItem 目标项
type TargetItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// HostRuleForm 新增/编辑托管规则表单
type HostRuleForm struct {
	ID                 uint        `json:"id"`
	RuleName           string      `json:"ruleName" binding:"required,max=50"`
	Target             string      `json:"target" binding:"required,oneof=account project promotion creative"`
	ScopeType          string      `json:"scopeType"`
	SelectedAgentID    string      `json:"selectedAgentId"` // nb_agent_token.account_id
	SelectedAccountIds []string    `json:"selectedAccountIds"`
	SelectedTargetIds  []string    `json:"selectedTargetIds"`
	ConditionConfig    interface{} `json:"conditionConfig"`
	Action             string      `json:"action"`
	CheckFreq          int         `json:"checkFreq"`
	DateRange          []string    `json:"dateRange"`
	NotifyMethods      []string    `json:"notifyMethods"`
	DingtalkWebhookUrl string      `json:"dingtalkWebhookUrl"`
}

func (a *HostRuleForm) Validate() error {
	return nil
}

func (a *HostRuleForm) FillTo(item *HostRule) {
	item.RuleName = a.RuleName
	item.Target = a.Target
	item.AgentID, _ = strconv.ParseInt(a.SelectedAgentID, 10, 64)
	// 使用条件与操作 JSON 序列化，不转义 HTML 字符
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(a.ConditionConfig)
	item.TriggerCondition = string(bytes.TrimRight(buf.Bytes(), "\n"))
	item.ExecuteAction = a.Action
	// logic 映射到 operate_method：and=0, or=1
	if cfg, ok := a.ConditionConfig.(map[string]interface{}); ok {
		if logic, ok := cfg["logic"].(string); ok && logic == "or" {
			item.OperateMethod = 1
		}
	}
	item.TriggerFrequency = a.CheckFreq
	item.DingtalkWebhookUrl = a.DingtalkWebhookUrl
	// 生效日期
	if len(a.DateRange) == 2 {
		start, _ := time.Parse("2006-01-02", a.DateRange[0])
		end, _ := time.Parse("2006-01-02", a.DateRange[1])
		item.TriggerStartDate = start
		item.TriggerEndDate = end
	}
	// 账户列表
	accountsJSON, _ := json.Marshal(a.SelectedAccountIds)
	item.TargetAccounts = string(accountsJSON)
	// 目标列表（按 target 类型写入对应字段）
	targetJSON, _ := json.Marshal(a.SelectedTargetIds)
	switch a.Target {
	case "promotion":
		item.TargetPromotion = string(targetJSON)
	case "creative":
		item.TargetMaterial = string(targetJSON)
	case "project":
		item.TargetProjects = string(targetJSON)
	}
	// 通知方式映射到 NotifyFrequency
	notifyVal := 0
	for _, m := range a.NotifyMethods {
		switch m {
		case "sms":
			notifyVal |= 1
		case "dingtalk":
			notifyVal |= 2
		}
	}
	item.NotifyFrequency = notifyVal
}
