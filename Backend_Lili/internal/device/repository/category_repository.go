package repository

import (
	"Backend_Lili/internal/device/model"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

// GetAllCategories 获取所有分类
func (r *CategoryRepository) GetAllCategories() ([]*model.Category, error) {
	o := orm.NewOrm()
	var categories []*model.Category
	_, err := o.QueryTable("categories").Filter("is_active", true).OrderBy("sort_order", "created_at").All(&categories)
	return categories, err
}

// GetCategoryByID 根据ID获取分类
func (r *CategoryRepository) GetCategoryByID(categoryID int) (*model.Category, error) {
	o := orm.NewOrm()
	category := &model.Category{}
	err := o.QueryTable("categories").Filter("id", categoryID).Filter("is_active", true).One(category)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return category, err
}

// GetCategoriesByParentID 根据父分类ID获取子分类
func (r *CategoryRepository) GetCategoriesByParentID(parentID int) ([]*model.Category, error) {
	o := orm.NewOrm()
	var categories []*model.Category
	_, err := o.QueryTable("categories").
		Filter("parent_id", parentID).
		Filter("is_active", true).
		OrderBy("sort_order", "created_at").
		All(&categories)
	return categories, err
}

// CreateCategory 创建分类
func (r *CategoryRepository) CreateCategory(category *model.Category) error {
	o := orm.NewOrm()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	_, err := o.Insert(category)
	return err
}

// UpdateCategory 更新分类
func (r *CategoryRepository) UpdateCategory(category *model.Category) error {
	o := orm.NewOrm()
	category.UpdatedAt = time.Now()
	_, err := o.Update(category)
	return err
} 