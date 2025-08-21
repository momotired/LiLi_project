package controller

import (
	base "Backend_Lili/internal/auth/controller"
	"Backend_Lili/internal/device/service"
	"Backend_Lili/pkg/utils"
	"strconv"
)

type DeviceController struct {
	base.BaseController
	deviceService *service.DeviceService
}

func NewDeviceController() *DeviceController {
	return &DeviceController{}
}

func (c *DeviceController) Prepare() {
	// 调用父类的Prepare方法
	c.BaseController.Prepare()
	
	// 初始化设备服务
	c.deviceService = service.NewDeviceService()
}

// GetDevicesList 获取设备列表
// @router /devices [get]
func (c *DeviceController) GetDevicesList() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析查询参数
	req := &service.GetDevicesListRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	response, err := c.deviceService.GetDevicesList(userID, req)
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

// GetDeviceDetail 获取设备详情
// @router /devices/:deviceId [get]
func (c *DeviceController) GetDeviceDetail() {
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
	device, err := c.deviceService.GetDeviceDetail(deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, device)
}

// CreateDevice 创建设备
// @router /devices [post]
func (c *DeviceController) CreateDevice() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析请求参数
	req := &service.CreateDeviceRequest{}
	if err := c.BindJSON(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	device, err := c.deviceService.CreateDevice(userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, device)
}

// UpdateDevice 更新设备
// @router /devices/:deviceId [put]
func (c *DeviceController) UpdateDevice() {
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

	// 解析请求参数
	req := &service.UpdateDeviceRequest{}
	if err := c.BindJSON(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	device, err := c.deviceService.UpdateDevice(deviceID, userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, device)
}

// DeleteDevice 删除设备
// @router /devices/:deviceId [delete]
func (c *DeviceController) DeleteDevice() {
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
	err = c.deviceService.DeleteDevice(deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, nil)
}

// UpdateDeviceStatus 更新设备状态
// @router /devices/:deviceId/status [patch]
func (c *DeviceController) UpdateDeviceStatus() {
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

	// 解析请求参数
	req := &service.UpdateDeviceStatusRequest{}
	if err := c.BindJSON(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	err = c.deviceService.UpdateDeviceStatus(deviceID, userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, nil)
}

// GetDeviceValuation 获取设备价值评估
// @router /devices/:deviceId/valuation [get]
func (c *DeviceController) GetDeviceValuation() {
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
	valuation, err := c.deviceService.GetDeviceValuation(deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, valuation)
}

// BatchImportDevices 批量导入设备
// @router /devices/import [post]
func (c *DeviceController) BatchImportDevices() {
	// 从JWT中获取用户ID
	userID, ok := c.Ctx.Input.GetData("user_id").(int)
	if !ok {
		utils.WriteError(c.Ctx, utils.ERROR_AUTH, "用户认证失败")
		return
	}

	// 解析请求参数
	req := &service.BatchImportDevicesRequest{}
	if err := c.BindJSON(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	response, err := c.deviceService.BatchImportDevices(userID, req)
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

// GetDeviceImages 获取设备图片
// @router /devices/:deviceId/images [get]
func (c *DeviceController) GetDeviceImages() {
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
	images, err := c.deviceService.GetDeviceImages(deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, images)
}

// UploadDeviceImage 上传设备图片
// @router /devices/:deviceId/images [post]
func (c *DeviceController) UploadDeviceImage() {
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

	// 获取上传的文件
	file, header, err := c.GetFile("image")
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "获取上传文件失败")
		return
	}
	defer file.Close()

	// 解析其他参数
	imageType := c.GetString("image_type", "normal")
	sortOrder, _ := c.GetInt("sort_order", 0)

	// 构建请求对象
	req := &service.UploadDeviceImageRequest{
		ImageType: imageType,
		SortOrder: int(sortOrder),
	}

	// 调用服务层
	response, err := c.deviceService.UploadDeviceImage(deviceID, userID, file, header, req)
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

// DeleteDeviceImage 删除设备图片
// @router /devices/:deviceId/images/:imageId [delete]
func (c *DeviceController) DeleteDeviceImage() {
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

	// 获取图片ID
	imageIDStr := c.Ctx.Input.Param(":imageId")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "图片ID格式错误")
		return
	}

	// 调用服务层
	err = c.deviceService.DeleteDeviceImage(imageID, deviceID, userID)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, nil)
}

// PredictDevicePrice 预测设备价格
// @router /devices/:deviceId/prediction [get]
func (c *DeviceController) PredictDevicePrice() {
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
	req := &service.PricePredictionRequest{}
	if err := c.ParseForm(req); err != nil {
		utils.WriteError(c.Ctx, utils.ERROR_PARAM, "参数解析失败")
		return
	}

	// 调用服务层
	prediction, err := c.deviceService.PredictDevicePrice(deviceID, userID, req)
	if err != nil {
		if businessErr, ok := err.(*utils.BusinessError); ok {
			utils.WriteError(c.Ctx, businessErr.Code, businessErr.Message)
		} else {
			utils.WriteError(c.Ctx, utils.ERROR_SERVER, "服务器内部错误")
		}
		return
	}

	utils.WriteSuccess(c.Ctx, prediction)
}
