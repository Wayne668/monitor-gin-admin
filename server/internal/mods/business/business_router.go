package business

import (
	"context"

	"monitor-gin-admin/internal/config"
	"monitor-gin-admin/internal/mods/business/api"
	"monitor-gin-admin/internal/mods/business/schema"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Business struct {
	DB          *gorm.DB
	CategoryAPI *api.Category
}

func (a *Business) AutoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(
		new(schema.Category),
	)
}

func (a *Business) Init(ctx context.Context) error {
	if config.C.Storage.DB.AutoMigrate {
		if err := a.AutoMigrate(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *Business) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	category := v1.Group("categories")
	{
		category.GET("", a.CategoryAPI.Query)
		category.GET(":id", a.CategoryAPI.Get)
		category.POST("", a.CategoryAPI.Create)
		category.PUT(":id", a.CategoryAPI.Update)
		category.DELETE(":id", a.CategoryAPI.Delete)
	}
	return nil
}

func (a *Business) Release(ctx context.Context) error {
	return nil
}
