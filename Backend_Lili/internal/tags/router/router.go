package router

import (
    "Backend_Lili/internal/auth/middleware"
    "Backend_Lili/internal/tags/controller"

    "github.com/beego/beego/v2/server/web"
)

// InitTagsRoutes 初始化标签模块路由
func InitTagsRoutes() {
    tagsController := controller.NewTagsController()

    ns := web.NewNamespace("/api/v1",
        web.NSNamespace("/tags",
            web.NSBefore(middleware.JWTAuth),

            web.NSRouter("/", tagsController, "get:ListTags;post:CreateTag"),
            web.NSRouter("/:tagId", tagsController, "get:GetTag;put:UpdateTag;delete:DeleteTag"),

            web.NSRouter("/system", tagsController, "get:GetSystemTags"),
            web.NSRouter("/custom", tagsController, "get:GetCustomTags"),
            web.NSRouter("/popular", tagsController, "get:GetPopularTags"),
            web.NSRouter("/search", tagsController, "get:SearchTags"),
            web.NSRouter("/categories", tagsController, "get:GetTagCategories"),
            web.NSRouter("/recommendations", tagsController, "get:GetRecommendations"),
            web.NSRouter("/statistics", tagsController, "get:GetTagStatistics"),
        ),
    )

    web.AddNamespace(ns)
}


