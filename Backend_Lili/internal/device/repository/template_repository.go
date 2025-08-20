package repository

import (
	"Backend_Lili/internal/device/model"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type TemplateRepository struct{}

func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{}
}

// GetAllTemplates 获取所有设备模板
func (r *TemplateRepository) GetAllTemplates() ([]*model.DeviceTemplate, error) {
	o := orm.NewOrm()
	var templates []*model.DeviceTemplate
	_, err := o.QueryTable("device_templates").
		Filter("is_active", true).
		RelatedSel().
		OrderBy("category_id", "created_at").
		All(&templates)
	return templates, err
}

// GetTemplatesByCategory 根据分类获取模板
func (r *TemplateRepository) GetTemplatesByCategory(categoryID int) ([]*model.DeviceTemplate, error) {
	o := orm.NewOrm()
	var templates []*model.DeviceTemplate
	_, err := o.QueryTable("device_templates").
		Filter("category_id", categoryID).
		Filter("is_active", true).
		RelatedSel().
		OrderBy("created_at").
		All(&templates)
	return templates, err
}

// GetTemplateByID 根据ID获取模板
func (r *TemplateRepository) GetTemplateByID(templateID int) (*model.DeviceTemplate, error) {
	o := orm.NewOrm()
	template := &model.DeviceTemplate{}
	err := o.QueryTable("device_templates").
		Filter("id", templateID).
		Filter("is_active", true).
		RelatedSel().
		One(template)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return template, err
}

// CreateTemplate 创建模板
func (r *TemplateRepository) CreateTemplate(template *model.DeviceTemplate) error {
	o := orm.NewOrm()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	_, err := o.Insert(template)
	return err
}

// UpdateTemplate 更新模板
func (r *TemplateRepository) UpdateTemplate(template *model.DeviceTemplate) error {
	o := orm.NewOrm()
	template.UpdatedAt = time.Now()
	_, err := o.Update(template)
	return err
}

// GetTemplatesWithPagination 分页获取模板列表
func (r *TemplateRepository) GetTemplatesWithPagination(categoryID int, active bool, page, limit int) ([]*model.DeviceTemplate, int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("device_templates")

	// 添加条件
	if categoryID > 0 {
		qs = qs.Filter("category_id", categoryID)
	}
	if active {
		qs = qs.Filter("is_active", true)
	}
	qs = qs.Filter("deleted_at__isnull", true)

	// 获取总数
	total, err := qs.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	var templates []*model.DeviceTemplate
	offset := (page - 1) * limit
	_, err = qs.RelatedSel().
		OrderBy("-use_count", "created_at").
		Limit(limit, offset).
		All(&templates)

	return templates, int(total), err
}

// GetPopularTemplates 获取热门模板
func (r *TemplateRepository) GetPopularTemplates(limit int) ([]*model.DeviceTemplate, error) {
	o := orm.NewOrm()
	var templates []*model.DeviceTemplate
	_, err := o.QueryTable("device_templates").
		Filter("is_active", true).
		Filter("deleted_at__isnull", true).
		RelatedSel().
		OrderBy("-use_count", "-created_at").
		Limit(limit).
		All(&templates)
	return templates, err
}

// DeleteTemplate 软删除模板
func (r *TemplateRepository) DeleteTemplate(templateID int) error {
	o := orm.NewOrm()
	template := &model.DeviceTemplate{}
	err := o.QueryTable("device_templates").Filter("id", templateID).One(template)
	if err != nil {
		return err
	}

	template.DeletedAt = time.Now()
	template.IsActive = false
	_, err = o.Update(template, "deleted_at", "is_active")
	return err
}

// IncrementUseCount 增加模板使用次数
func (r *TemplateRepository) IncrementUseCount(templateID int) error {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE device_templates SET use_count = use_count + 1 WHERE id = ?", templateID).Exec()
	return err
}

// GetTemplateStatistics 获取模板统计信息
func (r *TemplateRepository) GetTemplateStatistics(templateID int) (map[string]interface{}, error) {
	o := orm.NewOrm()
	stats := make(map[string]interface{})

	// 获取基于该模板创建的设备数量
	deviceCount, err := o.QueryTable("devices").Filter("template_id", templateID).Filter("deleted_at__isnull", true).Count()
	if err != nil {
		return nil, err
	}
	stats["device_count"] = deviceCount

	// 获取使用该模板的用户数量
	var userCount int64
	err = o.Raw("SELECT COUNT(DISTINCT user_id) FROM devices WHERE template_id = ? AND deleted_at IS NULL", templateID).QueryRow(&userCount)
	if err != nil {
		return nil, err
	}
	stats["user_count"] = userCount

	return stats, nil
}

// GetRecommendedTemplates 获取推荐模板
func (r *TemplateRepository) GetRecommendedTemplates(userID int, limit int) ([]*model.DeviceTemplate, error) {
	o := orm.NewOrm()

	// 根据用户历史创建的设备分类推荐模板
	var templates []*model.DeviceTemplate
	_, err := o.Raw(`
		SELECT DISTINCT dt.* FROM device_templates dt
		INNER JOIN categories c ON dt.category_id = c.id
		WHERE dt.is_active = 1 AND dt.deleted_at IS NULL
		AND c.id IN (
			SELECT DISTINCT d.category_id FROM devices d 
			WHERE d.user_id = ? AND d.deleted_at IS NULL
		)
		ORDER BY dt.use_count DESC, dt.created_at DESC
		LIMIT ?
	`, userID, limit).QueryRows(&templates)

	return templates, err
}

// SearchTemplates 搜索模板
func (r *TemplateRepository) SearchTemplates(keyword string, categoryID int) ([]*model.DeviceTemplate, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("device_templates").
		Filter("is_active", true).
		Filter("deleted_at__isnull", true)

	if keyword != "" {
		qs = qs.Filter("name__icontains", keyword)
	}
	if categoryID > 0 {
		qs = qs.Filter("category_id", categoryID)
	}

	var templates []*model.DeviceTemplate
	_, err := qs.RelatedSel().OrderBy("-use_count", "name").All(&templates)
	return templates, err
}
