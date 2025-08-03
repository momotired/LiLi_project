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