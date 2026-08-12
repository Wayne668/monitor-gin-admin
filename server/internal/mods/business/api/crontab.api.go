package api

import (
	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/pkg/util"
	"strconv"

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
	agentIDStr := c.Query("agent_id")
	advertiserIDStr := c.Query("advertiser_id")

	var agentID, advertiserID int64
	if agentIDStr != "" {
		var err error
		agentID, err = strconv.ParseInt(agentIDStr, 10, 64)
		if err != nil {
			util.ResError(c, err)
			return
		}
	}
	if advertiserIDStr != "" {
		var err error
		advertiserID, err = strconv.ParseInt(advertiserIDStr, 10, 64)
		if err != nil {
			util.ResError(c, err)
			return
		}
	}

	if err := t.CrontabBIZ.SyncPromotionMaterial(agentID, advertiserID); err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}

func (t *Crontab) SyncHostTriggerRecord(c *gin.Context) {
	if err := t.CrontabBIZ.SyncHostTriggerRecord(); err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}
