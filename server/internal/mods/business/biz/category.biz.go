package biz

import (
	"context"
	"time"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/errors"
	"monitor-gin-admin/pkg/util"
)

// Category management for business
type Category struct {
	Trans       *util.Trans
	CategoryDAL *dal.Category
}

// Query categories from the data access object based on the provided parameters and options.
func (a *Category) Query(ctx context.Context, params schema.CategoryQueryParam) (*schema.CategoryQueryResult, error) {
	result, err := a.CategoryDAL.Query(ctx, params, schema.CategoryQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.CategoriesOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get the specified category from the data access object.
func (a *Category) Get(ctx context.Context, id string) (*schema.Category, error) {
	category, err := a.CategoryDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if category == nil {
		return nil, errors.NotFound("", "Category not found")
	}
	return category, nil
}

// Create a new category in the data access object.
func (a *Category) Create(ctx context.Context, formItem *schema.CategoryForm) (*schema.Category, error) {
	category := &schema.Category{
		ID:        util.NewXID(),
		CreatedAt: time.Now(),
	}

	if err := formItem.FillTo(category); err != nil {
		return nil, err
	}

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.CategoryDAL.Create(ctx, category); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return category, nil
}

// Update the specified category in the data access object.
func (a *Category) Update(ctx context.Context, id string, formItem *schema.CategoryForm) error {
	category, err := a.CategoryDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if category == nil {
		return errors.NotFound("", "Category not found")
	}

	if err := formItem.FillTo(category); err != nil {
		return err
	}
	category.UpdatedAt = time.Now()

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.CategoryDAL.Update(ctx, category); err != nil {
			return err
		}
		return nil
	})
}

// Delete the specified category from the data access object.
func (a *Category) Delete(ctx context.Context, id string) error {
	exists, err := a.CategoryDAL.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "Category not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.CategoryDAL.Delete(ctx, id); err != nil {
			return err
		}
		return nil
	})
}
