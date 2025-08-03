package model

import (
	"github.com/beego/beego/v2/client/orm"
	userModel "Backend_Lili/internal/user/model"
)

func Init() {
	// 确保用户模型已经注册
	orm.RegisterModel(
		new(userModel.User),
	)
	
	// 注册设备模块的所有模型
	orm.RegisterModel(
		new(Device),
		new(DeviceTemplate),
		new(Category),
		new(DeviceImage),
		new(PriceHistory),
	)
}
