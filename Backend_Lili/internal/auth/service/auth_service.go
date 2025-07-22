package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"Backend_Lili/internal/auth/model"
	"Backend_Lili/internal/auth/repository"
	"Backend_Lili/pkg/utils"

	"github.com/beego/beego/v2/core/logs"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		authRepo: repository.NewAuthRepository(),
	}
}

// 微信登录
func (s *AuthService) WechatLogin(req *model.WechatLoginRequest) (*model.LoginResponse, error) {
	// 1. 调用微信API获取用户信息
	wechatInfo, err := utils.GetWechatUserInfo(req.Code)
	if err != nil {
		logs.Error("获取微信用户信息失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_CODE_INVALID, "微信授权码无效或已过期")
	}

	// 2. 查询用户是否已存在
	user, err := s.authRepo.GetUserByOpenID(wechatInfo.OpenID)
	if err != nil {
		// 用户不存在，创建新用户
		user, err = s.authRepo.CreateUser(wechatInfo.OpenID)
		if err != nil {
			logs.Error("创建用户失败:", err)
			return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "用户创建失败")
		}
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, utils.NewBusinessError(utils.ERROR_USER_DISABLED, "用户已被禁用")
	}

	// 3. 处理加密数据（如果提供了encryptedData和iv）
	if req.EncryptedData != "" && req.IV != "" {
		err := s.processEncryptedData(user.ID, req.EncryptedData, req.IV, wechatInfo.SessionKey)
		if err != nil {
			logs.Error("处理加密数据失败:", err)
			// 不中断登录流程，仅记录错误
		}
	}

	// 4. 更新最后登录时间
	err = s.authRepo.UpdateLastLoginTime(user.ID)
	if err != nil {
		logs.Error("更新最后登录时间失败:", err)
		// 不中断登录流程，仅记录错误
	}

	// 5. 生成Token
	accessToken, err := utils.GenerateToken(user.ID, user.OpenID, time.Hour*24)
	if err != nil {
		logs.Error("生成AccessToken失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "Token生成失败")
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.OpenID, time.Hour*24*7)
	if err != nil {
		logs.Error("生成RefreshToken失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "Token生成失败")
	}

	// 6. 存储RefreshToken
	err = s.authRepo.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		logs.Error("存储RefreshToken失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "Token存储失败")
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 3600, // 24小时
		TokenType:    "Bearer",
		UserInfo:     user,
	}, nil
}

// 处理微信加密数据
func (s *AuthService) processEncryptedData(userID int, encryptedData, iv, sessionKey string) error {
	// 解密用户信息
	userInfo, err := utils.DecryptWechatUserInfo(encryptedData, iv, sessionKey)
	if err != nil {
		// 如果解密用户信息失败，尝试解密手机号
		phoneInfo, phoneErr := utils.DecryptWechatPhoneInfo(encryptedData, iv, sessionKey)
		if phoneErr != nil {
			return utils.NewBusinessError(utils.ERROR_DECRYPT_FAILED, "数据解密失败")
		}

		// 更新用户手机号
		updates := map[string]interface{}{
			"phone": phoneInfo.PhoneNumber,
		}
		return s.authRepo.UpdateUserInfo(userID, updates)
	}

	// 更新用户信息
	updates := map[string]interface{}{
		"nickname": userInfo.NickName,
		"avatar":   userInfo.AvatarUrl,
		"gender":   userInfo.Gender,
		"city":     userInfo.City,
		"province": userInfo.Province,
		"country":  userInfo.Country,
		"language": userInfo.Language,
	}

	if userInfo.UnionID != "" {
		updates["unionid"] = userInfo.UnionID
	}

	return s.authRepo.UpdateUserInfo(userID, updates)
}

// 刷新Token
func (s *AuthService) RefreshToken(refreshToken string) (*model.LoginResponse, error) {
	// 1. 验证RefreshToken
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		logs.Error("RefreshToken验证失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_TOKEN_INVALID, "RefreshToken无效或已过期")
	}

	// 2. 检查RefreshToken是否存在
	exists, err := s.authRepo.CheckRefreshToken(claims.UserID, refreshToken)
	if err != nil {
		logs.Error("检查RefreshToken失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "Token验证失败")
	}

	if !exists {
		return nil, utils.NewBusinessError(utils.ERROR_TOKEN_INVALID, "RefreshToken无效")
	}

	// 3. 获取用户信息
	user, err := s.authRepo.GetUserByID(claims.UserID)
	if err != nil {
		logs.Error("获取用户信息失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_USER_NOT_FOUND, "用户不存在")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, utils.NewBusinessError(utils.ERROR_USER_DISABLED, "用户已被禁用")
	}

	// 4. 生成新的Token
	newAccessToken, err := utils.GenerateToken(user.ID, user.OpenID, time.Hour*24)
	if err != nil {
		logs.Error("生成新AccessToken失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "Token生成失败")
	}

	newRefreshToken, err := utils.GenerateToken(user.ID, user.OpenID, time.Hour*24*7)
	if err != nil {
		logs.Error("生成新RefreshToken失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "Token生成失败")
	}

	// 5. 更新RefreshToken
	err = s.authRepo.UpdateRefreshToken(claims.UserID, refreshToken, newRefreshToken)
	if err != nil {
		logs.Error("更新RefreshToken失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "Token更新失败")
	}

	return &model.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    24 * 3600,
		TokenType:    "Bearer",
		UserInfo:     user,
	}, nil
}

// 登出
func (s *AuthService) Logout(userID int, token string) error {
	// 1. 将Token加入黑名单
	err := s.authRepo.AddTokenToBlacklist(token)
	if err != nil {
		logs.Error("添加Token到黑名单失败:", err)
		return utils.NewBusinessError(utils.ERROR_DATABASE, "登出失败")
	}

	// 2. 删除RefreshToken
	err = s.authRepo.DeleteRefreshToken(userID)
	if err != nil {
		logs.Error("删除RefreshToken失败:", err)
		return utils.NewBusinessError(utils.ERROR_DATABASE, "登出失败")
	}

	return nil
}

// 验证Token
func (s *AuthService) VerifyToken(token string) (*model.TokenVerifyResponse, error) {
	// 1. 验证Token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		logs.Error("Token验证失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_TOKEN_INVALID, "Token无效或已过期")
	}

	// 2. 检查Token是否在黑名单中
	if s.authRepo.IsTokenBlacklisted(token) {
		return nil, utils.NewBusinessError(utils.ERROR_TOKEN_INVALID, "Token已失效")
	}

	// 3. 获取用户信息
	user, err := s.authRepo.GetUserByID(claims.UserID)
	if err != nil {
		logs.Error("获取用户信息失败:", err)
		return nil, utils.NewBusinessError(utils.ERROR_USER_NOT_FOUND, "用户不存在")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, utils.NewBusinessError(utils.ERROR_USER_DISABLED, "用户已被禁用")
	}

	// 4. 计算剩余时间
	remainingTime := claims.ExpiresAt - time.Now().Unix()

	return &model.TokenVerifyResponse{
		UserInfo:      user,
		ExpiresIn:     claims.ExpiresAt,
		RemainingTime: remainingTime,
	}, nil
}

// 生成随机字符串
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
