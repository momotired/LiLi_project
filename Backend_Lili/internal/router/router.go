package router

import (
	authCtrl "Backend_Lili/internal/auth/controller"
	"Backend_Lili/internal/auth/middleware"
	deviceCtrl "Backend_Lili/internal/device/controller"
	priceRouter "Backend_Lili/internal/price/router"
    statsRouter "Backend_Lili/internal/statistics/router"
    tagsRouter "Backend_Lili/internal/tags/router"
	userCtrl "Backend_Lili/internal/user/controller"

	beego "github.com/beego/beego/v2/server/web"
)

// Init 统一初始化所有路由
func Init() {
	// 创建控制器实例
	authController := &authCtrl.AuthController{}
	userController := userCtrl.NewUserController()

	// 添加全局CORS中间件
	beego.InsertFilter("*", beego.BeforeRouter, middleware.GlobalCORS) //用于注册全局或特定路由添加过滤器（Filter） 的核心函数
	
	// 添加全局访问日志中间件
	beego.InsertFilter("*", beego.BeforeRouter, middleware.AccessLog)

	// 创建设备控制器实例
	deviceController := deviceCtrl.NewDeviceController()
	categoryController := deviceCtrl.NewCategoryController()
	templateController := deviceCtrl.NewTemplateController()

	// 创建API命名空间
	ns := beego.NewNamespace("/api/v1",
		// 认证相关路由
		beego.NSNamespace("/auth",
			beego.NSBefore(middleware.ConditionalAuth), // 条件认证中间件
			beego.NSRouter("/login", authController, "post:Login"),
			beego.NSRouter("/refresh", authController, "post:RefreshToken"),
			beego.NSRouter("/logout", authController, "post:Logout"),
			beego.NSRouter("/verify", authController, "get:VerifyToken"),
		),

		// 用户相关路由
		beego.NSNamespace("/users",
			// 所有用户接口都需要JWT认证
			beego.NSBefore(middleware.JWTAuth),

			// 1. 用户基本信息管理
			beego.NSRouter("/profile", userController, "get:GetProfile;put:UpdateProfile"),

			// 2. 用户偏好设置管理
			beego.NSRouter("/preferences", userController, "get:GetPreferences;put:UpdatePreferences"),

			// 3. 用户标签管理
			beego.NSRouter("/tags", userController, "get:GetTags;put:UpdateTags"),

			// 4. 用户统计信息
			beego.NSRouter("/statistics", userController, "get:GetStatistics"),

			// 5. 用户账号注销
			beego.NSRouter("/account", userController, "delete:DeleteAccount"),
		),

		// 设备管理相关路由 - 需要JWT认证
		beego.NSNamespace("/devices",
			beego.NSBefore(middleware.JWTAuth),

			// 设备CRUD操作
			beego.NSRouter("/", deviceController, "get:GetDevicesList;post:CreateDevice"),
			beego.NSRouter("/:deviceId", deviceController, "get:GetDeviceDetail;put:UpdateDevice;delete:DeleteDevice"),

			// 设备状态管理
			beego.NSRouter("/:deviceId/status", deviceController, "patch:UpdateDeviceStatus"),

			// 设备价值评估
			beego.NSRouter("/:deviceId/valuation", deviceController, "get:GetDeviceValuation"),

			// 设备价格预测
			beego.NSRouter("/:deviceId/prediction", deviceController, "get:PredictDevicePrice"),

			// 批量导入设备
			beego.NSRouter("/import", deviceController, "post:BatchImportDevices"),

			// 设备图片管理
			beego.NSRouter("/:deviceId/images", deviceController, "get:GetDeviceImages;post:UploadDeviceImage"),
			beego.NSRouter("/:deviceId/images/:imageId", deviceController, "delete:DeleteDeviceImage"),
		),

		// 设备模板相关路由 - 需要JWT认证
		beego.NSNamespace("/device-templates",
			beego.NSBefore(middleware.JWTAuth),

			// 获取热门模板 - 需要在具体ID路由之前
			beego.NSRouter("/popular", templateController, "get:GetPopularTemplates"),
			// 获取推荐模板
			beego.NSRouter("/recommendations", templateController, "get:GetRecommendedTemplates"),

			// 基本CRUD操作
			beego.NSRouter("/", templateController, "get:GetTemplatesList;post:CreateTemplate"),
			beego.NSRouter("/:template_id", templateController, "get:GetTemplateDetail;put:UpdateTemplate;delete:DeleteTemplate"),

			// 模板字段定义
			beego.NSRouter("/:template_id/fields", templateController, "get:GetTemplateFields"),

			// 数据验证
			beego.NSRouter("/:template_id/validate", templateController, "post:ValidateDeviceData"),

			// 模板统计
			beego.NSRouter("/:template_id/statistics", templateController, "get:GetTemplateStatistics"),
		),

		// 设备分类相关路由 - 需要JWT认证
		beego.NSNamespace("/categories",
			beego.NSBefore(middleware.JWTAuth),

			// 特殊路由 - 需要在具体ID路由之前
			beego.NSRouter("/system", categoryController, "get:GetSystemCategories"),
			beego.NSRouter("/custom", categoryController, "get:GetCustomCategories"),
			beego.NSRouter("/sort", categoryController, "put:SortCategories"),
			beego.NSRouter("/statistics", categoryController, "get:GetCategoryStatistics"),
			beego.NSRouter("/search", categoryController, "get:SearchCategories"),

			// 基本CRUD操作
			beego.NSRouter("/", categoryController, "get:GetCategoriesList;post:CreateCustomCategory"),
			beego.NSRouter("/:category_id", categoryController, "get:GetCategoryDetail;put:UpdateCustomCategory;delete:DeleteCustomCategory"),
		),
	)

	// 注册命名空间
	beego.AddNamespace(ns)

	// 初始化价格模块路由
	priceRouter.InitPriceRoutes()
    // 初始化统计模块路由
    statsRouter.InitStatisticsRoutes()
    // 初始化标签模块路由
    tagsRouter.InitTagsRoutes()

	// 健康检查路由
	beego.Router("/health", authController, "get:Health")
}
