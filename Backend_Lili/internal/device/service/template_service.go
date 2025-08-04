package service

import (
	"Backend_Lili/internal/device/model"
	"Backend_Lili/internal/device/repository"
	"Backend_Lili/pkg/utils"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TemplateService struct {
	templateRepo *repository.TemplateRepository
	categoryRepo *repository.CategoryRepository
	deviceRepo   *repository.DeviceRepository
}

func NewTemplateService() *TemplateService {
	return &TemplateService{
		templateRepo: repository.NewTemplateRepository(),
		categoryRepo: repository.NewCategoryRepository(),
		deviceRepo:   repository.NewDeviceRepository(),
	}
}

// GetTemplatesList 获取设备模板列表
func (s *TemplateService) GetTemplatesList(req *GetTemplatesListRequest) (*GetTemplatesListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	templates, total, err := s.templateRepo.GetTemplatesWithPagination(req.CategoryID, req.Active, req.Page, req.Limit)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取模板列表失败")
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &GetTemplatesListResponse{
		Templates:  templates,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetTemplateDetail 获取设备模板详情
func (s *TemplateService) GetTemplateDetail(templateID int) (*model.DeviceTemplate, error) {
	if templateID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "模板ID无效")
	}

	template, err := s.templateRepo.GetTemplateByID(templateID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取模板详情失败")
	}
	if template == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "模板不存在")
	}

	return template, nil
}

// GetTemplateFields 获取模板字段定义
func (s *TemplateService) GetTemplateFields(templateID int) ([]*model.TemplateField, error) {
	if templateID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "模板ID无效")
	}

	template, err := s.templateRepo.GetTemplateByID(templateID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取模板失败")
	}
	if template == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "模板不存在")
	}

	var fields []*model.TemplateField
	if template.Fields != "" {
		err = json.Unmarshal([]byte(template.Fields), &fields)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_SERVER, "解析模板字段失败")
		}
	}

	return fields, nil
}

// CreateTemplate 创建设备模板（管理员）
func (s *TemplateService) CreateTemplate(userID int, req *CreateTemplateRequest) (*model.DeviceTemplate, error) {
	// TODO: 验证管理员权限
	// if !s.isAdmin(userID) {
	//     return nil, utils.NewBusinessError(utils.ERROR_FORBIDDEN, "需要管理员权限")
	// }

	// 验证必填参数
	if req.Name == "" {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "模板名称不能为空")
	}
	if req.CategoryID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "分类ID无效")
	}
	if len(req.Fields) == 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "字段定义不能为空")
	}

	// 验证分类是否存在
	category, err := s.categoryRepo.GetCategoryByID(req.CategoryID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "验证分类失败")
	}
	if category == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "分类不存在")
	}

	// 验证字段定义
	err = s.validateTemplateFields(req.Fields)
	if err != nil {
		return nil, err
	}

	// 序列化字段定义
	fieldsJSON, err := json.Marshal(req.Fields)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "序列化字段定义失败")
	}

	// 创建模板
	template := &model.DeviceTemplate{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Icon:        req.Icon,
		Fields:      string(fieldsJSON),
		IsActive:    req.Active,
		UseCount:    0,
	}

	err = s.templateRepo.CreateTemplate(template)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "创建模板失败")
	}

	return template, nil
}

