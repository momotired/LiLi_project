package controller

import (
	"strings"

	"Backend_Lili/pkg/utils"

	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

// 预处理 - 处理CORS
func (c *BaseController) Prepare() {
	// 处理跨域
	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if c.Ctx.Request.Method == "OPTIONS" {
		c.Ctx.Output.SetStatus(200)
		c.StopRun()
		return
	}
}

// 获取当前用户信息
func (c *BaseController) GetCurrentUser() (*utils.Claims, error) {
	authHeader := c.Ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, utils.Error(utils.ERROR_AUTH, "authorization header required")
	}

	// 检查Bearer前缀
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, utils.Error(utils.ERROR_AUTH, "invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, utils.Error(utils.ERROR_AUTH, "invalid token")
	}

	return claims, nil
}

// 输出JSON响应
func (c *BaseController) WriteJSON(data interface{}) {
	utils.WriteJSON(c.Ctx, utils.Success(data))
}

// 输出错误响应
func (c *BaseController) WriteError(code int, message string) {
	utils.WriteJSON(c.Ctx, utils.Error(code, message))
}
