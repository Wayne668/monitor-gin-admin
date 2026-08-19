package schema

import "time"

// PromotionMaterial nb_promotion_material 表模型
type PromotionMaterial struct {
	ID             uint      `json:"id" gorm:"primarykey;column:id"`
	AdvertiserID   int64     `json:"advertiserId" gorm:"column:advertiser_id"`
	PromotionID    int64     `json:"promotionId" gorm:"column:promotion_id"`
	PromotionName  string    `json:"promotionName" gorm:"column:promotion_name"`
	StatusFirst    string    `json:"statusFirst" gorm:"column:status_first"`
	OptStatus      string    `json:"optStatus" gorm:"column:opt_status"`
	MaterialID     int64     `json:"materialId" gorm:"column:material_id"`
	MaterialStatus string    `json:"materialStatus" gorm:"column:material_status"`
	FileName       string    `json:"fileName" gorm:"->"` // JOIN nb_material_video 获取
	CreatedAt      time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (PromotionMaterial) TableName() string {
	return "nb_promotion_material"
}

// MaterialStatusOK 素材状态正常
const MaterialStatusOK = "MATERIAL_STATUS_OK"

// MaterialStatusOfflineAudit 素材状态：未审核下线
const MaterialStatusOfflineAudit = "MATERIAL_STATUS_OFFLINE_AUDIT"

// MaterialStatusDelete 已删除状态
const MaterialStatusDelete = "MATERIAL_STATUS_DELETE"
