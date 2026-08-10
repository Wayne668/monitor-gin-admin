package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetHostAccountDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.HostAccount))
}

// HostAccount 托管账户数据访问
type HostAccount struct {
	DB *gorm.DB
}

// Query 查询托管账户列表
func (a *HostAccount) Query(ctx context.Context, params schema.HostAccountQueryParam, opts ...schema.HostAccountQueryOptions) (*schema.HostAccountQueryResult, error) {
	var opt schema.HostAccountQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetHostAccountDB(ctx, a.DB).Where("status = ?", 1)

	if params.AgentID > 0 {
		db = db.Where("agent_id = ?", params.AgentID)
	}
	if params.AdvertiserID > 0 {
		db = db.Where("advertiser_id = ?", params.AdvertiserID)
	}
	if v := params.AdvertiserName; len(v) > 0 {
		db = db.Where("advertiser_name LIKE ?", "%"+v+"%")
	}

	var list schema.HostAccounts
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &schema.HostAccountQueryResult{
		PageResult: pageResult,
		Data:       list,
	}, nil
}

// Get 获取指定托管账户
func (a *HostAccount) Get(ctx context.Context, id uint) (*schema.HostAccount, error) {
	item := new(schema.HostAccount)
	ok, err := util.FindOne(ctx, GetHostAccountDB(ctx, a.DB).Where("id=? AND status=?", id, 1), util.QueryOptions{}, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Create 新增托管账户
func (a *HostAccount) Create(ctx context.Context, item *schema.HostAccount) error {
	result := GetHostAccountDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update 更新托管账户
func (a *HostAccount) Update(ctx context.Context, item *schema.HostAccount) error {
	result := GetHostAccountDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete 软删除（status=0）
func (a *HostAccount) Delete(ctx context.Context, id uint) error {
	result := GetHostAccountDB(ctx, a.DB).Where("id=?", id).Update("status", 0)
	return errors.WithStack(result.Error)
}

// GetByAgentID 根据代理商ID获取账户列表
func (a *HostAccount) GetByAgentID(ctx context.Context, agentID int64) (schema.HostAccounts, error) {
	var list schema.HostAccounts
	db := GetHostAccountDB(ctx, a.DB).Where("agent_id = ? AND status = ?", agentID, 1)
	result := db.Find(&list)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return list, nil
}
