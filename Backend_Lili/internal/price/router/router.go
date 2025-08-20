package router

import (
	"Backend_Lili/internal/auth/middleware"
	"Backend_Lili/internal/price/controller"

	"github.com/beego/beego/v2/server/web"
)

// InitPriceRoutes 初始化价格模块路由
func InitPriceRoutes() {
	// 创建价格控制器实例
	priceController := controller.NewPriceController()

	// 创建价格路由组
    priceGroup := web.NewNamespace("/api/v1/prices",
        // 应用认证中间件
        web.NSBefore(middleware.JWTAuth),

		// 设备价格相关路由
		web.NSRouter("/device/:deviceId", priceController, "get:GetDevicePrice"),                 // 获取设备价格信息
		web.NSRouter("/device/:deviceId/history", priceController, "get:GetPriceHistory"),        // 获取价格历史记录
		web.NSRouter("/device/:deviceId/trend", priceController, "get:GetPriceTrend"),            // 获取价格趋势分析
		web.NSRouter("/device/:deviceId/prediction", priceController, "get:GetPricePrediction"),  // 获取价格预测
		web.NSRouter("/device/:deviceId/update", priceController, "post:UpdatePrice"),            // 手动更新价格
		web.NSRouter("/device/:deviceId/comparison", priceController, "get:GetMarketComparison"), // 获取市场价格对比

		// 价格预警相关路由
		web.NSRouter("/device/:deviceId/alerts", priceController, "post:CreatePriceAlert"), // 设置价格预警
		web.NSRouter("/alerts", priceController, "get:GetPriceAlerts"),                     // 获取价格预警列表
		web.NSRouter("/alerts/:alertId", priceController, "put:UpdatePriceAlert"),          // 更新价格预警
		web.NSRouter("/alerts/:alertId", priceController, "delete:DeletePriceAlert"),       // 删除价格预警

		// 其他功能路由
		web.NSRouter("/sources", priceController, "get:GetPriceSources"),         // 获取价格数据源列表
		web.NSRouter("/batch-update", priceController, "post:BatchUpdatePrices"), // 批量更新价格
	)

	// 注册路由组
	web.AddNamespace(priceGroup)
}
