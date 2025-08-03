package controller

import (
	"Backend_Lili/internal/device/service"
	"Backend_Lili/pkg/utils"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type TemplateController struct {
	web.Controller
	templateService *service.TemplateService
}

func NewTemplateController() *TemplateController {
	return &TemplateController{
		templateService: service.NewTemplateService(),
	}
}

// GetAllTemplates 获取所有设备模板
// @router /templates [get]
func (c *TemplateController) GetAllTemplates() {
	templates, err := c.templateService.GetAllTemplates()
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, templates)
}

// GetTemplatesByCategory 根据分类获取模板
// @router /templates/category/:categoryId [get]
func (c *TemplateController) GetTemplatesByCategory() {
	categoryIDStr := c.Ctx.Input.Param(":categoryId")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "分类ID格式错误")
		return
	}

	templates, err := c.templateService.GetTemplatesByCategory(categoryID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, templates)
}

// GetTemplateByID 根据ID获取模板
// @router /templates/:templateId [get]
func (c *TemplateController) GetTemplateByID() {
	templateIDStr := c.Ctx.Input.Param(":templateId")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "模板ID格式错误")
		return
	}

	template, err := c.templateService.GetTemplateByID(templateID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, template)
}
