package controller

import (
	"strings"

	"lili-backend/lilibd/internal/auth/model"
	"lili-backend/lilibd/internal/auth/service"
	"lili-backend/lilibd/pkg/utils"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

type AuthController struct {
	beego.Controller
	authService *service.AuthService
}

func (c *AuthController) Prepare() {
	c.authService = service.NewAuthService()

	// 处理CORS
	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if c.Ctx.Request.Method == "OPTIONS" {
		c.Ctx.Output.SetStatus(200)
		c.StopRun()
		return
	}
}

// POST /auth/login - 微信登录
func (c *AuthController) Login() {
	var req model.WechatLoginRequest
	if err := c.ParseForm(&req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "请求参数格式错误")
		return
	}

	// 参数验证
	valid := validation.Validation{}
	if b, err := valid.Valid(&req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数验证失败")
		return
	} else if !b {
		// 收集具体的验证错误信息
		var errors []string
		for _, err := range valid.Errors {
			errors = append(errors, err.Message)
		}
		utils.WriteErrorWithDetails(c.Ctx, utils.ERROR_PARAM, "参数验证失败", map[string]interface{}{
			"errors": errors,
		})
		return
	}

	// 调用服务层进行登录
	loginResp, err := c.authService.WechatLogin(&req)
	if err != nil {
		logs.Error("微信登录失败:", err)
		utils.HandleBusinessError(c.Ctx, err)
		return
	}

	utils.WriteSuccess(c.Ctx, loginResp)
}

// POST /auth/refresh - 刷新Token
func (c *AuthController) RefreshToken() {
	var req model.RefreshTokenRequest
	if err := c.ParseForm(&req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "请求参数格式错误")
		return
	}

	// 参数验证
	valid := validation.Validation{}
	if b, err := valid.Valid(&req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数验证失败")
		return
	} else if !b {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "RefreshToken不能为空")
		return
	}

	// 调用服务层刷新Token
	loginResp, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		logs.Error("刷新Token失败:", err)
		utils.HandleBusinessError(c.Ctx, err)
		return
	}

	utils.WriteSuccess(c.Ctx, loginResp)
}

// POST /auth/logout - 登出
func (c *AuthController) Logout() {
	// 从中间件获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteErrorWithCode(c.Ctx, utils.ERROR_AUTH)
		return
	}

	// 获取Token
	token := getTokenFromHeader(c.Ctx)
	if token == "" {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "Token不能为空")
		return
	}

	// 调用服务层登出
	err := c.authService.Logout(userID, token)
	if err != nil {
		logs.Error("登出失败:", err)
		utils.HandleBusinessError(c.Ctx, err)
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]string{
		"message": "登出成功",
	})
}

// GET /auth/verify - 验证Token
func (c *AuthController) VerifyToken() {
	// 获取Token
	token := getTokenFromHeader(c.Ctx)
	if token == "" {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "Token不能为空")
		return
	}

	// 调用服务层验证Token
	verifyResp, err := c.authService.VerifyToken(token)
	if err != nil {
		logs.Error("Token验证失败:", err)
		utils.HandleBusinessError(c.Ctx, err)
		return
	}

	utils.WriteSuccess(c.Ctx, verifyResp)
}

// 从请求头获取Token
func getTokenFromHeader(ctx *context.Context) string {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Bearer token格式
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return authHeader
}
