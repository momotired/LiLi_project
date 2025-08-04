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

// GetCategoriesByType 根据类型获取分类
func (r *CategoryRepository) GetCategoriesByType(categoryType string, userID int, includeCount bool) ([]*model.Category, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("categories").Filter("is_active", true).Filter("deleted_at__isnull", true)

	switch categoryType {
	case "system":
		qs = qs.Filter("type", "system")
	case "custom":
		qs = qs.Filter("type", "custom").Filter("user_id", userID)
	case "all":
		// 系统分类或用户自定义分类 - 需要分别查询后合并
		var systemCategories []*model.Category
		var customCategories []*model.Category

		// 查询系统分类
		_, err := qs.Filter("type", "system").All(&systemCategories)
		if err != nil {
			return nil, err
		}

		// 查询用户自定义分类
		_, err = qs.Filter("type", "custom").Filter("user_id", userID).All(&customCategories)
		if err != nil {
			return nil, err
		}

		// 合并结果
		allCategories := append(systemCategories, customCategories...)

		// 如果需要包含设备数量
		if includeCount {
			for _, category := range allCategories {
				count, err := r.GetDeviceCountByCategory(category.ID, userID)
				if err == nil {
					category.DeviceCount = count
				}
			}
		}

		return allCategories, nil
	}

	var categories []*model.Category
	_, err := qs.OrderBy("sort_order", "created_at").All(&categories)
	if err != nil {
		return nil, err
	}

	// 如果需要包含设备数量
	if includeCount {
		for _, category := range categories {
			count, err := r.GetDeviceCountByCategory(category.ID, userID)
			if err == nil {
				category.DeviceCount = count
			}
		}
	}

	return categories, nil
}

// GetDeviceCountByCategory 获取分类下的设备数量
func (r *CategoryRepository) GetDeviceCountByCategory(categoryID, userID int) (int, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable("devices").
		Filter("category_id", categoryID).
		Filter("user_id", userID).
		Filter("deleted_at__isnull", true).
		Count()
	return int(count), err
}

// GetSystemCategories 获取系统默认分类
func (r *CategoryRepository) GetSystemCategories() ([]*model.Category, error) {
	o := orm.NewOrm()
	var categories []*model.Category
	_, err := o.QueryTable("categories").
		Filter("type", "system").
		Filter("is_active", true).
		Filter("deleted_at__isnull", true).
		OrderBy("sort_order", "created_at").
		All(&categories)
	return categories, err
}

// GetCustomCategories 获取用户自定义分类
func (r *CategoryRepository) GetCustomCategories(userID int) ([]*model.Category, error) {
	o := orm.NewOrm()
	var categories []*model.Category
	_, err := o.QueryTable("categories").
		Filter("type", "custom").
		Filter("user_id", userID).
		Filter("is_active", true).
		Filter("deleted_at__isnull", true).
		OrderBy("sort_order", "created_at").
		All(&categories)
	return categories, err
}

// DeleteCategory 软删除分类
func (r *CategoryRepository) DeleteCategory(categoryID int, userID int) error {
	o := orm.NewOrm()

	// 只能删除自定义分类
	category := &model.Category{}
	err := o.QueryTable("categories").
		Filter("id", categoryID).
		Filter("type", "custom").
		Filter("user_id", userID).
		One(category)
	if err != nil {
		return err
	}

	category.DeletedAt = time.Now()
	category.IsActive = false
	_, err = o.Update(category, "deleted_at", "is_active")
	return err
}

// UpdateCategoriesSort 批量更新分类排序
func (r *CategoryRepository) UpdateCategoriesSort(userID int, categoryOrders []map[string]int) error {
	o := orm.NewOrm()

	for _, order := range categoryOrders {
		categoryID := order["category_id"]
		sortOrder := order["sort_order"]

		// 验证分类归属权
		category := &model.Category{}
		err := o.QueryTable("categories").
			Filter("id", categoryID).
			Filter("user_id", userID).
			Filter("type", "custom").
			One(category)
		if err != nil {
			return err
		}

		// 更新排序
		category.SortOrder = sortOrder
		category.UpdatedAt = time.Now()
		_, err = o.Update(category, "sort_order", "updated_at")
		if err != nil {
			return err
		}
	}

	return nil
}

