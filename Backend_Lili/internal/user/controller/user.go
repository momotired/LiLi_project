package controllers

import (
	"encoding/json"

	"Backend_Lili/internal/user/model"
	"Backend_Lili/pkg/utils"
)

type UserController struct {
	BaseController
}

// 更新用户信息请求结构
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
}

// 获取用户信息
func (c *UserController) GetProfile() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "authentication required")
		return
	}

	user, err := model.GetUserByOpenID(claims.OpenID)
	if err != nil {
		c.WriteError(utils.ERROR_NOT_FOUND, "user not found")
		return
	}

	c.WriteJSON(user)
}

// 更新用户信息
func (c *UserController) UpdateProfile() {
	claims, err := c.GetCurrentUser()
	if err != nil {
		c.WriteError(utils.ERROR_AUTH, "authentication required")
		return
	}

	var req UpdateUserRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.WriteError(utils.ERROR_PARAM, "invalid request parameters")
		return
	}

	user, err := models.GetUserByOpenID(claims.OpenID)
	if err != nil {
		c.WriteError(utils.ERROR_NOT_FOUND, "user not found")
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Gender >= 0 && req.Gender <= 2 {
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

	if err := models.UpdateUser(user); err != nil {
		c.WriteError(utils.ERROR_SERVER, "failed to update user")
		return
	}

	c.WriteJSON(user)
}
