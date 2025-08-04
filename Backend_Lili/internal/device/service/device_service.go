package service

import (
	"Backend_Lili/internal/device/model"
	"Backend_Lili/internal/device/repository"
	"Backend_Lili/pkg/utils"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type DeviceService struct {
	deviceRepo   *repository.DeviceRepository
	categoryRepo *repository.CategoryRepository
	templateRepo *repository.TemplateRepository
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		deviceRepo:   repository.NewDeviceRepository(),
		categoryRepo: repository.NewCategoryRepository(),
		templateRepo: repository.NewTemplateRepository(),
	}
}

// GetDevicesList 获取设备列表
func (s *DeviceService) GetDevicesList(userID int, req *GetDevicesListRequest) (*GetDevicesListResponse, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	// 构建查询参数
	params := make(map[string]interface{})
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // 限制最大页大小
	}

	params["page"] = req.Page
	params["limit"] = req.Limit
	if req.CategoryID > 0 {
		params["category_id"] = req.CategoryID
	}
	if req.Status != "" {
		params["status"] = req.Status
	}
	if req.Search != "" {
		params["search"] = strings.TrimSpace(req.Search)
	}
	if req.Sort != "" {
		// 验证排序字段
		validSorts := map[string]bool{
			"created_at":     true,
			"purchase_price": true,
			"current_value":  true,
			"purchase_date":  true,
		}
		if validSorts[req.Sort] {
			params["sort"] = req.Sort
			if req.Order == "asc" || req.Order == "desc" {
				params["order"] = req.Order
			}
		}
	}

	devices, total, err := s.deviceRepo.GetDevicesList(userID, params)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备列表失败")
	}

	// 计算分页信息
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &GetDevicesListResponse{
		Devices:    devices,
		Total:      int(total),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetDeviceDetail 获取设备详情
func (s *DeviceService) GetDeviceDetail(deviceID, userID int) (*model.Device, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	device, err := s.deviceRepo.GetDeviceByID(deviceID, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备详情失败")
	}
	if device == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备不存在")
	}

	// 获取价格历史，用于计算趋势
	priceHistories, err := s.deviceRepo.GetPriceHistory(deviceID, userID, 10)
	if err == nil && len(priceHistories) > 0 {
		// 计算最新市场价格作为当前估值
		latestPrice := priceHistories[0].Price
		if device.CurrentValue != latestPrice {
			s.deviceRepo.UpdateDeviceCurrentValue(deviceID, latestPrice)
			device.CurrentValue = latestPrice
		}
	}

	return device, nil
}

// CreateDevice 创建设备
func (s *DeviceService) CreateDevice(userID int, req *CreateDeviceRequest) (*model.Device, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	// 验证必填参数
	if err := s.validateCreateDeviceRequest(req); err != nil {
		return nil, err
	}

	// 验证模板和分类
	template, err := s.templateRepo.GetTemplateByID(req.TemplateID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "查询设备模板失败")
	}
	if template == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备模板不存在")
	}

	category, err := s.categoryRepo.GetCategoryByID(req.CategoryID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "查询设备分类失败")
	}
	if category == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备分类不存在")
	}

	// 解析日期
	purchaseDate, err := time.Parse("2006-01-02", req.PurchaseDate)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "购买日期格式错误")
	}

	// 构建设备对象
	device := &model.Device{
		UserID:        &userID,
		TemplateID:    &req.TemplateID,
		CategoryID:    &req.CategoryID,
		Name:          req.Name,
		Brand:         req.Brand,
		Model:         req.Model,
		SerialNumber:  req.SerialNumber,
		Color:         req.Color,
		Storage:       req.Storage,
		Memory:        req.Memory,
		Processor:     req.Processor,
		ScreenSize:    req.ScreenSize,
		PurchasePrice: req.PurchasePrice,
		CurrentValue:  req.PurchasePrice, // 初始估值等于购买价格
		PurchaseDate:  purchaseDate,
		Condition:     req.Condition,
		Status:        "active",
		Notes:         req.Notes,
	}

	// 处理保修日期
	if req.WarrantyDate != "" {
		warrantyDate, err := time.Parse("2006-01-02", req.WarrantyDate)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "保修日期格式错误")
		}
		device.WarrantyDate = warrantyDate
	}

	// 处理规格参数
	if req.Specifications != nil {
		specJSON, err := json.Marshal(req.Specifications)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "规格参数格式错误")
		}
		device.Specifications = string(specJSON)
	}

	// 创建设备
	err = s.deviceRepo.CreateDevice(device)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "创建设备失败")
	}

	// 添加初始价格记录
	initialPrice := &model.PriceHistory{
		DeviceID:    device.ID,
		Source:      "manual",
		Platform:    "initial",
		Price:       req.PurchasePrice,
		Condition:   req.Condition,
		Description: "初始购买价格",
		RecordDate:  purchaseDate,
	}
	s.deviceRepo.AddPriceHistory(initialPrice)

	return device, nil
}

