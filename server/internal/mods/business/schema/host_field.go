package schema

import (
	"monitor-gin-admin/pkg/util"
)

var (
	HostFieldsOrderParams = []util.OrderByParam{
		{Field: "id", Direction: util.ASC},
	}
)

// HostField 托管字段模型
type HostField struct {
	ID      uint   `json:"id" gorm:"primarykey;autoIncrement"`
	Field   string `json:"field" gorm:"column:field;type:varchar(50);not null;default:''"`
	Name    string `json:"name" gorm:"column:name;type:varchar(50);not null;default:''"`
	Cate    string `json:"cate" gorm:"column:cate;type:varchar(10);not null;default:''"`
	Stash   int8   `json:"stash" gorm:"column:stash;type:tinyint(1);not null;default:0"`
	Unit    string `json:"unit" gorm:"column:unit;type:varchar(5);not null;default:''"`
	Formula string `json:"formula" gorm:"column:formula;type:varchar(255);not null;default:''"`
}

func (HostField) TableName() string {
	return "nb_host_field"
}

// HostFieldQueryParam 查询参数
type HostFieldQueryParam struct {
	util.PaginationParam
	Cate string `form:"cate"`
}

// HostFieldQueryOptions 查询选项
type HostFieldQueryOptions struct {
	util.QueryOptions
}

// HostFieldQueryResult 查询结果
type HostFieldQueryResult struct {
	Data       HostFields
	PageResult *util.PaginationResult
}

// HostFields 列表类型
type HostFields []*HostField

// HostFieldForm 新增/编辑表单
type HostFieldForm struct {
	Field   string `json:"field" binding:"required,max=50"`
	Name    string `json:"name" binding:"required,max=50"`
	Cate    string `json:"cate" binding:"required,oneof=dimension metric"`
	Stash   int8   `json:"stash" binding:"required,oneof=1 2"`
	Unit    string `json:"unit" binding:"max=5"`
	Formula string `json:"formula" binding:"max=255"`
}

func (a *HostFieldForm) Validate() error {
	return nil
}

func (a *HostFieldForm) FillTo(item *HostField) {
	item.Field = a.Field
	item.Name = a.Name
	item.Cate = a.Cate
	item.Stash = a.Stash
	item.Unit = a.Unit
	item.Formula = a.Formula
}