// UpdateTemplate 更新设备模板（管理员）
func (s *TemplateService) UpdateTemplate(userID, templateID int, req *UpdateTemplateRequest) (*model.DeviceTemplate, error) {
	// TODO: 验证管理员权限
	// if !s.isAdmin(userID) {
	//     return nil, utils.NewBusinessError(utils.ERROR_FORBIDDEN, "需要管理员权限")
	// }

	if templateID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "模板ID无效")
	}

	// 获取现有模板
	template, err := s.templateRepo.GetTemplateByID(templateID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取模板失败")
	}
	if template == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "模板不存在")
	}

	// 更新字段
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.CategoryID > 0 {
		// 验证分类是否存在
		category, err := s.categoryRepo.GetCategoryByID(req.CategoryID)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_SERVER, "验证分类失败")
		}
		if category == nil {
			return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "分类不存在")
		}
		template.CategoryID = req.CategoryID
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.Icon != "" {
		template.Icon = req.Icon
	}
	if len(req.Fields) > 0 {
		// 验证字段定义
		err = s.validateTemplateFields(req.Fields)
		if err != nil {
			return nil, err
		}

		// 序列化字段定义
		fieldsJSON, err := json.Marshal(req.Fields)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_SERVER, "序列化字段定义失败")
		}
		template.Fields = string(fieldsJSON)
	}
	template.IsActive = req.Active

	err = s.templateRepo.UpdateTemplate(template)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "更新模板失败")
	}

	return template, nil
}

// DeleteTemplate 删除设备模板（管理员）
func (s *TemplateService) DeleteTemplate(userID, templateID int) error {
	// TODO: 验证管理员权限
	// if !s.isAdmin(userID) {
	//     return utils.NewBusinessError(utils.ERROR_FORBIDDEN, "需要管理员权限")
	// }

	if templateID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "模板ID无效")
	}

	// 检查模板是否被使用
	// TODO: 实现检查逻辑

	err := s.templateRepo.DeleteTemplate(templateID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_SERVER, "删除模板失败")
	}

	return nil
}

// GetPopularTemplates 获取热门设备模板
func (s *TemplateService) GetPopularTemplates(limit int) ([]*model.DeviceTemplate, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	templates, err := s.templateRepo.GetPopularTemplates(limit)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取热门模板失败")
	}

	return templates, nil
}

// ValidateDeviceData 根据模板验证设备数据
func (s *TemplateService) ValidateDeviceData(templateID int, req *ValidateDeviceDataRequest) (*ValidateDeviceDataResponse, error) {
	if templateID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "模板ID无效")
	}

	// 获取模板字段定义
	fields, err := s.GetTemplateFields(templateID)
	if err != nil {
		return nil, err
	}

	response := &ValidateDeviceDataResponse{
		Valid:   true,
		Errors:  make(map[string]string),
		Missing: make([]string, 0),
		Invalid: make(map[string]interface{}),
	}

	// 验证每个字段
	for _, field := range fields {
		value, exists := req.DeviceData[field.FieldName]

		// 检查必填字段
		if field.Required && (!exists || value == nil || value == "") {
			response.Valid = false
			response.Missing = append(response.Missing, field.FieldName)
			response.Errors[field.FieldName] = field.FieldLabel + "是必填字段"
			continue
		}

		// 如果字段不存在或为空，跳过验证
		if !exists || value == nil {
			continue
		}

		// 根据字段类型验证
		err := s.validateFieldValue(field, value)
		if err != nil {
			response.Valid = false
			response.Invalid[field.FieldName] = value
			response.Errors[field.FieldName] = err.Error()
		}
	}

	return response, nil
}

// GetRecommendedTemplates 获取推荐模板
func (s *TemplateService) GetRecommendedTemplates(userID int) ([]*model.DeviceTemplate, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	templates, err := s.templateRepo.GetRecommendedTemplates(userID, 10)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取推荐模板失败")
	}

	return templates, nil
}

// GetTemplateStatistics 获取模板统计信息
func (s *TemplateService) GetTemplateStatistics(templateID int) (*TemplateStatisticsResponse, error) {
	if templateID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "模板ID无效")
	}

	// 获取模板基本信息
	template, err := s.templateRepo.GetTemplateByID(templateID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取模板失败")
	}
	if template == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "模板不存在")
	}

	// 获取统计信息
	stats, err := s.templateRepo.GetTemplateStatistics(templateID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "获取统计信息失败")
	}

	response := &TemplateStatisticsResponse{
		TemplateID:  templateID,
		UseCount:    template.UseCount,
		DeviceCount: int(stats["device_count"].(int64)),
		UserCount:   0, // TODO: 实现用户数量统计
	}

	return response, nil
}