// UpdateDevice 更新设备
func (s *DeviceService) UpdateDevice(deviceID, userID int, req *UpdateDeviceRequest) (*model.Device, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 获取现有设备
	device, err := s.deviceRepo.GetDeviceByID(deviceID, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备信息失败")
	}
	if device == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备不存在")
	}

	// 更新字段
	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Brand != "" {
		device.Brand = req.Brand
	}
	if req.Model != "" {
		device.Model = req.Model
	}
	if req.SerialNumber != "" {
		device.SerialNumber = req.SerialNumber
	}
	if req.Color != "" {
		device.Color = req.Color
	}
	if req.Storage != "" {
		device.Storage = req.Storage
	}
	if req.Memory != "" {
		device.Memory = req.Memory
	}
	if req.Processor != "" {
		device.Processor = req.Processor
	}
	if req.ScreenSize != "" {
		device.ScreenSize = req.ScreenSize
	}
	if req.PurchasePrice > 0 {
		device.PurchasePrice = req.PurchasePrice
	}
	if req.PurchaseDate != "" {
		purchaseDate, err := time.Parse("2006-01-02", req.PurchaseDate)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "购买日期格式错误")
		}
		device.PurchaseDate = purchaseDate
	}
	if req.WarrantyDate != "" {
		warrantyDate, err := time.Parse("2006-01-02", req.WarrantyDate)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "保修日期格式错误")
		}
		device.WarrantyDate = warrantyDate
	}
	if req.Condition != "" {
		device.Condition = req.Condition
	}
	if req.Notes != "" {
		device.Notes = req.Notes
	}
	if req.Specifications != nil {
		specJSON, err := json.Marshal(req.Specifications)
		if err != nil {
			return nil, utils.NewBusinessError(utils.ERROR_PARAM, "规格参数格式错误")
		}
		device.Specifications = string(specJSON)
	}

	// 更新设备
	err = s.deviceRepo.UpdateDevice(device)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "更新设备失败")
	}

	return device, nil
}

// DeleteDevice 删除设备
func (s *DeviceService) DeleteDevice(deviceID, userID int) error {
	if deviceID <= 0 || userID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备存在
	device, err := s.deviceRepo.GetDeviceByID(deviceID, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备信息失败")
	}
	if device == nil {
		return utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备不存在")
	}

	// 软删除设备
	err = s.deviceRepo.SoftDeleteDevice(deviceID, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_DATABASE, "删除设备失败")
	}

	return nil
}

// UpdateDeviceStatus 更新设备状态
func (s *DeviceService) UpdateDeviceStatus(deviceID, userID int, req *UpdateDeviceStatusRequest) error {
	if deviceID <= 0 || userID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证状态值
	validStatuses := map[string]bool{
		"active": true,
		"sold":   true,
		"broken": true,
		"lost":   true,
	}
	if !validStatuses[req.Status] {
		return utils.NewBusinessError(utils.ERROR_PARAM, "无效的设备状态")
	}

	var salePrice *float64
	var saleDate *time.Time

	// 如果是出售状态，验证出售信息
	if req.Status == "sold" {
		if req.SalePrice <= 0 {
			return utils.NewBusinessError(utils.ERROR_PARAM, "出售价格必须大于0")
		}
		if req.SaleDate == "" {
			return utils.NewBusinessError(utils.ERROR_PARAM, "出售日期不能为空")
		}

		parsedSaleDate, err := time.Parse("2006-01-02", req.SaleDate)
		if err != nil {
			return utils.NewBusinessError(utils.ERROR_PARAM, "出售日期格式错误")
		}

		salePrice = &req.SalePrice
		saleDate = &parsedSaleDate
	}

	// 更新设备状态
	err := s.deviceRepo.UpdateDeviceStatus(deviceID, userID, req.Status, salePrice, saleDate, req.Notes)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_DATABASE, "更新设备状态失败")
	}

	return nil
}

