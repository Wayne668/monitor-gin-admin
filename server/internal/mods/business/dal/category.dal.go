package dal

import (
	"context"

	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"
	"gorm.io/gorm"
)

// Get category storage instance
func GetCategoryDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.Category))
}

// Category management for business
type Category struct {
	DB *gorm.DB
}

// Query categories from the database based on the provided parameters and options.
func (a *Category) Query(ctx context.Context, params schema.CategoryQueryParam, opts ...schema.CategoryQueryOptions) (*schema.CategoryQueryResult, error) {
	var opt schema.CategoryQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetCategoryDB(ctx, a.DB)

	if v := params.LikeName; len(v) > 0 {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; len(v) > 0 {
		db = db.Where("status = ?", v)
	}

	var list schema.Categories
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &schema.CategoryQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified category from the database.
func (a *Category) Get(ctx context.Context, id string, opts ...schema.CategoryQueryOptions) (*schema.Category, error) {
	var opt schema.CategoryQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.Category)
	ok, err := util.FindOne(ctx, GetCategoryDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Checks if the specified category exists in the database.
func (a *Category) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetCategoryDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

// Create a new category.
func (a *Category) Create(ctx context.Context, item *schema.Category) error {
	result := GetCategoryDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified category in the database.
func (a *Category) Update(ctx context.Context, item *schema.Category) error {
	result := GetCategoryDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified category from the database.
func (a *Category) Delete(ctx context.Context, id string) error {
	result := GetCategoryDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.Category))
	return errors.WithStack(result.Error)
}
