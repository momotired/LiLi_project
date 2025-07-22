package main

import (
	"Backend_Lili/internal/user/model"
	_ "Backend_Lili/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

// 空白导入 不直接调用包里的标识符 目的是触发包的init()函数执行
func init() {
	// 初始化数据库
	model.Init()
}

func main() {
	// 开发模式下开启ORM调试
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	// 启动服务器
	beego.Run()
}
