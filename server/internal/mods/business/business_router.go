package business

import (
	"context"

	"monitor-gin-admin/internal/mods/business/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Business struct {
	DB                         *gorm.DB
	HostRuleAPI                *api.HostRule
	HostFieldAPI               *api.HostField
	AccountInfoAPI             *api.AccountInfo
	AgentTokenAPI              *api.AgentToken
	DeleteUnauditedMaterialAPI *api.DeleteUnauditedMaterial
	CrontabAPI                 *api.Crontab
}

func (a *Business) Init(ctx context.Context) error {
	return nil
}

func (a *Business) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	crontab := v1.Group("crontab")
	{
		crontab.POST("refresh-token", a.CrontabAPI.RefreshToken)
		crontab.POST("handle-host-rule", a.CrontabAPI.HandleHostRule)
	}

	hostRule := v1.Group("host-rules")
	{
		hostRule.GET("", a.HostRuleAPI.Query)
		hostRule.GET(":id", a.HostRuleAPI.Get)
		hostRule.POST("", a.HostRuleAPI.SaveHostRule)
		hostRule.PATCH(":id/status", a.HostRuleAPI.UpdateStatus)
	}

	hostField := v1.Group("host-fields")
	{
		hostField.GET("", a.HostFieldAPI.Query)
		hostField.GET(":id", a.HostFieldAPI.Get)
		hostField.POST("", a.HostFieldAPI.Create)
		hostField.PUT(":id", a.HostFieldAPI.Update)
		hostField.DELETE(":id", a.HostFieldAPI.Delete)
	}

	agentToken := v1.Group("agent-tokens")
	{
		agentToken.GET("", a.AgentTokenAPI.Query)
		agentToken.GET(":id", a.AgentTokenAPI.Get)
		agentToken.POST("", a.AgentTokenAPI.Create)
		agentToken.PUT(":id", a.AgentTokenAPI.Update)
		agentToken.DELETE(":id", a.AgentTokenAPI.Delete)
	}

	{
		v1.GET("account-list", a.AccountInfoAPI.FindEnabledAdvertisers)
		v1.POST("get-target-by-account", a.HostRuleAPI.GetTargetsByAccount)
		v1.GET("get-unaudited-material", a.DeleteUnauditedMaterialAPI.GetUnAudititedMaterial)
		v1.POST("delete-unaudited-material", a.DeleteUnauditedMaterialAPI.DeleteUnAudititedMaterial)
		v1.POST("retry-failed-delete", a.DeleteUnauditedMaterialAPI.RetryFailedDelete)
	}

	deleteMaterial := v1.Group("delete-unaudited-material-records")
	{
		deleteMaterial.GET("", a.DeleteUnauditedMaterialAPI.Query)
		deleteMaterial.GET(":id", a.DeleteUnauditedMaterialAPI.Get)
	}

	return nil
}

func (a *Business) Release(ctx context.Context) error {
	return nil
}
