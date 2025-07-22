package middleware

import (
	"strings"

	"Backend_Lili/internal/auth/repository"
	"Backend_Lili/pkg/utils"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

// JWT认证中间件
func JWTAuth(ctx *context.Context) {
	// 获取Token
	token := getTokenFromHeader(ctx)
	if token == "" {
		utils.WriteError(ctx, utils.ERROR_AUTH, "Token不能为空")
		return
	}

	// 检查Token是否在黑名单中
	authRepo := repository.NewAuthRepository()
	if authRepo.IsTokenBlacklisted(token) {
		utils.WriteError(ctx, utils.ERROR_AUTH, "Token已失效")
		return
	}

	// 验证Token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		logs.Error("Token验证失败:", err)
		utils.WriteError(ctx, utils.ERROR_AUTH, "Token无效或已过期")
		return
	}

	// 将用户信息存储到上下文中
	ctx.Input.SetData("user_id", claims.UserID)
	ctx.Input.SetData("openid", claims.OpenID)
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

// 获取当前用户ID
func GetCurrentUserID(ctx *context.Context) int {
	userID := ctx.Input.GetData("user_id")
	if userID == nil {
		return 0
	}
	return userID.(int)
}

// 获取当前用户OpenID
func GetCurrentOpenID(ctx *context.Context) string {
	openID := ctx.Input.GetData("openid")
	if openID == nil {
		return ""
	}
	return openID.(string)
}

// 获取当前Token
func GetCurrentToken(ctx *context.Context) string {
	token := ctx.Input.GetData("token")
	if token == nil {
		return ""
	}
	return token.(string)
}