// SearchCategories 搜索分类
func (r *CategoryRepository) SearchCategories(keyword, categoryType string, userID int) ([]*model.Category, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("categories").
		Filter("is_active", true).
		Filter("deleted_at__isnull", true).
		Filter("name__icontains", keyword)

	switch categoryType {
	case "system":
		qs = qs.Filter("type", "system")
	case "custom":
		qs = qs.Filter("type", "custom").Filter("user_id", userID)
	case "all":
		// 查询系统分类和用户自定义分类需要分别处理
		var systemCategories []*model.Category
		var customCategories []*model.Category

		// 查询系统分类
		_, err := o.QueryTable("categories").
			Filter("is_active", true).
			Filter("deleted_at__isnull", true).
			Filter("type", "system").
			Filter("name__icontains", keyword).
			OrderBy("sort_order", "name").
			All(&systemCategories)
		if err != nil {
			return nil, err
		}

		// 查询用户自定义分类
		_, err = o.QueryTable("categories").
			Filter("is_active", true).
			Filter("deleted_at__isnull", true).
			Filter("type", "custom").
			Filter("user_id", userID).
			Filter("name__icontains", keyword).
			OrderBy("sort_order", "name").
			All(&customCategories)
		if err != nil {
			return nil, err
		}

		// 合并结果
		allCategories := append(systemCategories, customCategories...)
		return allCategories, nil
	}

	var categories []*model.Category
	_, err := qs.OrderBy("sort_order", "name").All(&categories)
	return categories, err
}

// GetCategoryStatistics 获取分类统计信息
func (r *CategoryRepository) GetCategoryStatistics(userID int, period string) (map[string]interface{}, error) {
	o := orm.NewOrm()
	stats := make(map[string]interface{})

	// 获取各分类的设备数量和总价值
	var categoryStats []orm.Params
	sql := `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COUNT(d.id) as device_count,
			COALESCE(SUM(d.current_value), 0) as total_value,
			COALESCE(AVG(d.current_value), 0) as avg_value
		FROM categories c
		LEFT JOIN devices d ON c.id = d.category_id AND d.user_id = ? AND d.deleted_at IS NULL
		WHERE c.is_active = 1 AND c.deleted_at IS NULL 
		AND (c.type = 'system' OR (c.type = 'custom' AND c.user_id = ?))
		GROUP BY c.id, c.name
		ORDER BY device_count DESC
	`
	_, err := o.Raw(sql, userID, userID).Values(&categoryStats)
	if err != nil {
		return nil, err
	}

	stats["categories"] = categoryStats

	// 获取总设备数和总价值
	var totalStats orm.Params
	totalSql := `
		SELECT 
			COUNT(*) as total_devices,
			COALESCE(SUM(current_value), 0) as total_value
		FROM devices 
		WHERE user_id = ? AND deleted_at IS NULL
	`
	err = o.Raw(totalSql, userID).QueryRow(&totalStats)
	if err != nil {
		return nil, err
	}

	stats["total_devices"] = totalStats["total_devices"]
	stats["total_value"] = totalStats["total_value"]

	return stats, nil
}

// CheckCategoryExists 检查分类是否存在
func (r *CategoryRepository) CheckCategoryExists(name string, userID int, excludeID int) (bool, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("categories").
		Filter("name", name).
		Filter("is_active", true).
		Filter("deleted_at__isnull", true)

	if excludeID > 0 {
		qs = qs.Exclude("id", excludeID)
	}

	// 检查系统分类或用户自定义分类 - 简化为只检查用户自定义分类
	qs = qs.Filter("type", "custom").Filter("user_id", userID)

	count, err := qs.Count()
	return count > 0, err
}

// GetCategoryWithChildren 获取分类及其子分类
func (r *CategoryRepository) GetCategoryWithChildren(categoryID int) (*model.Category, error) {
	o := orm.NewOrm()
	category := &model.Category{}
	err := o.QueryTable("categories").
		Filter("id", categoryID).
		Filter("is_active", true).
		Filter("deleted_at__isnull", true).
		One(category)
	if err != nil {
		return nil, err
	}

	// 获取子分类
	var children []*model.Category
	_, err = o.QueryTable("categories").
		Filter("parent_id", categoryID).
		Filter("is_active", true).
		Filter("deleted_at__isnull", true).
		OrderBy("sort_order", "created_at").
		All(&children)
	if err == nil {
		category.Children = children
	}

	return category, nil
}
