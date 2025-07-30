package router

import (
	authCtrl "Backend_Lili/internal/auth/controller"
	"Backend_Lili/internal/auth/middleware"
	userCtrl "Backend_Lili/internal/user/controller"

	beego "github.com/beego/beego/v2/server/web"
)

// Init 统一初始化所有路由
func Init() {
	// 创建控制器实例
	authController := &authCtrl.AuthController{}
	userController := userCtrl.NewUserController()

	// 创建API命名空间
	ns := beego.NewNamespace("/api/v1",
		// 认证相关路由
		beego.NSNamespace("/auth",
			// 不需要认证的接口
			beego.NSRouter("/login", authController, "post:Login"),
			beego.NSRouter("/refresh", authController, "post:RefreshToken"),

			// 需要认证的接口
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
	)

	// 注册命名空间
	beego.AddNamespace(ns)

	// 健康检查路由
	beego.Router("/health", authController, "get:Health")
}