// GetDeviceValuation 获取设备价值评估
func (s *DeviceService) GetDeviceValuation(deviceID, userID int) (*DeviceValuationResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 获取设备信息
	device, err := s.deviceRepo.GetDeviceByID(deviceID, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备信息失败")
	}
	if device == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备不存在")
	}

	// 获取价格历史
	priceHistories, err := s.deviceRepo.GetPriceHistory(deviceID, userID, 30)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取价格历史失败")
	}

	// 计算贬值信息
	currentValue := device.CurrentValue
	if len(priceHistories) > 0 {
		currentValue = priceHistories[0].Price
	}

	depreciation := device.PurchasePrice - currentValue
	depreciationRate := 0.0
	if device.PurchasePrice > 0 {
		depreciationRate = (depreciation / device.PurchasePrice) * 100
	}

	// 计算持有天数
	holdingDays := int(time.Since(device.PurchaseDate).Hours() / 24)
	if holdingDays == 0 {
		holdingDays = 1
	}

	// 计算日均贬值
	dailyDepreciation := depreciation / float64(holdingDays)

	return &DeviceValuationResponse{
		DeviceID:          deviceID,
		PurchasePrice:     device.PurchasePrice,
		CurrentValue:      currentValue,
		Depreciation:      depreciation,
		DepreciationRate:  depreciationRate,
		HoldingDays:       holdingDays,
		DailyDepreciation: dailyDepreciation,
		LastUpdateTime:    time.Now(),
		PriceHistories:    priceHistories,
	}, nil
}

// BatchImportDevices 批量导入设备
func (s *DeviceService) BatchImportDevices(userID int, req *BatchImportDevicesRequest) (*BatchImportDevicesResponse, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}
	if len(req.Devices) == 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "设备列表不能为空")
	}
	if len(req.Devices) > 100 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "批量导入设备数量不能超过100个")
	}

	var devices []*model.Device
	var errors []string

	for i, deviceReq := range req.Devices {
		// 验证每个设备
		if err := s.validateCreateDeviceRequest(&deviceReq); err != nil {
			errors = append(errors, "第"+strconv.Itoa(i+1)+"个设备: "+err.Error())
			continue
		}

		// 解析日期
		purchaseDate, err := time.Parse("2006-01-02", deviceReq.PurchaseDate)
		if err != nil {
			errors = append(errors, "第"+strconv.Itoa(i+1)+"个设备: 购买日期格式错误")
			continue
		}

		device := &model.Device{
			UserID:        &userID,
			TemplateID:    &deviceReq.TemplateID,
			CategoryID:    &deviceReq.CategoryID,
			Name:          deviceReq.Name,
			Brand:         deviceReq.Brand,
			Model:         deviceReq.Model,
			SerialNumber:  deviceReq.SerialNumber,
			Color:         deviceReq.Color,
			Storage:       deviceReq.Storage,
			Memory:        deviceReq.Memory,
			Processor:     deviceReq.Processor,
			ScreenSize:    deviceReq.ScreenSize,
			PurchasePrice: deviceReq.PurchasePrice,
			CurrentValue:  deviceReq.PurchasePrice,
			PurchaseDate:  purchaseDate,
			Condition:     deviceReq.Condition,
			Status:        "active",
			Notes:         deviceReq.Notes,
		}

		// 处理保修日期
		if deviceReq.WarrantyDate != "" {
			warrantyDate, err := time.Parse("2006-01-02", deviceReq.WarrantyDate)
			if err != nil {
				errors = append(errors, "第"+strconv.Itoa(i+1)+"个设备: 保修日期格式错误")
				continue
			}
			device.WarrantyDate = warrantyDate
		}

		// 处理规格参数
		if deviceReq.Specifications != nil {
			specJSON, err := json.Marshal(deviceReq.Specifications)
			if err != nil {
				errors = append(errors, "第"+strconv.Itoa(i+1)+"个设备: 规格参数格式错误")
				continue
			}
			device.Specifications = string(specJSON)
		}

		devices = append(devices, device)
	}

	// 批量导入
	successCount, err := s.deviceRepo.BatchImportDevices(devices)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "批量导入设备失败")
	}

	return &BatchImportDevicesResponse{
		TotalCount:   len(req.Devices),
		SuccessCount: successCount,
		FailCount:    len(req.Devices) - successCount,
		Errors:       errors,
	}, nil
}

