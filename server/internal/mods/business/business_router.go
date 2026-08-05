package business

import (
	"context"

	"monitor-gin-admin/internal/mods/business/api"
	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/dal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Business struct {
	DB             *gorm.DB
	HostRuleAPI    *api.HostRule
	HostFieldAPI   *api.HostField
	AccountInfoAPI *api.AccountInfo
	CronTab        *biz.Crontab
}

func (a *Business) Init(ctx context.Context) error {
	a.CronTab = &biz.Crontab{AccountInfo: &dal.AccountInfo{DB: a.DB}}
	return nil
}

func (a *Business) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	hostRule := v1.Group("host-rules")
	{
		hostRule.GET("", a.HostRuleAPI.Query)
		hostRule.GET(":id", a.HostRuleAPI.Get)
		hostRule.PATCH(":id/status", a.HostRuleAPI.UpdateStatus)
	}

	v1.POST("get-target-by-account", a.HostRuleAPI.GetTargetsByAccount)

	hostField := v1.Group("host-fields")
	{
		hostField.GET("", a.HostFieldAPI.Query)
		hostField.GET(":id", a.HostFieldAPI.Get)
		hostField.POST("", a.HostFieldAPI.Create)
		hostField.PUT(":id", a.HostFieldAPI.Update)
		hostField.DELETE(":id", a.HostFieldAPI.Delete)
	}

	v1.GET("account-list", a.AccountInfoAPI.FindEnabledAdvertisers)
	return nil
}

func (a *Business) Release(ctx context.Context) error {
	return nil
}
