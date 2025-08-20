package router

import (
	"Backend_Lili/internal/auth/middleware"
	"Backend_Lili/internal/user/controller"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// 创建用户控制器实例
	userController := controller.NewUserController()

	ns := web.NewNamespace("/api/v1",
		web.NSNamespace("/users",
			// 所有用户接口都需要JWT认证
			web.NSBefore(middleware.JWTAuth),

			// 1. 用户基本信息管理
			web.NSRouter("/profile", userController, "get:GetProfile;put:UpdateProfile"),

			// 2. 用户偏好设置管理
			web.NSRouter("/preferences", userController, "get:GetPreferences;put:UpdatePreferences"),

			// 3. 用户标签管理
			web.NSRouter("/tags", userController, "get:GetTags;put:UpdateTags"),

			// 4. 用户统计信息
			web.NSRouter("/statistics", userController, "get:GetStatistics"),

			// 5. 用户账号注销
			web.NSRouter("/account", userController, "delete:DeleteAccount"),
		),
	)
	web.AddNamespace(ns)
}
