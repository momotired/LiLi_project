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

		// 设备分类相关路由 - 不需要认证，公共数据
		web.NSNamespace("/categories",
			web.NSRouter("/", categoryController, "get:GetAllCategories"),
			web.NSRouter("/tree", categoryController, "get:GetCategoryTree"),
			web.NSRouter("/:categoryId", categoryController, "get:GetCategoryByID"),
			web.NSRouter("/:parentId/children", categoryController, "get:GetCategoriesByParentID"),
		),

		// 设备模板相关路由 - 不需要认证，公共数据
		web.NSNamespace("/templates",
			web.NSRouter("/", templateController, "get:GetAllTemplates"),
			web.NSRouter("/:templateId", templateController, "get:GetTemplateByID"),
			web.NSRouter("/category/:categoryId", templateController, "get:GetTemplatesByCategory"),
		),
	)

	web.AddNamespace(ns)
}
