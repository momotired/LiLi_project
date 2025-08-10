package router

import (
    "Backend_Lili/internal/auth/middleware"
    "Backend_Lili/internal/statistics/controller"

    "github.com/beego/beego/v2/server/web"
)

// InitStatisticsRoutes 初始化统计模块路由
func InitStatisticsRoutes() {
    statsController := controller.NewStatisticsController()

    ns := web.NewNamespace("/api/v1",
        web.NSNamespace("/statistics",
            web.NSBefore(middleware.JWTAuth),

            web.NSRouter("/dashboard", statsController, "get:GetDashboard"),
            web.NSRouter("/devices", statsController, "get:GetDevicesStatistics"),
            web.NSRouter("/value-analysis", statsController, "get:GetValueAnalysis"),
            web.NSRouter("/price-trends", statsController, "get:GetPriceTrends"),
            web.NSRouter("/brands", statsController, "get:GetBrandsStatistics"),
            web.NSRouter("/device-age", statsController, "get:GetDeviceAgeStatistics"),
            web.NSRouter("/depreciation", statsController, "get:GetDepreciationStatistics"),
            web.NSRouter("/spending", statsController, "get:GetSpendingStatistics"),
            web.NSRouter("/heatmap", statsController, "get:GetHeatmap"),
            web.NSRouter("/investment-return", statsController, "get:GetInvestmentReturn"),
            web.NSRouter("/custom", statsController, "post:PostCustomStatistics"),
            web.NSRouter("/export", statsController, "post:PostExportReport"),
            web.NSRouter("/insights", statsController, "get:GetInsights"),
            web.NSRouter("/comparison", statsController, "post:PostComparison"),
        ),
    )

    web.AddNamespace(ns)
}