// GetDeviceImages 获取设备图片列表
func (s *DeviceService) GetDeviceImages(deviceID, userID int) ([]*model.DeviceImage, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	images, err := s.deviceRepo.GetDeviceImages(deviceID, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备图片失败")
	}

	return images, nil
}

// UploadDeviceImage 上传设备图片
func (s *DeviceService) UploadDeviceImage(deviceID, userID int, file multipart.File, fileHeader *multipart.FileHeader, req *UploadDeviceImageRequest) (*UploadDeviceImageResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备是否存在并属于当前用户
	device, err := s.deviceRepo.GetDeviceByID(deviceID, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备信息失败")
	}
	if device == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备不存在")
	}

	// 验证文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	// 从文件扩展名判断类型
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	contentType := fileHeader.Header.Get("Content-Type")

	if !allowedTypes[contentType] && !allowedTypes["image/"+strings.TrimPrefix(ext, ".")] {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "不支持的文件类型，仅支持 JPEG、PNG、GIF 格式")
	}

	// 验证文件大小 (限制5MB)
	if fileHeader.Size > 5*1024*1024 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "文件大小不能超过 5MB")
	}

	// 这里应该实现文件保存到云存储或本地存储的逻辑
	// 暂时使用简单的URL生成模拟
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("device_%d_%d_%s", deviceID, timestamp, fileHeader.Filename)
	imageURL := fmt.Sprintf("/uploads/devices/%s", fileName)

	// 创建图片记录
	deviceImage := &model.DeviceImage{
		DeviceID:  deviceID,
		ImageURL:  imageURL,
		ImageType: req.ImageType,
		SortOrder: req.SortOrder,
	}

	// 保存到数据库
	err = s.deviceRepo.AddDeviceImage(deviceImage)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "保存图片记录失败")
	}

	return &UploadDeviceImageResponse{
		ImageID:   deviceImage.ID,
		ImageURL:  deviceImage.ImageURL,
		ImageType: deviceImage.ImageType,
		SortOrder: deviceImage.SortOrder,
	}, nil
}

// DeleteDeviceImage 删除设备图片
func (s *DeviceService) DeleteDeviceImage(imageID, deviceID, userID int) error {
	if imageID <= 0 || deviceID <= 0 || userID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 删除图片记录
	err := s.deviceRepo.DeleteDeviceImage(imageID, deviceID, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_DATABASE, "删除图片失败")
	}

	return nil
}

// PredictDevicePrice 预测设备价格 - 实现简单的价格预测算法
func (s *DeviceService) PredictDevicePrice(deviceID, userID int, req *PricePredictionRequest) (*PricePredictionResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 设置默认预测天数
	if req.Days <= 0 {
		req.Days = 30
	}
	if req.Days > 365 {
		req.Days = 365 // 限制最多预测一年
	}

	// 获取设备信息
	device, err := s.deviceRepo.GetDeviceByID(deviceID, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取设备信息失败")
	}
	if device == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "设备不存在")
	}

	// 获取价格历史数据（获取更多数据用于分析）
	priceHistories, err := s.deviceRepo.GetPriceHistory(deviceID, userID, 100)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取价格历史失败")
	}

	if len(priceHistories) < 2 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "价格历史数据不足，无法进行预测")
	}

	// 执行价格预测算法
	prediction := s.calculatePricePrediction(device, priceHistories, req.Days)

	return prediction, nil
}

