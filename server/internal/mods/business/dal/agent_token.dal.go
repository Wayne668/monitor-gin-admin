package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetAgentTokenDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.AgentToken))
}

// AgentToken 代理商账号授权数据访问
type AgentToken struct {
	DB *gorm.DB
}

// Query 查询代理商账号授权列表
func (a *AgentToken) Query(ctx context.Context, params schema.AgentTokenQueryParam, opts ...schema.AgentTokenQueryOptions) (*schema.AgentTokenQueryResult, error) {
	var opt schema.AgentTokenQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetAgentTokenDB(ctx, a.DB)

	if v := params.AccountName; len(v) > 0 {
		db = db.Where("account_name LIKE ?", "%"+v+"%")
	}
	if v := params.AccountID; len(v) > 0 {
		db = db.Where("account_id = ?", v)
	}

	var list schema.AgentTokens
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &schema.AgentTokenQueryResult{
		PageResult: pageResult,
		Data:       list,
	}, nil
}

// Get 获取指定代理商账号授权
func (a *AgentToken) Get(ctx context.Context, id uint) (*schema.AgentToken, error) {
	item := new(schema.AgentToken)
	ok, err := util.FindOne(ctx, GetAgentTokenDB(ctx, a.DB).Where("id=?", id), util.QueryOptions{}, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Create 新增代理商账号授权
func (a *AgentToken) Create(ctx context.Context, item *schema.AgentToken) error {
	result := GetAgentTokenDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update 更新代理商账号授权
func (a *AgentToken) Update(ctx context.Context, item *schema.AgentToken) error {
	result := GetAgentTokenDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete 删除代理商账号授权
func (a *AgentToken) Delete(ctx context.Context, id uint) error {
	result := GetAgentTokenDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.AgentToken))
	return errors.WithStack(result.Error)
}

// GetByAccountID 根据账号ID获取token
func (a *AgentToken) GetByAccountID(ctx context.Context, accountID string) (*schema.AgentToken, error) {
	item := new(schema.AgentToken)
	ok, err := util.FindOne(ctx, GetAgentTokenDB(ctx, a.DB).Where("account_id=?", accountID), util.QueryOptions{}, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}
