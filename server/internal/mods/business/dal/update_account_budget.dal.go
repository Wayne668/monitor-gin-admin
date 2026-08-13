package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetUpdateAccountBudgetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.UpdateAccountBudget))
}

// UpdateAccountBudget 次日生效预算数据访问
type UpdateAccountBudget struct {
	DB *gorm.DB
}

// Create 新增次日生效预算记录
func (a *UpdateAccountBudget) Create(ctx context.Context, item *schema.UpdateAccountBudget) error {
	result := GetUpdateAccountBudgetDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// FindPendingNextDay 查询 is_set=0 且 budget_mod='nextDay' 的记录
func (a *UpdateAccountBudget) FindPendingNextDay(ctx context.Context) ([]*schema.UpdateAccountBudget, error) {
	var list []*schema.UpdateAccountBudget
	result := GetUpdateAccountBudgetDB(ctx, a.DB).
		Where("is_set = ? AND budget_mode = ?", 0, "nextDay").
		Find(&list)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return list, nil
}

// UpdateSetStatus 更新 is_set 和 err_msg
func (a *UpdateAccountBudget) UpdateSetStatus(ctx context.Context, id uint, isSet int8, errMsg string) error {
	result := GetUpdateAccountBudgetDB(ctx, a.DB).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_set":  isSet,
			"err_msg": errMsg,
		})
	return errors.WithStack(result.Error)
}

// QueryByAdvertiserID 根据广告主ID查询预算修改记录
func (a *UpdateAccountBudget) QueryByAdvertiserID(ctx context.Context, advertiserID int64) ([]*schema.UpdateAccountBudget, error) {
	var list []*schema.UpdateAccountBudget
	result := GetUpdateAccountBudgetDB(ctx, a.DB).
		Where("advertiser_id = ?", advertiserID).
		Order("id DESC").
		Find(&list)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return list, nil
}
