package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

// 微信Code2Session响应结构
type WechatSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// 微信用户信息（解密后）
type WechatUserInfo struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	Language  string `json:"language"`
}

// 微信手机号信息（解密后）
type WechatPhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

// 通过Code获取用户OpenID和SessionKey
func GetWechatUserInfo(code string) (*WechatSessionResponse, error) {
	appID, _ := beego.AppConfig.String("wechat_app_id")
	appSecret, _ := beego.AppConfig.String("wechat_app_secret")

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result WechatSessionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wechat api error: %d, %s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}

// 解密微信 encryptedData
func DecryptWechatData(encryptedData, iv, sessionKey string) ([]byte, error) {
	if encryptedData == "" || iv == "" || sessionKey == "" {
		return nil, errors.New("encryptedData, iv, sessionKey are required")
	}

	// Base64 解码
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	// AES-128-CBC 解密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)

	// 去除 PKCS7 填充
	plainText := pkcs7Unpad(cipherText)
	return plainText, nil
}

// 解密微信用户信息
func DecryptWechatUserInfo(encryptedData, iv, sessionKey string) (*WechatUserInfo, error) {
	plainText, err := DecryptWechatData(encryptedData, iv, sessionKey)
	if err != nil {
		return nil, err
	}

	var userInfo WechatUserInfo
	err = json.Unmarshal(plainText, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// 解密微信手机号信息
func DecryptWechatPhoneInfo(encryptedData, iv, sessionKey string) (*WechatPhoneInfo, error) {
	plainText, err := DecryptWechatData(encryptedData, iv, sessionKey)
	if err != nil {
		return nil, err
	}

	var phoneInfo WechatPhoneInfo
	err = json.Unmarshal(plainText, &phoneInfo)
	if err != nil {
		return nil, err
	}

	return &phoneInfo, nil
}

// PKCS7 去填充
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}

	unpadding := int(data[length-1])
	if unpadding > length {
		return data
	}

	return data[:(length - unpadding)]
}
