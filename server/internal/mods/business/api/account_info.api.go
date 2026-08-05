package api

import (
	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// AccountInfo 账户信息API
type AccountInfo struct {
	AccountInfoBIZ *biz.AccountInfo
}

// FindEnabledAdvertisers 查询状态为 STATUS_ENABLE 的账户列表（前端渲染用）
func (a *AccountInfo) FindEnabledAdvertisers(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := a.AccountInfoBIZ.FindEnabledAdvertisers(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, list)
}
