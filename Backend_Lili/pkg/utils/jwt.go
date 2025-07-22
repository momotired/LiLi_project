package utils

import (
	"errors"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID int    `json:"user_id"`
	OpenID string `json:"openid"`
	jwt.StandardClaims
}

// 生成JWT Token - 支持自定义过期时间
func GenerateToken(userID int, openID string, duration time.Duration) (string, error) {
	jwtSecret, _ := beego.AppConfig.String("jwt_secret")

	nowTime := time.Now()
	expireTime := nowTime.Add(duration)

	claims := Claims{
		UserID: userID,
		OpenID: openID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Backend_Lili",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtSecret))

	return token, err
}

// 兼容旧版本的 GenerateToken 函数
func GenerateTokenWithDefaultExpire(userID int, openID string) (string, error) {
	expireHours, _ := beego.AppConfig.Int("jwt_expire_hours")
	if expireHours == 0 {
		expireHours = 24 // 默认24小时
	}

	return GenerateToken(userID, openID, time.Duration(expireHours)*time.Hour)
}

// 解析JWT Token
func ParseToken(token string) (*Claims, error) {
	jwtSecret, _ := beego.AppConfig.String("jwt_secret")

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// 验证Token有效性
func ValidateToken(token string) (*Claims, error) {
	if token == "" {
		return nil, errors.New("token is required")
	}

	claims, err := ParseToken(token)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
