package biz

import (
	"context"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
)

// AccountInfo 账户信息业务逻辑
type AccountInfo struct {
	AccountInfoDAL *dal.AccountInfo
}

// FindEnabledAdvertisers 查询状态为 STATUS_ENABLE 的账户列表
// fields 指定查询字段，为空则查询全部字段
// limit 限制返回条数，<=0 时由 DAL 层默认返回最新 100 条
func (a *AccountInfo) FindEnabledAdvertisers(ctx context.Context, fields []string, limit int) ([]schema.AccountInfo, error) {
	return a.AccountInfoDAL.FindEnabledAdvertisers(ctx, fields, limit)
}
