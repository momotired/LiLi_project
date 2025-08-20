package service

import (
	"Backend_Lili/internal/device/model"
	"Backend_Lili/internal/device/repository"
	"Backend_Lili/pkg/utils"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	deviceRepo   *repository.DeviceRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repository.NewCategoryRepository(),
		deviceRepo:   repository.NewDeviceRepository(),
	}
}

// GetCategoriesList 获取分类列表
func (s *CategoryService) GetCategoriesList(userID int, req *GetCategoriesListRequest) (*GetCategoriesListResponse, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	// 设置默认值
	if req.Type == "" {
		req.Type = "all"
	}

	var categories []*model.Category
	var err error

	if req.ParentID > 0 {
		// 获取子分类
		categories, err = s.categoryRepo.GetCategoriesByParentID(req.ParentID)
	} else {
		// 根据类型获取分类
		categories, err = s.categoryRepo.GetCategoriesByType(req.Type, userID, req.IncludeCount)
	}

	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取分类列表失败")
	}

	// 构建层级结构
	categories = s.buildCategoryTree(categories)

	return &GetCategoriesListResponse{
		Categories: categories,
	}, nil
}

// GetCategoryDetail 获取分类详情
func (s *CategoryService) GetCategoryDetail(userID, categoryID int) (*model.Category, error) {
	if userID <= 0 || categoryID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	category, err := s.categoryRepo.GetCategoryWithChildren(categoryID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取分类详情失败")
	}
	if category == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "分类不存在")
	}

	// 验证访问权限（自定义分类只能访问自己的）
	if category.Type == "custom" && category.UserID != userID {
		return nil, utils.NewBusinessError(utils.ERROR_FORBIDDEN, "无权访问此分类")
	}

	// 获取设备数量
	count, err := s.categoryRepo.GetDeviceCountByCategory(categoryID, userID)
	if err == nil {
		category.DeviceCount = count
	}

	return category, nil
}

// CreateCustomCategory 创建自定义分类
func (s *CategoryService) CreateCustomCategory(userID int, req *CreateCategoryRequest) (*model.Category, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	// 验证必填参数
	if req.Name == "" {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "分类名称不能为空")
	}

	// 验证分类名称唯一性
	exists, err := s.categoryRepo.CheckCategoryExists(req.Name, userID, 0)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "验证分类名称失败")
	}
	if exists {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "分类名称已存在")
	}

	// 验证父分类
	if req.ParentID > 0 {
		parentCategory, err := s.categoryRepo.GetCategoryByID(req.ParentID)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_SERVER, "验证父分类失败")
		}
		if parentCategory == nil {
			return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "父分类不存在")
		}
		// 自定义分类的父分类必须是自己的或系统的
		if parentCategory.Type == "custom" && parentCategory.UserID != userID {
			return nil, utils.NewBusinessError(utils.ERROR_FORBIDDEN, "无权使用此父分类")
		}
	}

	// 创建分类
	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Icon:        req.Icon,
		Color:       req.Color,
		SortOrder:   req.SortOrder,
		Type:        "custom",
		UserID:      userID,
		IsActive:    true,
	}

	err = s.categoryRepo.CreateCategory(category)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "创建分类失败")
	}

	return category, nil
}

// UpdateCustomCategory 更新自定义分类
func (s *CategoryService) UpdateCustomCategory(userID, categoryID int, req *UpdateCategoryRequest) (*model.Category, error) {
	if userID <= 0 || categoryID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 获取现有分类
	category, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取分类失败")
	}
	if category == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "分类不存在")
	}

	// 验证权限（只能修改自己的自定义分类）
	if category.Type != "custom" || category.UserID != userID {
		return nil, utils.NewBusinessError(utils.ERROR_FORBIDDEN, "无权修改此分类")
	}

	// 更新字段
	if req.Name != "" {
		// 验证名称唯一性
		exists, err := s.categoryRepo.CheckCategoryExists(req.Name, userID, categoryID)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_SERVER, "验证分类名称失败")
		}
		if exists {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "分类名称已存在")
		}
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID > 0 {
		// 验证父分类
		parentCategory, err := s.categoryRepo.GetCategoryByID(req.ParentID)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_SERVER, "验证父分类失败")
		}
		if parentCategory == nil {
			return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "父分类不存在")
		}
		// 不能设置自己为父分类
		if req.ParentID == categoryID {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "不能设置自己为父分类")
		}
		category.ParentID = req.ParentID
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Color != "" {
		category.Color = req.Color
	}
	if req.SortOrder > 0 {
		category.SortOrder = req.SortOrder
	}

	err = s.categoryRepo.UpdateCategory(category)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "更新分类失败")
	}

	return category, nil
}

