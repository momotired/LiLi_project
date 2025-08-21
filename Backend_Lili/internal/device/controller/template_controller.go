package controller

import (
	base "Backend_Lili/internal/auth/controller"
	"Backend_Lili/internal/device/service"
	"Backend_Lili/pkg/utils"
	"encoding/json"
	"strconv"
)

type TemplateController struct {
	base.BaseController
	templateService *service.TemplateService
}

func NewTemplateController() *TemplateController {
	return &TemplateController{}
}

func (c *TemplateController) Prepare() {
	// 调用父类的Prepare方法
	c.BaseController.Prepare()
	
	// 初始化模板服务
	c.templateService = service.NewTemplateService()
}

// GetTemplatesList 获取设备模板列表
// @router /device-templates [get]
func (c *TemplateController) GetTemplatesList() {
	// 解析查询参数
	req := &service.GetTemplatesListRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	response, err := c.templateService.GetTemplatesList(req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// GetTemplateDetail 获取设备模板详情
// @router /device-templates/:template_id [get]
func (c *TemplateController) GetTemplateDetail() {
	// 获取路径参数
	templateIDStr := c.Ctx.Input.Param(":template_id")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "模板ID格式错误")
		return
	}

	// 调用服务层
	template, err := c.templateService.GetTemplateDetail(templateID)
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

// GetTemplateFields 获取模板字段定义
// @router /device-templates/:template_id/fields [get]
func (c *TemplateController) GetTemplateFields() {
	// 获取路径参数
	templateIDStr := c.Ctx.Input.Param(":template_id")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "模板ID格式错误")
		return
	}

	// 调用服务层
	fields, err := c.templateService.GetTemplateFields(templateID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"fields": fields,
	})
}

// CreateTemplate 创建设备模板（管理员）
// @router /device-templates [post]
func (c *TemplateController) CreateTemplate() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析请求体
	req := &service.CreateTemplateRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "请求体解析失败")
		return
	}

	// 调用服务层
	template, err := c.templateService.CreateTemplate(userID, req)
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

// UpdateTemplate 更新设备模板（管理员）
// @router /device-templates/:template_id [put]
func (c *TemplateController) UpdateTemplate() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取路径参数
	templateIDStr := c.Ctx.Input.Param(":template_id")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "模板ID格式错误")
		return
	}

	// 解析请求体
	req := &service.UpdateTemplateRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "请求体解析失败")
		return
	}

	// 调用服务层
	template, err := c.templateService.UpdateTemplate(userID, templateID, req)
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

// DeleteTemplate 删除设备模板（管理员）
// @router /device-templates/:template_id [delete]
func (c *TemplateController) DeleteTemplate() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取路径参数
	templateIDStr := c.Ctx.Input.Param(":template_id")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "模板ID格式错误")
		return
	}

	// 调用服务层
	err = c.templateService.DeleteTemplate(userID, templateID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"message": "模板删除成功",
	})
}

// GetPopularTemplates 获取热门设备模板
// @router /device-templates/popular [get]
func (c *TemplateController) GetPopularTemplates() {
	// 解析查询参数
	limitStr := c.GetString("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// 调用服务层
	templates, err := c.templateService.GetPopularTemplates(limit)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"templates": templates,
	})
}

// ValidateDeviceData 根据模板验证设备数据
// @router /device-templates/:template_id/validate [post]
func (c *TemplateController) ValidateDeviceData() {
	// 获取路径参数
	templateIDStr := c.Ctx.Input.Param(":template_id")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "模板ID格式错误")
		return
	}

	// 解析请求体
	req := &service.ValidateDeviceDataRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "请求体解析失败")
		return
	}

	// 调用服务层
	response, err := c.templateService.ValidateDeviceData(templateID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// GetRecommendedTemplates 获取推荐模板
// @router /device-templates/recommendations [get]
func (c *TemplateController) GetRecommendedTemplates() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 调用服务层
	templates, err := c.templateService.GetRecommendedTemplates(userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"templates": templates,
	})
}

// GetTemplateStatistics 获取模板统计信息
// @router /device-templates/:template_id/statistics [get]
func (c *TemplateController) GetTemplateStatistics() {
	// 获取路径参数
	templateIDStr := c.Ctx.Input.Param(":template_id")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "模板ID格式错误")
		return
	}

	// 调用服务层
	statistics, err := c.templateService.GetTemplateStatistics(templateID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, statistics)
}
