package controller

import (
	"encoding/json"

	"Backend_Lili/internal/auth/controller"
	"Backend_Lili/internal/user/service"
	"Backend_Lili/pkg/utils"
)

type UserController struct {
	controller.BaseController
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// 1. 获取用户信息
// GET /users/profile
func (c *UserController) GetProfile() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "认证失败")
		return
	}

	user, err := c.userService.GetUserProfile(claims.OpenID)
	if err != nil {
		if bizErr, ok := err.(*utils.BusinessError); ok {
			c.WriteError(bizErr.Code, bizErr.Message)
		} else {
			c.WriteError(utils.ERROR_SERVER, "获取用户信息失败")
		}
		return
	}

	c.WriteJSON(user)
}

// 2. 更新用户信息
// PUT /users/profile
func (c *UserController) UpdateProfile() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "认证失败")
		return
	}

	var req service.UpdateUserProfileRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.WriteError(utils.ERROR_PARAM, "请求参数格式错误")
		return
	}

	user, err := c.userService.UpdateUserProfile(claims.OpenID, &req)
	if err != nil {
		if bizErr, ok := err.(*utils.BusinessError); ok {
			c.WriteError(bizErr.Code, bizErr.Message)
		} else {
			c.WriteError(utils.ERROR_SERVER, "更新用户信息失败")
		}
		return
	}

	c.WriteJSON(user)
}

// 3. 获取用户偏好设置
// GET /users/preferences
func (c *UserController) GetPreferences() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "认证失败")
		return
	}

	prefs, err := c.userService.GetUserPreferences(claims.UserID)
	if err != nil {
		if bizErr, ok := err.(*utils.BusinessError); ok {
			c.WriteError(bizErr.Code, bizErr.Message)
		} else {
			c.WriteError(utils.ERROR_SERVER, "获取用户偏好设置失败")
		}
		return
	}

	c.WriteJSON(prefs)
}

// 4. 更新用户偏好设置
// PUT /users/preferences
func (c *UserController) UpdatePreferences() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "认证失败")
		return
	}

	var req service.UpdateUserPreferencesRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.WriteError(utils.ERROR_PARAM, "请求参数格式错误")
		return
	}

	prefs, err := c.userService.UpdateUserPreferences(claims.UserID, &req)
	if err != nil {
		if bizErr, ok := err.(*utils.BusinessError); ok {
			c.WriteError(bizErr.Code, bizErr.Message)
		} else {
			c.WriteError(utils.ERROR_SERVER, "更新用户偏好设置失败")
		}
		return
	}

	c.WriteJSON(prefs)
}

// 5. 获取用户标签
// GET /users/tags
func (c *UserController) GetTags() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "认证失败")
		return
	}

	tags, err := c.userService.GetUserTags(claims.UserID)
	if err != nil {
		if bizErr, ok := err.(*utils.BusinessError); ok {
			c.WriteError(bizErr.Code, bizErr.Message)
		} else {
			c.WriteError(utils.ERROR_SERVER, "获取用户标签失败")
		}
		return
	}

	c.WriteJSON(tags)
}

// 6. 更新用户标签
// PUT /users/tags
func (c *UserController) UpdateTags() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "认证失败")
		return
	}

	var req service.UpdateUserTagsRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.WriteError(utils.ERROR_PARAM, "请求参数格式错误")
		return
	}

	tags, err := c.userService.UpdateUserTags(claims.UserID, &req)
	if err != nil {
		if bizErr, ok := err.(*utils.BusinessError); ok {
			c.WriteError(bizErr.Code, bizErr.Message)
		} else {
			c.WriteError(utils.ERROR_SERVER, "更新用户标签失败")
		}
		return
	}

	c.WriteJSON(tags)
}

// 7. 获取用户统计信息
// GET /users/statistics
func (c *UserController) GetStatistics() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "认证失败")
		return
	}

	statistics, err := c.userService.GetUserStatistics(claims.UserID)
	if err != nil {
		if bizErr, ok := err.(*utils.BusinessError); ok {
			c.WriteError(bizErr.Code, bizErr.Message)
		} else {
			c.WriteError(utils.ERROR_SERVER, "获取用户统计信息失败")
		}
		return
	}

	c.WriteJSON(statistics)
}

// 8. 注销用户账号
// DELETE /users/account
func (c *UserController) DeleteAccount() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "认证失败")
		return
	}

	var req service.DeleteUserAccountRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.WriteError(utils.ERROR_PARAM, "请求参数格式错误")
		return
	}

	err = c.userService.DeleteUserAccount(claims.OpenID, &req)
	if err != nil {
		if bizErr, ok := err.(*utils.BusinessError); ok {
			c.WriteError(bizErr.Code, bizErr.Message)
		} else {
			c.WriteError(utils.ERROR_SERVER, "注销用户账号失败")
		}
		return
	}

	c.WriteJSON(map[string]interface{}{
		"message": "用户账号已成功注销",
	})
}
