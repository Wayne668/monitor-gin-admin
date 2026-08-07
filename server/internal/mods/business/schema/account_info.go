package schema

import (
	"monitor-gin-admin/internal/config"
	"time"
)

type AccountInfo struct {
	Id                 *int32     `json:"id,omitempty" form:"id" gorm:"primarykey;column:id;"`
	AdvertiserId       *int64     `json:"advertiserId,omitempty" form:"advertiserId" gorm:"comment:账户ID;column:advertiser_id;NOT NULL;"`
	AdvertiserName     *string    `json:"advertiserName,omitempty" form:"advertiserName" gorm:"comment:账户名称;column:advertiser_name;size:255;NOT NULL;DEFAULT:'';"`
	PayCompanyId       *int       `json:"payCompanyId,omitempty" form:"payCompanyId" gorm:"comment:合作主体公司;column:pay_company_id;NOT NULL;DEFAULT:0;"`
	AdvCompanyId       *int64     `json:"advCompanyId,omitempty" form:"advCompanyId" gorm:"comment:开户公司id;column:adv_company_id;NOT NULL;DEFAULT:0;"`
	AdvCompanyName     *string    `json:"advCompanyName,omitempty" form:"advCompanyName" gorm:"comment:客户公司名称;column:adv_company_name;size:255;NOT NULL;DEFAULT:'';"`
	FirstIndustryName  *string    `json:"firstIndustryName,omitempty" form:"firstIndustryName" gorm:"comment:一级行业;column:first_industry_name;size:20;NOT NULL;DEFAULT:'';"`
	SecondIndustryName *string    `json:"secondIndustryName,omitempty" form:"secondIndustryName" gorm:"comment:二级行业;column:second_industry_name;size:20;NOT NULL;DEFAULT:'';"`
	CreateTime         *time.Time `json:"createTime,omitempty" form:"createTime" gorm:"comment:投放账户创建时间;column:create_time;"`
	AdvertiserStatus   *string    `json:"advertiserStatus,omitempty" form:"advertiserStatus" gorm:"comment:投放账户状态;column:advertiser_status;size:50;"`
	OperatorTag        *int32     `json:"operatorTag,omitempty" form:"operatorTag" gorm:"comment:运营标签;column:operator_tag;"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty" form:"updatedAt" gorm:"comment:更新时间;column:updated_at;"`
}

func (a *AccountInfo) TableName() string {
	return config.C.FormatTableName("nb_account_info")
}
