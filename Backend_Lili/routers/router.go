package routers

import (
	// 导入各模块的路由
	_ "Backend_Lili/internal/auth/router"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	// 全局中间件
	beego.InsertFilter("*", beego.BeforeRouter, cors)
}

// CORS中间件
func cors(ctx *context.Context) {
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if ctx.Request.Method == "OPTIONS" {
		ctx.Output.SetStatus(200)
		return
	}
}
