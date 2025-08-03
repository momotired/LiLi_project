package service

import (
	"Backend_Lili/internal/device/model"
	"Backend_Lili/internal/device/repository"
	"Backend_Lili/pkg/utils"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repository.NewCategoryRepository(),
	}
}

// GetAllCategories 获取所有分类
func (s *CategoryService) GetAllCategories() ([]*model.Category, error) {
	categories, err := s.categoryRepo.GetAllCategories()
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取分类列表失败")
	}
	return categories, nil
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(categoryID int) (*model.Category, error) {
	if categoryID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "分类ID无效")
	}

	category, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取分类信息失败")
	}
	if category == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "分类不存在")
	}

	return category, nil
}

// GetCategoriesByParentID 根据父分类ID获取子分类
func (s *CategoryService) GetCategoriesByParentID(parentID int) ([]*model.Category, error) {
	categories, err := s.categoryRepo.GetCategoriesByParentID(parentID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取子分类失败")
	}
	return categories, nil
}

// GetCategoryTree 获取分类树结构
func (s *CategoryService) GetCategoryTree() ([]*CategoryTreeNode, error) {
	// 获取所有分类
	allCategories, err := s.categoryRepo.GetAllCategories()
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取分类列表失败")
	}

	// 构建分类树
	categoryMap := make(map[int]*CategoryTreeNode)
	var rootCategories []*CategoryTreeNode

	// 第一遍遍历：创建所有节点
	for _, category := range allCategories {
		node := &CategoryTreeNode{
			ID:        category.ID,
			Name:      category.Name,
			ParentID:  category.ParentID,
			Icon:      category.Icon,
			SortOrder: category.SortOrder,
			IsActive:  category.IsActive,
			Children:  make([]*CategoryTreeNode, 0),
		}
		categoryMap[category.ID] = node
	}

	// 第二遍遍历：建立父子关系
	for _, category := range allCategories {
		node := categoryMap[category.ID]
		if category.ParentID == 0 {
			// 根分类
			rootCategories = append(rootCategories, node)
		} else {
			// 子分类
			if parent, exists := categoryMap[category.ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return rootCategories, nil
}

// 分类树节点结构
type CategoryTreeNode struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	ParentID  int                 `json:"parent_id"`
	Icon      string              `json:"icon"`
	SortOrder int                 `json:"sort_order"`
	IsActive  bool                `json:"is_active"`
	Children  []*CategoryTreeNode `json:"children"`
} 