// DeleteCustomCategory 删除自定义分类
func (s *CategoryService) DeleteCustomCategory(userID, categoryID int) error {
	if userID <= 0 || categoryID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 获取分类信息
	category, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_SERVER, "获取分类失败")
	}
	if category == nil {
		return utils.NewBusinessError(utils.ERROR_NOT_FOUND, "分类不存在")
	}

	// 验证权限
	if category.Type != "custom" || category.UserID != userID {
		return utils.NewBusinessError(utils.ERROR_FORBIDDEN, "无权删除此分类")
	}

	// 检查分类下是否有设备
	count, err := s.categoryRepo.GetDeviceCountByCategory(categoryID, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_SERVER, "检查分类使用情况失败")
	}
	if count > 0 {
		return utils.NewBusinessError(utils.ERROR_BUSINESS, "分类下有设备，无法删除")
	}

	// 检查是否有子分类
	children, err := s.categoryRepo.GetCategoriesByParentID(categoryID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_SERVER, "检查子分类失败")
	}
	if len(children) > 0 {
		return utils.NewBusinessError(utils.ERROR_BUSINESS, "分类下有子分类，无法删除")
	}

	err = s.categoryRepo.DeleteCategory(categoryID, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_SERVER, "删除分类失败")
	}

	return nil
}

// GetSystemCategories 获取系统默认分类
func (s *CategoryService) GetSystemCategories() ([]*model.Category, error) {
	categories, err := s.categoryRepo.GetSystemCategories()
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取系统分类失败")
	}

	// 构建层级结构
	categories = s.buildCategoryTree(categories)

	return categories, nil
}

// GetCustomCategories 获取用户自定义分类
func (s *CategoryService) GetCustomCategories(userID int) ([]*model.Category, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	categories, err := s.categoryRepo.GetCustomCategories(userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取自定义分类失败")
	}

	// 构建层级结构
	categories = s.buildCategoryTree(categories)

	return categories, nil
}

// SortCategories 分类排序
func (s *CategoryService) SortCategories(userID int, req *SortCategoriesRequest) error {
	if userID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}
	if len(req.CategoryOrders) == 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "排序数据不能为空")
	}

	// 转换数据格式
	categoryOrders := make([]map[string]int, len(req.CategoryOrders))
	for i, order := range req.CategoryOrders {
		categoryOrders[i] = map[string]int{
			"category_id": order.CategoryID,
			"sort_order":  order.SortOrder,
		}
	}

	err := s.categoryRepo.UpdateCategoriesSort(userID, categoryOrders)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_SERVER, "更新分类排序失败")
	}

	return nil
}

// GetCategoryStatistics 获取分类统计
func (s *CategoryService) GetCategoryStatistics(userID int, req *GetCategoryStatisticsRequest) (*CategoryStatisticsResponse, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	// 设置默认周期
	if req.Period == "" {
		req.Period = "month"
	}

	stats, err := s.categoryRepo.GetCategoryStatistics(userID, req.Period)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取分类统计失败")
	}

	// 构建响应
	response := &CategoryStatisticsResponse{
		Period:       req.Period,
		TotalDevices: int(stats["total_devices"].(int64)),
		TotalValue:   stats["total_value"].(float64),
		Categories:   make([]*CategoryStatInfo, 0),
	}

	// 处理分类统计数据
	if categoryStats, ok := stats["categories"].([]map[string]interface{}); ok {
		for _, stat := range categoryStats {
			categoryInfo := &CategoryStatInfo{
				CategoryID:   int(stat["category_id"].(int64)),
				CategoryName: stat["category_name"].(string),
				DeviceCount:  int(stat["device_count"].(int64)),
				TotalValue:   stat["total_value"].(float64),
				AvgValue:     stat["avg_value"].(float64),
			}
			// 计算占比
			if response.TotalDevices > 0 {
				categoryInfo.Percentage = float64(categoryInfo.DeviceCount) / float64(response.TotalDevices) * 100
			}
			response.Categories = append(response.Categories, categoryInfo)
		}
	}

	return response, nil
}

// SearchCategories 搜索分类
func (s *CategoryService) SearchCategories(userID int, req *SearchCategoriesRequest) ([]*model.Category, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}
	if req.Keyword == "" {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "搜索关键词不能为空")
	}

	// 设置默认类型
	if req.Type == "" {
		req.Type = "all"
	}

	categories, err := s.categoryRepo.SearchCategories(req.Keyword, req.Type, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "搜索分类失败")
	}

	return categories, nil
}

// buildCategoryTree 构建分类层级结构
func (s *CategoryService) buildCategoryTree(categories []*model.Category) []*model.Category {
	// 创建ID到分类的映射
	categoryMap := make(map[int]*model.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
		category.Children = make([]*model.Category, 0)
	}

	// 构建树形结构
	var rootCategories []*model.Category
	for _, category := range categories {
		if category.ParentID == 0 {
			rootCategories = append(rootCategories, category)
		} else if parent, exists := categoryMap[category.ParentID]; exists {
			parent.Children = append(parent.Children, category)
		}
	}

	return rootCategories
}
