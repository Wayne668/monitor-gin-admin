package api

import (
	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type Crontab struct {
	CrontabBIZ *biz.Crontab
}

func (t *Crontab) RefreshToken(c *gin.Context) {
	if err := t.CrontabBIZ.RefreshToken(); err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}

func (t *Crontab) HandleHostRule(c *gin.Context) {
	if err := t.CrontabBIZ.HandleHostRule(); err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}

func (t *Crontab) SyncPromotionMaterial(c *gin.Context) {
	if err := t.CrontabBIZ.SyncPromotionMaterial(); err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}
