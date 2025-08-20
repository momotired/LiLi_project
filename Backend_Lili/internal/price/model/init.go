package model

import (
	"github.com/beego/beego/v2/client/orm"
)


// Init 初始化价格模块数据模型
func Init() {
	// 注册模型
	orm.RegisterModel(
		new(Price),
		new(PriceHistory),
		new(PriceAlert),
		new(PriceSource),
		new(PricePrediction),
	)
}