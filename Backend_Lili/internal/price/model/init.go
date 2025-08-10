package model

import (
	"github.com/beego/beego/v2/client/orm"
)

// 注册价格相关模型
func init() {
	// 注册模型
	orm.RegisterModel(
		new(Price),
		new(PriceHistory),
		new(PriceAlert),
		new(PriceSource),
		new(PricePrediction),
	)
}