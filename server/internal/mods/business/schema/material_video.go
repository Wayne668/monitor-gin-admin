package schema

import "time"

// MaterialVideo nb_material_video 表模型
type MaterialVideo struct {
	ID           uint      `json:"id" gorm:"primarykey;column:id"`
	AdvertiserID int64     `json:"advertiserId" gorm:"column:advertiser_id"`
	MaterialID   int64     `json:"materialId" gorm:"column:material_id"`
	Signature    string    `json:"signature" gorm:"column:signature"`
	FileName     string    `json:"fileName" gorm:"column:file_name"`
	PosterURL    string    `json:"posterUrl" gorm:"column:poster_url"`
	Labels       string    `json:"labels" gorm:"column:labels"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (MaterialVideo) TableName() string {
	return "nb_material_video"
}
