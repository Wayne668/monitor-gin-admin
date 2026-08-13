package schema

import "time"

// UpdateAccountBudget 次日生效预算记录
type UpdateAccountBudget struct {
	ID           uint      `json:"id" gorm:"primarykey;autoIncrement;column:id"`
	AdvertiserID int64     `json:"advertiserId" gorm:"column:advertiser_id;NOT NULL;DEFAULT:0"`
	Budget       float64   `json:"budget" gorm:"column:budget;type:decimal(10,2);NOT NULL;DEFAULT:0.00"`
	BudgetMode   string    `json:"budgetMode" gorm:"column:budget_mod;size:10;NOT NULL;DEFAULT:''"`
	IsSet        int8      `json:"isSet" gorm:"column:is_set;NOT NULL;DEFAULT:0"`
	ErrMsg       string    `json:"errMsg" gorm:"column:err_msg;size:255;NOT NULL;DEFAULT:''"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;NOT NULL;DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;NOT NULL;DEFAULT:CURRENT_TIMESTAMP"`
}

func (UpdateAccountBudget) TableName() string {
	return "nb_update_account_budget"
}

// UpdateAccountBudgetForm 新增次日生效预算表单
type UpdateAccountBudgetForm struct {
	AdvertiserID int64   `json:"advertiserId" binding:"required"`
	Budget       float64 `json:"budget" binding:"required"`
	BudgetMode   string  `json:"budgetMode"`
}

func (a *UpdateAccountBudgetForm) Validate() error {
	return nil
}

func (a *UpdateAccountBudgetForm) FillTo(item *UpdateAccountBudget) {
	item.AdvertiserID = a.AdvertiserID
	item.Budget = a.Budget
	item.BudgetMode = a.BudgetMode
}

// AdvertiserBudgetUpdateForm 立即生效预算更新表单
type AdvertiserBudgetUpdateForm struct {
	AccountID    int64   `json:"accountId" binding:"required"`
	AdvertiserID int64   `json:"advertiserId" binding:"required"`
	BudgetMode   string  `json:"budgetMode" binding:"required"`
	Budget       float64 `json:"budget"`
}

func (a *AdvertiserBudgetUpdateForm) Validate() error {
	return nil
}
