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
	ID        int       `orm:"column(id);auto;pk" json:"id"`
	Token     string    `orm:"column(token);size(1000)" json:"token"`
	ExpiresAt time.Time `orm:"column(expires_at);type(datetime)" json:"expires_at"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
}

func (t *TokenBlacklist) TableName() string {
	return "token_blacklist"
}

// 用户会话
type UserSession struct {
	ID           int       `orm:"column(id);auto;pk" json:"id"`
	UserID       int       `orm:"column(user_id)" json:"user_id"`
	OpenID       string    `orm:"column(openid);size(100);null" json:"openid"`
	SessionKey   string    `orm:"column(session_key);size(100);null" json:"session_key"`
	AccessToken  string    `orm:"column(access_token);size(1000);null" json:"access_token"`
	RefreshToken string    `orm:"column(refresh_token);size(1000);null" json:"refresh_token"`
	ExpiresAt    time.Time `orm:"column(expires_at);null;type(datetime)" json:"expires_at"`
	LastLoginAt  time.Time `orm:"column(last_login_at);null;type(datetime)" json:"last_login_at"`
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (u *UserSession) TableName() string {
	return "user_session"
}
