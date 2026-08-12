package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetHostRuleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.HostRule)).Select("nb_host_rule.*, nb_user.name AS user_name").Joins("LEFT JOIN nb_user ON nb_host_rule.userid = nb_user.id")
}

// HostRule 托管规则数据访问
type HostRule struct {
	DB *gorm.DB
}

// Query 查询托管规则列表
func (a *HostRule) Query(ctx context.Context, params schema.HostRuleQueryParam, opts ...schema.HostRuleQueryOptions) (*schema.HostRuleQueryResult, error) {
	var opt schema.HostRuleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetHostRuleDB(ctx, a.DB)

	if v := params.LikeRuleName; len(v) > 0 {
		db = db.Where("rule_name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v != nil {
		db = db.Where("status = ?", *v)
	}

	var list schema.HostRules
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &schema.HostRuleQueryResult{
		PageResult: pageResult,
		Data:       list,
	}, nil
}

// Get 获取指定托管规则
func (a *HostRule) Get(ctx context.Context, id uint) (*schema.HostRule, error) {
	item := new(schema.HostRule)
	ok, err := util.FindOne(ctx, GetHostRuleDB(ctx, a.DB).Where("id=?", id), util.QueryOptions{}, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// UpdateStatus 更新托管规则状态
func (a *HostRule) UpdateStatus(ctx context.Context, id uint, status int8) error {
	result := GetHostRuleDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}

// Create 新增托管规则
func (a *HostRule) Create(ctx context.Context, item *schema.HostRule) error {
	result := GetHostRuleDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update 更新托管规则
func (a *HostRule) Update(ctx context.Context, item *schema.HostRule) error {
	result := GetHostRuleDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Updates(item)
	return errors.WithStack(result.Error)
}

// QueryAllEnabled 查询所有启用状态的托管规则
func (a *HostRule) QueryAllEnabled(ctx context.Context) ([]*schema.HostRule, error) {
	var list []*schema.HostRule
	db := GetHostRuleDB(ctx, a.DB).Where("status = ?", schema.HostRuleStatusEnabled)
	result := db.Find(&list)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return list, nil
}
