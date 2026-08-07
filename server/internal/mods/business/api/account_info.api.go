package api

import (
	"strconv"
	"strings"

	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// AccountInfo 账户信息API
type AccountInfo struct {
	AccountInfoBIZ *biz.AccountInfo
}

// FindEnabledAdvertisers 查询状态为 STATUS_ENABLE 的账户列表（前端渲染用）
// Query 参数：
//   - fields: 逗号分隔的字段名，例如 "advertiser_id,advertiser_name"，为空则返回全部字段
//   - limit:  返回条数，<=0 或不传时默认返回最新 100 条（按 updated_at 倒序）
func (a *AccountInfo) FindEnabledAdvertisers(c *gin.Context) {
	ctx := c.Request.Context()

	var fields []string
	if v := c.Query("fields"); v != "" {
		for _, f := range strings.Split(v, ",") {
			if f = strings.TrimSpace(f); f != "" {
				fields = append(fields, f)
			}
		}
	}

	limit := 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	list, err := a.AccountInfoBIZ.FindEnabledAdvertisers(ctx, fields, limit)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, list)
}
