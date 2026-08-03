package schema

import (
	"time"

	"monitor-gin-admin/internal/config"
	"monitor-gin-admin/pkg/util"
)

const (
	CategoryStatusDisabled = "disabled"
	CategoryStatusEnabled  = "enabled"
)

var (
	CategoriesOrderParams = []util.OrderByParam{
		{Field: "created_at", Direction: util.DESC},
	}
)

// Category management for business
type Category struct {
	ID          string    `json:"id" gorm:"size:20;primarykey;"` // Unique ID
	Name        string    `json:"name" gorm:"size:128;index"`    // Display name of category
	Description string    `json:"description" gorm:"size:1024"`  // Details about category
	Status      string    `json:"status" gorm:"size:20;index"`   // Status of category (enabled, disabled)
	CreatedAt   time.Time `json:"created_at" gorm:"index;"`      // Create time
	UpdatedAt   time.Time `json:"updated_at" gorm:"index;"`      // Update time
}

func (a *Category) TableName() string {
	return config.C.FormatTableName("category")
}

// Defining the query parameters for the `Category` struct.
type CategoryQueryParam struct {
	util.PaginationParam
	LikeName string `form:"name"`   // Display name of category
	Status   string `form:"status"` // Status of category (disabled, enabled)
}

// Defining the query options for the `Category` struct.
type CategoryQueryOptions struct {
	util.QueryOptions
}

// Defining the query result for the `Category` struct.
type CategoryQueryResult struct {
	Data       Categories
	PageResult *util.PaginationResult
}

// Defining the slice of `Category` struct.
type Categories []*Category

// Defining the data structure for creating a `Category` struct.
type CategoryForm struct {
	Name        string `json:"name" binding:"required,max=128"`                  // Display name of category
	Description string `json:"description"`                                      // Details about category
	Status      string `json:"status" binding:"required,oneof=disabled enabled"` // Status of category (enabled, disabled)
}

// A validation function for the `CategoryForm` struct.
func (a *CategoryForm) Validate() error {
	return nil
}

func (a *CategoryForm) FillTo(category *Category) error {
	category.Name = a.Name
	category.Description = a.Description
	category.Status = a.Status
	return nil
}
