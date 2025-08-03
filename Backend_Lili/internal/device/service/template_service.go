package service

import (
	"Backend_Lili/internal/device/model"
	"Backend_Lili/internal/device/repository"
	"Backend_Lili/pkg/utils"
)

type TemplateService struct {
	templateRepo *repository.TemplateRepository
	categoryRepo *repository.CategoryRepository
}

func NewTemplateService() *TemplateService {
	return &TemplateService{
		templateRepo: repository.NewTemplateRepository(),
		categoryRepo: repository.NewCategoryRepository(),
	}
}

// GetAllTemplates 获取所有设备模板
func (s *TemplateService) GetAllTemplates() ([]*model.DeviceTemplate, error) {
	templates, err := s.templateRepo.GetAllTemplates()
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备模板列表失败")
	}
	return templates, nil
}

// GetTemplatesByCategory 根据分类获取模板
func (s *TemplateService) GetTemplatesByCategory(categoryID int) ([]*model.DeviceTemplate, error) {
	if categoryID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "分类ID无效")
	}

	// 验证分类是否存在
	category, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "查询分类信息失败")
	}
	if category == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "分类不存在")
	}

	templates, err := s.templateRepo.GetTemplatesByCategory(categoryID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备模板失败")
	}
	return templates, nil
}

// GetTemplateByID 根据ID获取模板
func (s *TemplateService) GetTemplateByID(templateID int) (*model.DeviceTemplate, error) {
	if templateID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "模板ID无效")
	}

	template, err := s.templateRepo.GetTemplateByID(templateID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备模板失败")
	}
	if template == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备模板不存在")
	}

	return template, nil
} 