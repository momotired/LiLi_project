package model

import "time"

// 微信登录请求
type WechatLoginRequest struct {
	Code          string `json:"code" valid:"Required"`
	EncryptedData string `json:"encryptedData"`
	IV            string `json:"iv"`
}

// 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" valid:"Required"`
}

// 登录响应
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	UserInfo     *UserInfo `json:"user_info"`
}

// Token验证响应
type TokenVerifyResponse struct {
	UserInfo      *UserInfo `json:"user_info"`
	ExpiresIn     int64     `json:"expires_in"`
	RemainingTime int64     `json:"remaining_time"`
}

// 用户信息
type UserInfo struct {
	ID       int    `json:"id"`
	OpenID   string `json:"openid"`
	NickName string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
}

// Token声明（用于JWT解析）
type TokenClaims struct {
	UserID    int    `json:"user_id"`
	OpenID    string `json:"openid"`
	ExpiresAt int64  `json:"exp"`
}

// Token黑名单
type TokenBlacklist struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// 用户会话
type UserSession struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	OpenID       string    `json:"openid"`
	SessionKey   string    `json:"session_key"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