// validateTemplateFields 验证模板字段定义
func (s *TemplateService) validateTemplateFields(fields []*model.TemplateField) error {
	fieldNames := make(map[string]bool)

	for _, field := range fields {
		// 检查字段名称
		if field.FieldName == "" {
			return utils.NewBusinessError(utils.ERROR_PARAM, "字段名称不能为空")
		}
		if field.FieldLabel == "" {
			return utils.NewBusinessError(utils.ERROR_PARAM, "字段显示名称不能为空")
		}

		// 检查字段名称重复
		if fieldNames[field.FieldName] {
			return utils.NewBusinessError(utils.ERROR_PARAM, "字段名称重复: "+field.FieldName)
		}
		fieldNames[field.FieldName] = true

		// 检查字段类型
		validTypes := []string{"text", "number", "select", "date", "textarea"}
		if !s.contains(validTypes, field.FieldType) {
			return utils.NewBusinessError(utils.ERROR_PARAM, "无效的字段类型: "+field.FieldType)
		}

		// 检查select类型的选项
		if field.FieldType == "select" && len(field.Options) == 0 {
			return utils.NewBusinessError(utils.ERROR_PARAM, "select类型字段必须提供选项列表")
		}
	}

	return nil
}

// validateFieldValue 验证字段值
func (s *TemplateService) validateFieldValue(field *model.TemplateField, value interface{}) error {
	strValue := fmt.Sprintf("%v", value)

	switch field.FieldType {
	case "number":
		_, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return fmt.Errorf("字段 %s 必须是数字", field.FieldLabel)
		}
	case "date":
		_, err := time.Parse("2006-01-02", strValue)
		if err != nil {
			return fmt.Errorf("字段 %s 必须是有效日期格式(YYYY-MM-DD)", field.FieldLabel)
		}
	case "select":
		if !s.contains(field.Options, strValue) {
			return fmt.Errorf("字段 %s 的值不在有效选项中", field.FieldLabel)
		}
	case "text", "textarea":
		// 可以添加长度验证等
		if len(strValue) > 1000 && field.FieldType == "text" {
			return fmt.Errorf("字段 %s 长度不能超过1000字符", field.FieldLabel)
		}
	}

	// 验证规则
	if field.ValidationRules != "" {
		err := s.validateByRules(field.ValidationRules, strValue, field.FieldLabel)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateByRules 根据验证规则验证值
func (s *TemplateService) validateByRules(rules, value, fieldLabel string) error {
	// 简单的验证规则解析
	ruleList := strings.Split(rules, "|")

	for _, rule := range ruleList {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		if strings.HasPrefix(rule, "min:") {
			minStr := strings.TrimPrefix(rule, "min:")
			min, err := strconv.Atoi(minStr)
			if err == nil && len(value) < min {
				return fmt.Errorf("字段 %s 长度不能少于 %d 字符", fieldLabel, min)
			}
		} else if strings.HasPrefix(rule, "max:") {
			maxStr := strings.TrimPrefix(rule, "max:")
			max, err := strconv.Atoi(maxStr)
			if err == nil && len(value) > max {
				return fmt.Errorf("字段 %s 长度不能超过 %d 字符", fieldLabel, max)
			}
		} else if strings.HasPrefix(rule, "regex:") {
			pattern := strings.TrimPrefix(rule, "regex:")
			matched, err := regexp.MatchString(pattern, value)
			if err != nil {
				return fmt.Errorf("字段 %s 验证规则错误", fieldLabel)
			}
			if !matched {
				return fmt.Errorf("字段 %s 格式不正确", fieldLabel)
			}
		}
	}

	return nil
}

// contains 检查slice中是否包含指定值
func (s *TemplateService) contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
