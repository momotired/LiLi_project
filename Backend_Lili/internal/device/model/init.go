package model

import (
	"github.com/beego/beego/v2/client/orm"
)

func Init() {
	
	
	// 注册设备模块的所有模型
	orm.RegisterModel(
		new(Device),
		new(DeviceTemplate),
		new(Category),
		new(DeviceImage),
		new(PriceHistory),
	)
}