// calculatePricePrediction 计算价格预测
func (s *DeviceService) calculatePricePrediction(device *model.Device, histories []*model.PriceHistory, days int) *PricePredictionResponse {
	// 1. 数据预处理 - 按日期排序并提取价格序列
	prices := make([]float64, len(histories))
	dates := make([]time.Time, len(histories))

	for i := len(histories) - 1; i >= 0; i-- { // 反向遍历，使数据按时间正序
		prices[len(histories)-1-i] = histories[i].Price
		dates[len(histories)-1-i] = histories[i].RecordDate
	}

	// 2. 趋势分析
	trendAnalysis := s.analyzeTrend(prices, dates)

	// 3. 线性回归预测
	predictionPoints := s.linearRegressionPredict(prices, dates, days)

	// 4. 计算预测准确度
	accuracy := s.calculateAccuracy(prices)

	// 5. 构建响应
	return &PricePredictionResponse{
		DeviceID:         device.ID,
		CurrentValue:     device.CurrentValue,
		PredictionDays:   days,
		Algorithm:        "Linear Regression with Trend Analysis",
		Accuracy:         accuracy,
		PredictionPoints: predictionPoints,
		TrendAnalysis:    trendAnalysis,
		CreatedAt:        time.Now(),
	}
}

// analyzeTrend 分析价格趋势
func (s *DeviceService) analyzeTrend(prices []float64, dates []time.Time) *TrendAnalysis {
	if len(prices) < 2 {
		return &TrendAnalysis{
			Trend:           "unknown",
			TrendStrength:   "weak",
			VolatilityLevel: "unknown",
			DailyChangeRate: 0,
			Confidence:      0,
		}
	}

	// 计算总体趋势
	firstPrice := prices[0]
	lastPrice := prices[len(prices)-1]
	totalDays := dates[len(dates)-1].Sub(dates[0]).Hours() / 24
	if totalDays == 0 {
		totalDays = 1
	}

	// 日均变化率
	dailyChangeRate := ((lastPrice - firstPrice) / firstPrice) * 100 / totalDays

	// 趋势方向
	var trend string
	if dailyChangeRate > 0.1 {
		trend = "rising"
	} else if dailyChangeRate < -0.1 {
		trend = "falling"
	} else {
		trend = "stable"
	}

	// 趋势强度
	var trendStrength string
	absChangeRate := dailyChangeRate
	if absChangeRate < 0 {
		absChangeRate = -absChangeRate
	}

	if absChangeRate > 2.0 {
		trendStrength = "strong"
	} else if absChangeRate > 0.5 {
		trendStrength = "moderate"
	} else {
		trendStrength = "weak"
	}

	// 计算波动性
	volatility := s.calculateVolatility(prices)
	var volatilityLevel string
	if volatility > 0.15 {
		volatilityLevel = "high"
	} else if volatility > 0.05 {
		volatilityLevel = "medium"
	} else {
		volatilityLevel = "low"
	}

	// 计算置信度
	confidence := s.calculateConfidence(prices, volatility)

	return &TrendAnalysis{
		Trend:           trend,
		TrendStrength:   trendStrength,
		VolatilityLevel: volatilityLevel,
		DailyChangeRate: dailyChangeRate,
		Confidence:      confidence * 100,
	}
}

// calculateVolatility 计算价格波动性
func (s *DeviceService) calculateVolatility(prices []float64) float64 {
	if len(prices) < 2 {
		return 0
	}

	// 计算价格变化的标准差
	var sum float64
	var sumSquares float64
	n := len(prices) - 1

	for i := 1; i < len(prices); i++ {
		change := (prices[i] - prices[i-1]) / prices[i-1]
		sum += change
		sumSquares += change * change
	}

	mean := sum / float64(n)
	variance := (sumSquares / float64(n)) - (mean * mean)

	if variance < 0 {
		variance = 0
	}

	return variance // 返回方差作为波动性指标
}

