package controller

import (
	base "Backend_Lili/internal/auth/controller"
	"Backend_Lili/internal/price/service"
	"Backend_Lili/pkg/utils"
	"strconv"
)

type PriceController struct {
	base.BaseController
	priceService *service.PriceService
}

func NewPriceController() *PriceController {
	return &PriceController{}
}

func (c *PriceController) Prepare() {
	// 调用父类的Prepare方法
	c.BaseController.Prepare()
	
	// 初始化价格服务
	c.priceService = service.NewPriceService()
}

// GetDevicePrice 获取设备价格信息
// @router /prices/device/:deviceId [get]
func (c *PriceController) GetDevicePrice() {
    // 从JWT中获取用户ID
    userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取设备ID
	deviceIDStr := c.Ctx.Input.Param(":deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "设备ID格式错误")
		return
	}

	// 调用服务层
	response, err := c.priceService.GetDevicePrice(deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// GetPriceHistory 获取价格历史记录
// @router /prices/device/:deviceId/history [get]
func (c *PriceController) GetPriceHistory() {
    // 从JWT中获取用户ID
    userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取设备ID
	deviceIDStr := c.Ctx.Input.Param(":deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "设备ID格式错误")
		return
	}

	// 解析查询参数
	req := &service.GetPriceHistoryRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	response, err := c.priceService.GetPriceHistory(deviceID, userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// GetPriceTrend 获取价格趋势分析
// @router /prices/device/:deviceId/trend [get]
func (c *PriceController) GetPriceTrend() {
    // 从JWT中获取用户ID
    userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取设备ID
	deviceIDStr := c.Ctx.Input.Param(":deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "设备ID格式错误")
		return
	}

	// 调用服务层
	response, err := c.priceService.GetPriceTrend(deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// GetPricePrediction 获取价格预测
// @router /prices/device/:deviceId/prediction [get]
func (c *PriceController) GetPricePrediction() {
    // 从JWT中获取用户ID
    userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取设备ID
	deviceIDStr := c.Ctx.Input.Param(":deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "设备ID格式错误")
		return
	}

	// 解析查询参数
	req := &service.GetPricePredictionRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	response, err := c.priceService.GetPricePrediction(deviceID, userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// UpdatePrice 手动更新价格
// @router /prices/device/:deviceId/update [post]
func (c *PriceController) UpdatePrice() {
    // 从JWT中获取用户ID
    userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取设备ID
	deviceIDStr := c.Ctx.Input.Param(":deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "设备ID格式错误")
		return
	}

	// 调用服务层
	response, err := c.priceService.UpdatePrice(deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// CreatePriceAlert 设置价格预警
// @router /prices/device/:deviceId/alerts [post]
func (c *PriceController) CreatePriceAlert() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取设备ID
	deviceIDStr := c.Ctx.Input.Param(":deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "设备ID格式错误")
		return
	}

	// 解析请求参数
	req := &service.CreatePriceAlertRequest{}
	if err := c.BindJSON(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	alert, err := c.priceService.CreatePriceAlert(deviceID, userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, alert)
}

// GetPriceAlerts 获取价格预警列表
// @router /prices/alerts [get]
func (c *PriceController) GetPriceAlerts() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析查询参数
	req := &service.GetPriceAlertsRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	response, err := c.priceService.GetPriceAlerts(userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// UpdatePriceAlert 更新价格预警
// @router /prices/alerts/:alertId [put]
func (c *PriceController) UpdatePriceAlert() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取预警ID
	alertIDStr := c.Ctx.Input.Param(":alertId")
	alertID, err := strconv.Atoi(alertIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "预警ID格式错误")
		return
	}

	// 解析请求参数
	req := &service.UpdatePriceAlertRequest{}
	if err := c.BindJSON(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	err = c.priceService.UpdatePriceAlert(alertID, userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{"message": "更新成功"})
}

// DeletePriceAlert 删除价格预警
// @router /prices/alerts/:alertId [delete]
func (c *PriceController) DeletePriceAlert() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取预警ID
	alertIDStr := c.Ctx.Input.Param(":alertId")
	alertID, err := strconv.Atoi(alertIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "预警ID格式错误")
		return
	}

	// 调用服务层
	err = c.priceService.DeletePriceAlert(alertID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, map[string]interface{}{"message": "删除成功"})
}

// GetMarketComparison 获取市场价格对比
// @router /prices/device/:deviceId/comparison [get]
func (c *PriceController) GetMarketComparison() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("userID").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 获取设备ID
	deviceIDStr := c.Ctx.Input.Param(":deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "设备ID格式错误")
		return
	}

	// 调用服务层
	response, err := c.priceService.GetMarketComparison(deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// GetPriceSources 获取价格数据源列表
// @router /prices/sources [get]
func (c *PriceController) GetPriceSources() {
    // 从JWT中获取用户ID（验证用户身份）
    _, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 调用服务层
	response, err := c.priceService.GetPriceSources()
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}

// BatchUpdatePrices 批量更新价格
// @router /prices/batch-update [post]
func (c *PriceController) BatchUpdatePrices() {
    // 从JWT中获取用户ID
    userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析请求参数
	req := &service.BatchUpdatePricesRequest{}
	if err := c.BindJSON(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	response, err := c.priceService.BatchUpdatePrices(userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, response)
}