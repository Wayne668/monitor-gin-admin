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
func (a *AccountInfo) FindEnabledAdvertisers(ctx context.Context) ([]schema.AccountInfo, error) {
	return a.AccountInfoDAL.FindEnabledAdvertisers(ctx)
}
