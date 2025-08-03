package controller

import (
	"Backend_Lili/internal/device/service"
	"Backend_Lili/pkg/utils"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type CategoryController struct {
	web.Controller
	categoryService *service.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: service.NewCategoryService(),
	}
}

// GetAllCategories 获取所有分类
// @router /categories [get]
func (c *CategoryController) GetAllCategories() {
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, categories)
}

// GetCategoryByID 根据ID获取分类
// @router /categories/:categoryId [get]
func (c *CategoryController) GetCategoryByID() {
	categoryIDStr := c.Ctx.Input.Param(":categoryId")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "分类ID格式错误")
		return
	}

	category, err := c.categoryService.GetCategoryByID(categoryID)
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

// GetCategoriesByParentID 根据父分类ID获取子分类
// @router /categories/:parentId/children [get]
func (c *CategoryController) GetCategoriesByParentID() {
	parentIDStr := c.Ctx.Input.Param(":parentId")
	parentID, err := strconv.Atoi(parentIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "父分类ID格式错误")
		return
	}

	categories, err := c.categoryService.GetCategoriesByParentID(parentID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, categories)
}

// GetCategoryTree 获取分类树结构
// @router /categories/tree [get]
func (c *CategoryController) GetCategoryTree() {
	categoryTree, err := c.categoryService.GetCategoryTree()
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, categoryTree)
}
