package router

import (
	"Backend_Lili/internal/auth/controller"
	"Backend_Lili/internal/auth/middleware"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// API版本前缀
	ns := beego.NewNamespace("/api/v1",
		// 认证相关路由
		beego.NSNamespace("/auth",
			// 不需要认证的接口
			beego.NSRouter("/login", &controller.AuthController{}, "post:Login"),
			beego.NSRouter("/refresh", &controller.AuthController{}, "post:RefreshToken"),

			// 需要认证的接口
			beego.NSRouter("/logout", &controller.AuthController{}, "post:Logout"),
			beego.NSRouter("/verify", &controller.AuthController{}, "get:VerifyToken"),
		),
	)

	// 为需要认证的接口添加中间件
	beego.InsertFilter("/api/v1/auth/logout", beego.BeforeRouter, middleware.JWTAuth)
	beego.InsertFilter("/api/v1/auth/verify", beego.BeforeRouter, middleware.JWTAuth)

	beego.AddNamespace(ns)

	// 健康检查路由
	beego.Router("/health", &controller.AuthController{}, "get:Get")
}
