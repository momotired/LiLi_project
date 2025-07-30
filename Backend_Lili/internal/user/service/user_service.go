package service

import (
	"Backend_Lili/internal/user/model"
	"Backend_Lili/internal/user/repository"
	"Backend_Lili/pkg/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// 获取用户信息
func (s *UserService) GetUserProfile(openID string) (*model.User, error) {
	if openID == "" {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "openID不能为空")
	}

	user, err := s.userRepo.GetUserByOpenID(openID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "查询用户信息失败")
	}

	if user == nil {
		return nil, utils.NewBusinessError(utils.ERROR_USER_NOT_FOUND, "用户不存在")
	}

	if user.Status == 0 {
		return nil, utils.NewBusinessError(utils.ERROR_USER_DISABLED, "用户已被禁用")
	}

	return user, nil
}

// 更新用户信息
func (s *UserService) UpdateUserProfile(openID string, req *UpdateUserProfileRequest) (*model.User, error) {
	user, err := s.GetUserProfile(openID)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if req.Gender < 0 || req.Gender > 2 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "性别参数无效")
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Gender >= 0 {
		user.Gender = req.Gender
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.Province != "" {
		user.Province = req.Province
	}
	if req.Country != "" {
		user.Country = req.Country
	}

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "更新用户信息失败")
	}

	return user, nil
}

// 获取用户偏好设置
func (s *UserService) GetUserPreferences(userID int) (*model.UserPreferences, error) {
	prefs, err := s.userRepo.GetPreferencesByUserID(userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "查询用户偏好设置失败")
	}

	// 如果没有偏好设置，返回默认设置
	if prefs == nil {
		prefs = &model.UserPreferences{
			UserID:                  userID,
			NotificationEnabled:     false,
			PriceAlertEnabled:       false,
			WarrantyReminderEnabled: false,
			Theme:                   "light",
			Language:                "zh-CN",
		}
	}

	return prefs, nil
}

// 更新用户偏好设置
func (s *UserService) UpdateUserPreferences(userID int, req *UpdateUserPreferencesRequest) (*model.UserPreferences, error) {
	prefs, err := s.userRepo.GetPreferencesByUserID(userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "查询用户偏好设置失败")
	}

	// 如果不存在，创建新的偏好设置
	if prefs == nil {
		prefs = &model.UserPreferences{
			UserID:                  userID,
			NotificationEnabled:     false,
			PriceAlertEnabled:       false,
			WarrantyReminderEnabled: false,
			Theme:                   "light",
			Language:                "zh-CN",
		}
	}

	// 更新偏好设置
	if req.NotificationEnabled != nil {
		prefs.NotificationEnabled = *req.NotificationEnabled
	}
	if req.PriceAlertEnabled != nil {
		prefs.PriceAlertEnabled = *req.PriceAlertEnabled
	}
	if req.WarrantyReminderEnabled != nil {
		prefs.WarrantyReminderEnabled = *req.WarrantyReminderEnabled
	}
	if req.Theme != "" {
		if req.Theme != "light" && req.Theme != "dark" {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "主题参数无效")
		}
		prefs.Theme = req.Theme
	}
	if req.Language != "" {
		prefs.Language = req.Language
	}

	err = s.userRepo.CreateOrUpdatePreferences(prefs)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "更新用户偏好设置失败")
	}

	return prefs, nil
}

// 获取用户标签
func (s *UserService) GetUserTags(userID int) ([]model.Tag, error) {
	tags, err := s.userRepo.GetTagsByUserID(userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "查询用户标签失败")
	}

	return tags, nil
}

// 更新用户标签
func (s *UserService) UpdateUserTags(userID int, req *UpdateUserTagsRequest) ([]model.Tag, error) {
	if req.TagIDs == nil {
		req.TagIDs = []int{}
	}

	// 验证标签ID有效性
	if len(req.TagIDs) > 0 {
		valid, err := s.userRepo.ValidateTagIDs(req.TagIDs)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "验证标签ID失败")
		}
		if !valid {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "存在无效的标签ID")
		}
	}

	// 更新用户标签关联
	err := s.userRepo.UpdateUserTags(userID, req.TagIDs)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "更新用户标签失败")
	}

	// 返回更新后的标签列表
	return s.GetUserTags(userID)
}

// 获取用户统计信息
func (s *UserService) GetUserStatistics(userID int) (*UserStatistics, error) {
	// 获取设备总数
	deviceCount, err := s.userRepo.GetUserDeviceCount(userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备总数失败")
	}

	// 获取各类设备数量
	categoryCount, err := s.userRepo.GetUserDeviceCountByCategory(userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备分类统计失败")
	}

	// 获取设备总价值
	totalValue, err := s.userRepo.GetUserDevicesTotalValue(userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备总价值失败")
	}

	// 获取平均持有时间
	avgHoldDays, err := s.userRepo.GetUserDeviceAverageHoldTime(userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取平均持有时间失败")
	}

	statistics := &UserStatistics{
		TotalDevices:        deviceCount,
		DevicesByCategory:   categoryCount,
		TotalValue:          totalValue,
		AverageHoldTimeDays: avgHoldDays,
	}

	return statistics, nil
}

// 注销用户账号
func (s *UserService) DeleteUserAccount(openID string, req *DeleteUserAccountRequest) error {
	if !req.Confirm {
		return utils.NewBusinessError(utils.ERROR_PARAM, "必须确认注销操作")
	}

	user, err := s.GetUserProfile(openID)
	if err != nil {
		return err
	}

	// 软删除用户账号
	err = s.userRepo.SoftDeleteUser(user.ID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_DATABASE, "注销用户账号失败")
	}

	return nil
}

// 请求和响应结构体
type UpdateUserProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
}

type UpdateUserPreferencesRequest struct {
	NotificationEnabled     *bool  `json:"notification_enabled"`
	PriceAlertEnabled       *bool  `json:"price_alert_enabled"`
	WarrantyReminderEnabled *bool  `json:"warranty_reminder_enabled"`
	Theme                   string `json:"theme"`
	Language                string `json:"language"`
}

type UpdateUserTagsRequest struct {
	TagIDs []int `json:"tag_ids"`
}

type DeleteUserAccountRequest struct {
	Confirm bool `json:"confirm"`
}

type UserStatistics struct {
	TotalDevices        int            `json:"total_devices"`
	DevicesByCategory   map[string]int `json:"devices_by_category"`
	TotalValue          float64        `json:"total_value"`
	AverageHoldTimeDays int            `json:"average_hold_time_days"`
}
