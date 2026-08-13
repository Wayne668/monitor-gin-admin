package schema

import "time"

// Promotion nb_promotions 表模型
type Promotion struct {
	ID            uint      `json:"id" gorm:"primarykey;column:id"`
	AdvertiserID  int64     `json:"advertiserId" gorm:"column:advertiser_id"`
	PromotionID   int64     `json:"promotionId" gorm:"column:promotion_id"`
	PromotionName string    `json:"promotionName" gorm:"column:promotion_name"`
	StatusFirst   string    `json:"statusFirst" gorm:"column:status_first"`
	StatusSecond  string    `json:"statusSecond" gorm:"column:status_second"`
	OptStatus     string    `json:"optStatus" gorm:"column:opt_status"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (Promotion) TableName() string {
	return "nb_promotions"
}

// PromotionStatusEnable 启用状态
const PromotionStatusEnable = "PROMOTION_STATUS_ENABLE"

// PromotionStatusAll 不限包含已删除
const PromotionStatusAll = "PROMOTION_STATUS_ALL"

// PromotionStatusDeleted 已删除状态
const PromotionStatusDeleted = "PROMOTION_STATUS_DELETED"

// PromotionStatusNotDelete 不限不包含已删除
const PromotionStatusNotDelete = "PROMOTION_STATUS_NOT_DELETE"

// MaterialStatusDelete 已删除状态
const MaterialStatusDelete = "MATERIAL_STATUS_DELETE"
