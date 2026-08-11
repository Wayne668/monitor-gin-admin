package schema

import (
	"monitor-gin-admin/pkg/util"
)

var (
	DeleteUnauditedMaterialOrderParams = []util.OrderByParam{
		{Field: "id", Direction: util.DESC},
	}
)

// DeleteUnauditedMaterial 素材删除记录
type DeleteUnauditedMaterial struct {
	ID           uint   `json:"id" gorm:"primarykey;autoIncrement;column:id"`
	MaterialID   int64  `json:"materialId" gorm:"column:material_id;type:bigint;not null;default:0"`
	AdvertiserID int64  `json:"advertiserId" gorm:"column:advertiser_id;type:bigint;not null;default:0"`
	PromotionID  int64  `json:"promotionId" gorm:"column:promotion_id;type:bigint;not null;default:0"`
	MaterialName string `json:"materialName" gorm:"column:material_name;type:varchar(100);not null;default:''"`
	IsDeleted    string `json:"isDeleted" gorm:"column:is_deleted;type:varchar(6);not null;default:'';comment:pending,deleted,failed"`
	ErrorMsg     string `json:"errorMsg" gorm:"column:error_msg;type:varchar(255);not null;default:''"`
	RetryTimes   int    `json:"retryTimes" gorm:"column:retry_times;type:tinyint(1);not null;default:0;comment:重试次数<=3"`
	CreatedAt    string `json:"createdAt" gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt    string `json:"updatedAt" gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP"`
}

func (DeleteUnauditedMaterial) TableName() string {
	return "nb_delete_unaudited_material"
}

// DeleteUnauditedMaterialQueryParam 查询参数
type DeleteUnauditedMaterialQueryParam struct {
	util.PaginationParam
	MaterialID   *int64 `form:"materialId"`
	AdvertiserID *int64 `form:"advertiserId"`
	IsDeleted    string `form:"isDeleted"`
	MaterialName string `form:"materialName"`
}

// DeleteUnauditedMaterialQueryOptions 查询选项
type DeleteUnauditedMaterialQueryOptions struct {
	util.QueryOptions
}

// DeleteUnauditedMaterialQueryResult 查询结果
type DeleteUnauditedMaterialQueryResult struct {
	Data       DeleteUnauditedMaterials
	PageResult *util.PaginationResult
}

// DeleteUnauditedMaterials 列表类型
type DeleteUnauditedMaterials []*DeleteUnauditedMaterial

// UnAudititedMaterialItem 未审核素材项
type UnAudititedMaterialItem struct {
	MaterialID   int64  `json:"materialId"`
	PromotionID  int64  `json:"promotionId"`
	AdvertiserID int64  `json:"advertiserId"`
	MaterialName string `json:"materialName"`
}

// UnAudititedMaterialReq 删除未审核素材请求
type UnAudititedMaterialReq struct {
	AccountID int64                     `json:"accountId"`
	Materials []UnAudititedMaterialItem `json:"materials"`
}