// linearRegressionPredict 使用线性回归进行价格预测
func (s *DeviceService) linearRegressionPredict(prices []float64, dates []time.Time, days int) []*PredictionPoint {
	if len(prices) < 2 {
		return []*PredictionPoint{}
	}

	// 将日期转换为数值（以第一个日期为基准的天数）
	baseDate := dates[0]
	x := make([]float64, len(dates))
	for i, date := range dates {
		x[i] = date.Sub(baseDate).Hours() / 24
	}

	// 线性回归计算
	n := float64(len(prices))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for i := 0; i < len(prices); i++ {
		sumX += x[i]
		sumY += prices[i]
		sumXY += x[i] * prices[i]
		sumX2 += x[i] * x[i]
	}

	// 计算回归系数
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	// 生成预测点
	predictions := make([]*PredictionPoint, days)
	lastDate := dates[len(dates)-1]

	for i := 0; i < days; i++ {
		futureDate := lastDate.AddDate(0, 0, i+1)
		futureDays := futureDate.Sub(baseDate).Hours() / 24
		predictedPrice := slope*futureDays + intercept

		// 确保预测价格不为负数
		if predictedPrice < 0 {
			predictedPrice = 0
		}

		predictions[i] = &PredictionPoint{
			Date:  futureDate.Format("2006-01-02"),
			Price: predictedPrice,
		}
	}

	return predictions
}

// calculateAccuracy 计算预测准确度
func (s *DeviceService) calculateAccuracy(prices []float64) float64 {
	if len(prices) < 5 {
		return 0.5 // 数据不足时返回中等准确度
	}

	// 简单的准确度估算：基于数据量和波动性
	dataPoints := float64(len(prices))
	volatility := s.calculateVolatility(prices)

	// 数据点越多，波动性越小，准确度越高
	baseAccuracy := 0.3 + (dataPoints/100)*0.4 // 基础准确度 0.3-0.7
	volatilityPenalty := volatility * 2        // 波动性惩罚

	accuracy := baseAccuracy - volatilityPenalty
	if accuracy > 0.9 {
		accuracy = 0.9
	}
	if accuracy < 0.1 {
		accuracy = 0.1
	}

	return accuracy
}

// calculateConfidence 计算预测置信度
func (s *DeviceService) calculateConfidence(prices []float64, volatility float64) float64 {
	if len(prices) < 3 {
		return 0.3
	}

	// 基于数据量和波动性计算置信度
	dataFactor := float64(len(prices)) / 50.0 // 50个数据点为满分
	if dataFactor > 1.0 {
		dataFactor = 1.0
	}

	volatilityFactor := 1.0 - (volatility * 5.0) // 波动性越高，置信度越低
	if volatilityFactor < 0.1 {
		volatilityFactor = 0.1
	}

	confidence := (dataFactor + volatilityFactor) / 2.0
	return confidence
}

// 验证创建设备请求
func (s *DeviceService) validateCreateDeviceRequest(req *CreateDeviceRequest) error {
	if req.TemplateID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "设备模板ID不能为空")
	}
	if req.CategoryID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "设备分类ID不能为空")
	}
	if strings.TrimSpace(req.Name) == "" {
		return utils.NewBusinessError(utils.ERROR_PARAM, "设备名称不能为空")
	}
	if strings.TrimSpace(req.Brand) == "" {
		return utils.NewBusinessError(utils.ERROR_PARAM, "品牌不能为空")
	}
	if strings.TrimSpace(req.Model) == "" {
		return utils.NewBusinessError(utils.ERROR_PARAM, "型号不能为空")
	}
	if req.PurchasePrice <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "购买价格必须大于0")
	}
	if strings.TrimSpace(req.PurchaseDate) == "" {
		return utils.NewBusinessError(utils.ERROR_PARAM, "购买日期不能为空")
	}
	if req.Condition != "" {
		validConditions := map[string]bool{
			"new":  true,
			"good": true,
			"fair": true,
			"poor": true,
		}
		if !validConditions[req.Condition] {
			return utils.NewBusinessError(utils.ERROR_PARAM, "无效的设备状态")
		}
	}

	return nil
}
