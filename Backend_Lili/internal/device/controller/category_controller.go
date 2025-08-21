package controller

import (
	base "Backend_Lili/internal/auth/controller"
	"Backend_Lili/internal/device/service"
	"Backend_Lili/pkg/utils"
	"encoding/json"
	"strconv"
)

type CategoryController struct {
	base.BaseController
	categoryService *service.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

func (c *CategoryController) Prepare() {
	// 调用父类的Prepare方法
	c.BaseController.Prepare()
	
	// 初始化分类服务
	c.categoryService = service.NewCategoryService()
}

// GetCategoriesList 获取分类列表
// @router /categories [get]
func (c *CategoryController) GetCategoriesList() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析查询参数
	req := &service.GetCategoriesListRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	response, err := c.categoryService.GetCategoriesList(userID, req)
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

// GetCategoryDetail 获取分类详情
// @router /categories/:category_id [get]
func (c *CategoryController) GetCategoryDetail() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取路径参数
	categoryIDStr := c.Ctx.Input.Param(":category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "分类ID格式错误")
		return
	}

	// 调用服务层
	category, err := c.categoryService.GetCategoryDetail(userID, categoryID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, category)
}

// CreateCustomCategory 创建自定义分类
// @router /categories [post]
func (c *CategoryController) CreateCustomCategory() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析请求体
	req := &service.CreateCategoryRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "请求体解析失败")
		return
	}

	// 调用服务层
	category, err := c.categoryService.CreateCustomCategory(userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, category)
}

// UpdateCustomCategory 更新自定义分类
// @router /categories/:category_id [put]
func (c *CategoryController) UpdateCustomCategory() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取路径参数
	categoryIDStr := c.Ctx.Input.Param(":category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "分类ID格式错误")
		return
	}

	// 解析请求体
	req := &service.UpdateCategoryRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "请求体解析失败")
		return
	}

	// 调用服务层
	category, err := c.categoryService.UpdateCustomCategory(userID, categoryID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, category)
}

// DeleteCustomCategory 删除自定义分类
// @router /categories/:category_id [delete]
func (c *CategoryController) DeleteCustomCategory() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取路径参数
	categoryIDStr := c.Ctx.Input.Param(":category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "分类ID格式错误")
		return
	}

	// 调用服务层
	err = c.categoryService.DeleteCustomCategory(userID, categoryID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"message": "分类删除成功",
	})
}

// GetSystemCategories 获取系统默认分类
// @router /categories/system [get]
func (c *CategoryController) GetSystemCategories() {
	// 调用服务层
	categories, err := c.categoryService.GetSystemCategories()
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"categories": categories,
	})
}

// GetCustomCategories 获取用户自定义分类
// @router /categories/custom [get]
func (c *CategoryController) GetCustomCategories() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 调用服务层
	categories, err := c.categoryService.GetCustomCategories(userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"categories": categories,
	})
}

// SortCategories 分类排序
// @router /categories/sort [put]
func (c *CategoryController) SortCategories() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析请求体
	req := &service.SortCategoriesRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "请求体解析失败")
		return
	}

	// 调用服务层
	err := c.categoryService.SortCategories(userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"message": "分类排序更新成功",
	})
}

// GetCategoryStatistics 获取分类统计
// @router /categories/statistics [get]
func (c *CategoryController) GetCategoryStatistics() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析查询参数
	req := &service.GetCategoryStatisticsRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	statistics, err := c.categoryService.GetCategoryStatistics(userID, req)
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

// SearchCategories 搜索分类
// @router /categories/search [get]
func (c *CategoryController) SearchCategories() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析查询参数
	req := &service.SearchCategoriesRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	categories, err := c.categoryService.SearchCategories(userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{
		"categories": categories,
	})
}
