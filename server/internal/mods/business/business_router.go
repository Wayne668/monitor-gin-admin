package business

import (
	"context"

	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/dal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Business struct {
	DB      *gorm.DB
	CronTab *biz.CronTab
}

func (a *Business) Init(ctx context.Context) error {
	a.CronTab = &biz.CronTab{Acc: &dal.AccountInfo{DB: a.DB}}
	return nil
}

func (a *Business) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	return nil
}

func (a *Business) Release(ctx context.Context) error {
	return nil
}
