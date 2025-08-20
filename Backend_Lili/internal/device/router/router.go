package router

import (
	"Backend_Lili/internal/auth/middleware"
	"Backend_Lili/internal/device/controller"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// 创建控制器实例
	deviceController := controller.NewDeviceController()
	categoryController := controller.NewCategoryController()
	templateController := controller.NewTemplateController()

	// 设备模块路由
	ns := web.NewNamespace("/api/v1",
		// 设备管理相关路由 - 需要JWT认证
		web.NSNamespace("/devices",
			web.NSBefore(middleware.JWTAuth),

			// 设备CRUD操作
			web.NSRouter("/", deviceController, "get:GetDevicesList;post:CreateDevice"),
			web.NSRouter("/:deviceId", deviceController, "get:GetDeviceDetail;put:UpdateDevice;delete:DeleteDevice"),

			// 设备状态管理
			web.NSRouter("/:deviceId/status", deviceController, "patch:UpdateDeviceStatus"),

			// 设备价值评估
			web.NSRouter("/:deviceId/valuation", deviceController, "get:GetDeviceValuation"),

			// 设备价格预测
			web.NSRouter("/:deviceId/prediction", deviceController, "get:PredictDevicePrice"),

			// 批量导入设备
			web.NSRouter("/import", deviceController, "post:BatchImportDevices"),

			// 设备图片管理
			web.NSRouter("/:deviceId/images", deviceController, "get:GetDeviceImages;post:UploadDeviceImage"),
			web.NSRouter("/:deviceId/images/:imageId", deviceController, "delete:DeleteDeviceImage"),
		),

		// 设备模板相关路由 - 需要JWT认证
		web.NSNamespace("/device-templates",
			web.NSBefore(middleware.JWTAuth),

			// 获取热门模板 - 需要在具体ID路由之前
			web.NSRouter("/popular", templateController, "get:GetPopularTemplates"),
			// 获取推荐模板
			web.NSRouter("/recommendations", templateController, "get:GetRecommendedTemplates"),

			// 基本CRUD操作
			web.NSRouter("/", templateController, "get:GetTemplatesList;post:CreateTemplate"),
			web.NSRouter("/:template_id", templateController, "get:GetTemplateDetail;put:UpdateTemplate;delete:DeleteTemplate"),

			// 模板字段定义
			web.NSRouter("/:template_id/fields", templateController, "get:GetTemplateFields"),

			// 数据验证
			web.NSRouter("/:template_id/validate", templateController, "post:ValidateDeviceData"),

			// 模板统计
			web.NSRouter("/:template_id/statistics", templateController, "get:GetTemplateStatistics"),
		),

		// 设备分类相关路由 - 需要JWT认证
		web.NSNamespace("/categories",
			web.NSBefore(middleware.JWTAuth),

			// 特殊路由 - 需要在具体ID路由之前
			web.NSRouter("/system", categoryController, "get:GetSystemCategories"),
			web.NSRouter("/custom", categoryController, "get:GetCustomCategories"),
			web.NSRouter("/sort", categoryController, "put:SortCategories"),
			web.NSRouter("/statistics", categoryController, "get:GetCategoryStatistics"),
			web.NSRouter("/search", categoryController, "get:SearchCategories"),

			// 基本CRUD操作
			web.NSRouter("/", categoryController, "get:GetCategoriesList;post:CreateCustomCategory"),
			web.NSRouter("/:category_id", categoryController, "get:GetCategoryDetail;put:UpdateCustomCategory;delete:DeleteCustomCategory"),
		),
	)

	web.AddNamespace(ns)
}
