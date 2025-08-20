package service

import (
	"Backend_Lili/internal/price/model"
	"Backend_Lili/internal/price/repository"
	"Backend_Lili/pkg/utils"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"
)

type PriceService struct {
	priceRepo *repository.PriceRepository
}

func NewPriceService() *PriceService {
	return &PriceService{
		priceRepo: repository.NewPriceRepository(),
	}
}

// GetDevicePrice 获取设备价格信息
func (s *PriceService) GetDevicePrice(deviceID, userID int) (*GetDevicePriceResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备归属权
	if !s.verifyDeviceOwnership(deviceID, userID) {
		return nil, utils.NewBusinessError(utils.ERROR_AUTH, "无权访问该设备")
	}

	price, err := s.priceRepo.GetDevicePrice(deviceID, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取价格信息失败")
	}

	if price == nil {
		return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "价格信息不存在")
	}

	response := &GetDevicePriceResponse{
		DeviceID:     price.DeviceID,
		CurrentPrice: price.CurrentPrice,
		MarketPrice:  price.MarketPrice,
		AveragePrice: price.AveragePrice,
		MinPrice:     price.MinPrice,
		MaxPrice:     price.MaxPrice,
		PriceChange:  price.PriceChange,
		ChangeRate:   price.ChangeRate,
		TrendStatus:  price.TrendStatus,
	}

	if !price.LastUpdateAt.IsZero() {
		response.LastUpdateAt = price.LastUpdateAt.Format("2006-01-02 15:04:05")
	}

	return response, nil
}

// GetPriceHistory 获取价格历史记录
func (s *PriceService) GetPriceHistory(deviceID, userID int, req *GetPriceHistoryRequest) (*GetPriceHistoryResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备归属权
	if !s.verifyDeviceOwnership(deviceID, userID) {
		return nil, utils.NewBusinessError(utils.ERROR_AUTH, "无权访问该设备")
	}

	// 设置默认值
	if req.Period == "" {
		req.Period = "30d"
	}
	if req.Granularity == "" {
		req.Granularity = "day"
	}

	params := map[string]interface{}{
		"period": req.Period,
	}
	if req.Source != "" {
		params["source"] = req.Source
	}

	histories, err := s.priceRepo.GetPriceHistory(deviceID, userID, params)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取价格历史失败")
	}

	// 计算统计信息
	statistics := s.calculatePriceStatistics(histories)

	// 根据粒度聚合数据
	if req.Granularity != "day" {
		histories = s.aggregateHistoryData(histories, req.Granularity)
	}

	return &GetPriceHistoryResponse{
		DeviceID:   deviceID,
		Period:     req.Period,
		Histories:  histories,
		Statistics: statistics,
	}, nil
}

// GetPriceTrend 获取价格趋势分析
func (s *PriceService) GetPriceTrend(deviceID, userID int) (*GetPriceTrendResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备归属权
	if !s.verifyDeviceOwnership(deviceID, userID) {
		return nil, utils.NewBusinessError(utils.ERROR_AUTH, "无权访问该设备")
	}

	// 获取不同时间段的价格历史
	shortTermParams := map[string]interface{}{"period": "7d"}
	mediumTermParams := map[string]interface{}{"period": "30d"}
	longTermParams := map[string]interface{}{"period": "90d"}

	shortTermHistories, _ := s.priceRepo.GetPriceHistory(deviceID, userID, shortTermParams)
	mediumTermHistories, _ := s.priceRepo.GetPriceHistory(deviceID, userID, mediumTermParams)
	longTermHistories, _ := s.priceRepo.GetPriceHistory(deviceID, userID, longTermParams)

	// 分析趋势
	shortTrend := s.analyzeTrend(shortTermHistories)
	mediumTrend := s.analyzeTrend(mediumTermHistories)
	longTrend := s.analyzeTrend(longTermHistories)

	// 综合判断主要趋势
	overallTrend := s.determineOverallTrend(shortTrend, mediumTrend, longTrend)

	return &GetPriceTrendResponse{
		DeviceID:      deviceID,
		TrendStatus:   overallTrend.Direction,
		TrendStrength: s.getTrendStrength(overallTrend.ChangeRate),
		ChangeRate:    overallTrend.ChangeRate,
		Confidence:    overallTrend.Reliability,
		Analysis: &TrendAnalysisData{
			ShortTerm:  shortTrend,
			MediumTerm: mediumTrend,
			LongTerm:   longTrend,
		},
	}, nil
}

// GetPricePrediction 获取价格预测
func (s *PriceService) GetPricePrediction(deviceID, userID int, req *GetPricePredictionRequest) (*GetPricePredictionResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备归属权
	if !s.verifyDeviceOwnership(deviceID, userID) {
		return nil, utils.NewBusinessError(utils.ERROR_AUTH, "无权访问该设备")
	}

	// 设置默认预测周期
	if req.Period == "" {
		req.Period = "30d"
	}

	// 检查是否有有效的预测缓存
	prediction, err := s.priceRepo.GetPricePrediction(deviceID, userID, req.Period)
	if err == nil && prediction != nil && prediction.ValidUntil.After(time.Now()) {
		return s.buildPredictionResponse(prediction), nil
	}

	// 获取历史数据进行预测
	params := map[string]interface{}{"period": "180d"} // 使用6个月历史数据
	histories, err := s.priceRepo.GetPriceHistory(deviceID, userID, params)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取历史数据失败")
	}

	if len(histories) < 5 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "历史数据不足，无法进行预测")
	}

	// 执行价格预测
	predictionResult := s.performPricePrediction(histories, req.Period)

	// 保存预测结果
	newPrediction := &model.PricePrediction{
		DeviceID:       deviceID,
		UserID:         userID,
		PredictionType: req.Period,
		PredictedPrice: predictionResult.PredictedPrice,
		Confidence:     predictionResult.Confidence,
		Algorithm:      predictionResult.Algorithm,
		ValidUntil:     time.Now().AddDate(0, 0, 7), // 预测有效期7天
	}

	if factorsJSON, err := json.Marshal(predictionResult.Factors); err == nil {
		newPrediction.Factors = string(factorsJSON)
	}

	s.priceRepo.CreatePricePrediction(newPrediction)

	return predictionResult, nil
}

// UpdatePrice 手动更新价格
func (s *PriceService) UpdatePrice(deviceID, userID int) (*UpdatePriceResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备归属权
	if !s.verifyDeviceOwnership(deviceID, userID) {
		return nil, utils.NewBusinessError(utils.ERROR_AUTH, "无权访问该设备")
	}

	// 获取当前价格
	currentPrice, _ := s.priceRepo.GetDevicePrice(deviceID, userID)
	oldPrice := 0.0
	if currentPrice != nil {
		oldPrice = currentPrice.CurrentPrice
	}

	// 调用价格抓取服务（这里简化处理）
	newPrice, err := s.fetchPriceFromSources(deviceID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_SERVER, "价格更新失败: "+err.Error())
	}

	// 更新价格信息
	price := &model.Price{
		DeviceID:     deviceID,
		UserID:       userID,
		CurrentPrice: newPrice,
		LastUpdateAt: time.Now(),
	}

	if currentPrice != nil {
		price.PriceChange = newPrice - oldPrice
		if oldPrice > 0 {
			price.ChangeRate = (newPrice - oldPrice) / oldPrice * 100
		}
		price.TrendStatus = s.determineTrendStatus(price.ChangeRate)
	}

	err = s.priceRepo.UpdateDevicePrice(price)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "更新价格失败")
	}

	// 记录价格历史
	history := &model.PriceHistory{
		DeviceID:    deviceID,
		UserID:      userID,
		Source:      "api",
		Price:       newPrice,
		Condition:   "unknown",
		Description: "自动更新",
		RecordDate:  time.Now(),
	}
	s.priceRepo.CreatePriceHistory(history)

	// 检查价格预警
	s.checkAndTriggerAlerts(deviceID, newPrice)

	return &UpdatePriceResponse{
		DeviceID:    deviceID,
		OldPrice:    oldPrice,
		NewPrice:    newPrice,
		PriceChange: newPrice - oldPrice,
		ChangeRate:  price.ChangeRate,
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Source:      "api",
	}, nil
}

// CreatePriceAlert 创建价格预警
func (s *PriceService) CreatePriceAlert(deviceID, userID int, req *CreatePriceAlertRequest) (*model.PriceAlert, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备归属权
	if !s.verifyDeviceOwnership(deviceID, userID) {
		return nil, utils.NewBusinessError(utils.ERROR_AUTH, "无权访问该设备")
	}

	// 验证预警参数
	if err := s.validateAlertRequest(req); err != nil {
		return nil, err
	}

	// 序列化通知方式
	notificationMethodsJSON := ""
	if len(req.NotificationMethods) > 0 {
		if data, err := json.Marshal(req.NotificationMethods); err == nil {
			notificationMethodsJSON = string(data)
		}
	}

	alert := &model.PriceAlert{
		DeviceID:            deviceID,
		UserID:              userID,
		AlertType:           req.AlertType,
		Threshold:           req.Threshold,
		ThresholdType:       req.ThresholdType,
		Enabled:             req.Enabled,
		NotificationMethods: notificationMethodsJSON,
		Status:              "active",
	}

	err := s.priceRepo.CreatePriceAlert(alert)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "创建价格预警失败")
	}

	return alert, nil
}

// GetPriceAlerts 获取价格预警列表
func (s *PriceService) GetPriceAlerts(userID int, req *GetPriceAlertsRequest) (*GetPriceAlertsResponse, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	params := map[string]interface{}{}
	if req.Status != "" {
		params["status"] = req.Status
	}
	if req.DeviceID > 0 {
		params["device_id"] = req.DeviceID
	}

	alerts, err := s.priceRepo.GetPriceAlerts(userID, params)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取价格预警列表失败")
	}

	alertInfos := make([]*PriceAlertInfo, 0, len(alerts))
	for _, alert := range alerts {
		alertInfo := &PriceAlertInfo{
			ID:            alert.ID,
			DeviceID:      alert.DeviceID,
			AlertType:     alert.AlertType,
			Threshold:     alert.Threshold,
			ThresholdType: alert.ThresholdType,
			Enabled:       alert.Enabled,
			TriggerCount:  alert.TriggerCount,
			Status:        alert.Status,
			CreatedAt:     alert.CreatedAt,
		}

		// 解析通知方式
		if alert.NotificationMethods != "" {
			var methods []string
			json.Unmarshal([]byte(alert.NotificationMethods), &methods)
			alertInfo.NotificationMethods = methods
		}

		// 格式化触发时间
		if !alert.LastTriggeredAt.IsZero() {
			triggeredAt := alert.LastTriggeredAt.Format("2006-01-02 15:04:05")
			alertInfo.LastTriggeredAt = &triggeredAt
		}

		alertInfos = append(alertInfos, alertInfo)
	}

	return &GetPriceAlertsResponse{
		Alerts: alertInfos,
		Total:  len(alertInfos),
	}, nil
}

// UpdatePriceAlert 更新价格预警
func (s *PriceService) UpdatePriceAlert(alertID, userID int, req *UpdatePriceAlertRequest) error {
	if alertID <= 0 || userID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 获取现有预警
	alert, err := s.priceRepo.GetPriceAlertByID(alertID, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_DATABASE, "获取价格预警失败")
	}
	if alert == nil {
		return utils.NewBusinessError(utils.ERROR_NOT_FOUND, "价格预警不存在")
	}

	// 更新字段
	if req.AlertType != "" {
		alert.AlertType = req.AlertType
	}
	if req.Threshold > 0 {
		alert.Threshold = req.Threshold
	}
	if req.ThresholdType != "" {
		alert.ThresholdType = req.ThresholdType
	}
	if req.Enabled != nil {
		alert.Enabled = *req.Enabled
	}
	if len(req.NotificationMethods) > 0 {
		if data, err := json.Marshal(req.NotificationMethods); err == nil {
			alert.NotificationMethods = string(data)
		}
	}

	err = s.priceRepo.UpdatePriceAlert(alert)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_DATABASE, "更新价格预警失败")
	}

	return nil
}

// DeletePriceAlert 删除价格预警
func (s *PriceService) DeletePriceAlert(alertID, userID int) error {
	if alertID <= 0 || userID <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	err := s.priceRepo.DeletePriceAlert(alertID, userID)
	if err != nil {
		return utils.NewBusinessError(utils.ERROR_DATABASE, "删除价格预警失败")
	}

	return nil
}

// GetMarketComparison 获取市场价格对比
func (s *PriceService) GetMarketComparison(deviceID, userID int) (*GetMarketComparisonResponse, error) {
	if deviceID <= 0 || userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效")
	}

	// 验证设备归属权
	if !s.verifyDeviceOwnership(deviceID, userID) {
		return nil, utils.NewBusinessError(utils.ERROR_AUTH, "无权访问该设备")
	}

	// 获取当前价格
	currentPrice, _ := s.priceRepo.GetDevicePrice(deviceID, userID)
	currentPriceValue := 0.0
	if currentPrice != nil {
		currentPriceValue = currentPrice.CurrentPrice
	}

	// 获取市场对比数据
	marketData, err := s.priceRepo.GetMarketComparison(deviceID, userID)
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取市场对比数据失败")
	}

	comparisons := make([]*MarketPriceComparison, 0, len(marketData))
	var bestPrice *MarketPriceComparison
	prices := make([]float64, 0, len(marketData))

	for _, data := range marketData {
		comparison := &MarketPriceComparison{
			Platform:    data.Platform,
			Price:       data.Price,
			Condition:   data.Condition,
			URL:         data.URL,
			RecordDate:  data.RecordDate.Format("2006-01-02"),
			Reliability: 0.8, // 简化处理
		}
		comparisons = append(comparisons, comparison)
		prices = append(prices, data.Price)

		// 找到最优价格
		if bestPrice == nil || data.Price < bestPrice.Price {
			bestPrice = comparison
		}
	}

	// 计算价格汇总
	summary := s.calculateMarketSummary(prices)

	return &GetMarketComparisonResponse{
		DeviceID:     deviceID,
		CurrentPrice: currentPriceValue,
		Comparisons:  comparisons,
		BestPrice:    bestPrice,
		PriceSummary: summary,
	}, nil
}

// GetPriceSources 获取价格数据源
func (s *PriceService) GetPriceSources() (*GetPriceSourcesResponse, error) {
	sources, err := s.priceRepo.GetPriceSources()
	if err != nil {
		return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取价格数据源失败")
	}

	sourceInfos := make([]*PriceSourceInfo, 0, len(sources))
	for _, source := range sources {
		sourceInfo := &PriceSourceInfo{
			ID:          source.ID,
			Name:        source.Name,
			Platform:    source.Platform,
			Status:      source.Status,
			Reliability: source.Reliability,
			UpdateFreq:  source.UpdateFreq,
		}

		if !source.LastSync.IsZero() {
			lastSync := source.LastSync.Format("2006-01-02 15:04:05")
			sourceInfo.LastSync = &lastSync
		}

		sourceInfos = append(sourceInfos, sourceInfo)
	}

	return &GetPriceSourcesResponse{
		Sources: sourceInfos,
	}, nil
}

// BatchUpdatePrices 批量更新价格
func (s *PriceService) BatchUpdatePrices(userID int, req *BatchUpdatePricesRequest) (*BatchUpdatePricesResponse, error) {
	if userID <= 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "用户ID无效")
	}

	if len(req.DeviceIDs) == 0 {
		return nil, utils.NewBusinessError(utils.ERROR_PARAM, "设备ID列表不能为空")
	}

	response := &BatchUpdatePricesResponse{
		TotalCount: len(req.DeviceIDs),
		Results:    make([]*BatchUpdateResult, 0, len(req.DeviceIDs)),
	}

	for _, deviceID := range req.DeviceIDs {
		result := &BatchUpdateResult{
			DeviceID: deviceID,
		}

		// 验证设备归属权
		if !s.verifyDeviceOwnership(deviceID, userID) {
			result.Success = false
			result.Error = "无权访问该设备"
			response.FailCount++
			response.Results = append(response.Results, result)
			continue
		}

		// 更新价格
		updateResult, err := s.UpdatePrice(deviceID, userID)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			response.FailCount++
		} else {
			result.Success = true
			result.OldPrice = updateResult.OldPrice
			result.NewPrice = updateResult.NewPrice
			result.PriceChange = updateResult.PriceChange
			response.SuccessCount++
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}

// ============= 私有辅助方法 =============

// verifyDeviceOwnership 验证设备归属权
func (s *PriceService) verifyDeviceOwnership(deviceID, userID int) bool {
	// 这里应该调用设备服务验证归属权
	// 简化处理，假设验证通过
	return true
}

// calculatePriceStatistics 计算价格统计信息
func (s *PriceService) calculatePriceStatistics(histories []*model.PriceHistory) *PriceHistoryStatistics {
	if len(histories) == 0 {
		return &PriceHistoryStatistics{}
	}

	prices := make([]float64, len(histories))
	for i, h := range histories {
		prices[i] = h.Price
	}

	sort.Float64s(prices)

	stats := &PriceHistoryStatistics{
		RecordCount: len(prices),
		MinPrice:    prices[0],
		MaxPrice:    prices[len(prices)-1],
		PriceRange:  prices[len(prices)-1] - prices[0],
	}

	// 计算平均价格
	sum := 0.0
	for _, price := range prices {
		sum += price
	}
	stats.AveragePrice = sum / float64(len(prices))

	// 计算价格波动性（标准差）
	variance := 0.0
	for _, price := range prices {
		variance += math.Pow(price-stats.AveragePrice, 2)
	}
	stats.Volatility = math.Sqrt(variance / float64(len(prices)))

	return stats
}

// aggregateHistoryData 根据粒度聚合历史数据
func (s *PriceService) aggregateHistoryData(histories []*model.PriceHistory, granularity string) []*model.PriceHistory {
	// 简化处理，实际应该按周/月聚合数据
	return histories
}

// analyzeTrend 分析价格趋势
func (s *PriceService) analyzeTrend(histories []*model.PriceHistory) *TrendInfo {
	if len(histories) < 2 {
		return &TrendInfo{
			Direction:   "stable",
			ChangeRate:  0,
			Volatility:  0,
			Reliability: 0,
		}
	}

	// 计算价格变化
	firstPrice := histories[len(histories)-1].Price // 最早的价格
	lastPrice := histories[0].Price                 // 最新的价格
	changeRate := (lastPrice - firstPrice) / firstPrice * 100

	direction := "stable"
	if changeRate > 1 {
		direction = "up"
	} else if changeRate < -1 {
		direction = "down"
	}

	// 计算波动性
	prices := make([]float64, len(histories))
	for i, h := range histories {
		prices[i] = h.Price
	}
	volatility := s.calculateVolatility(prices)

	return &TrendInfo{
		Direction:   direction,
		ChangeRate:  changeRate,
		Volatility:  volatility,
		Reliability: math.Min(float64(len(histories))/30.0, 1.0), // 数据点越多可靠性越高
	}
}

// determineOverallTrend 确定整体趋势
func (s *PriceService) determineOverallTrend(short, medium, long *TrendInfo) *TrendInfo {
	// 简化处理，以中期趋势为主
	return medium
}

// getTrendStrength 获取趋势强度
func (s *PriceService) getTrendStrength(changeRate float64) string {
	absRate := math.Abs(changeRate)
	if absRate >= 10 {
		return "strong"
	} else if absRate >= 5 {
		return "moderate"
	}
	return "weak"
}

// performPricePrediction 执行价格预测
func (s *PriceService) performPricePrediction(histories []*model.PriceHistory, period string) *GetPricePredictionResponse {
	// 简化的线性预测算法
	if len(histories) < 5 {
		return &GetPricePredictionResponse{
			PredictedPrice: 0,
			Confidence:     0,
			Algorithm:      "insufficient_data",
		}
	}

	// 计算预测天数
	days := 30
	switch period {
	case "90d":
		days = 90
	case "180d":
		days = 180
	}

	// 使用最近的价格计算趋势
	recentHistories := histories
	if len(histories) > 30 {
		recentHistories = histories[:30]
	}

	// 简单线性回归预测
	prices := make([]float64, len(recentHistories))
	for i, h := range recentHistories {
		prices[i] = h.Price
	}

	// 计算平均价格变化率
	totalChange := 0.0
	for i := 1; i < len(prices); i++ {
		if prices[i-1] > 0 {
			totalChange += (prices[i] - prices[i-1]) / prices[i-1]
		}
	}
	avgChangeRate := totalChange / float64(len(prices)-1)

	// 预测未来价格
	currentPrice := prices[0]
	predictedPrice := currentPrice * (1 + avgChangeRate*float64(days))

	// 生成预测数据点
	predictions := make([]*PredictionDataPoint, 0)
	for i := 1; i <= days; i += 7 { // 每周一个数据点
		price := currentPrice * (1 + avgChangeRate*float64(i))
		confidence := 0.1 * price // 置信区间为预测价格的10%
		predictions = append(predictions, &PredictionDataPoint{
			Date:  time.Now().AddDate(0, 0, i).Format("2006-01-02"),
			Price: price,
			Upper: price + confidence,
			Lower: price - confidence,
		})
	}

	return &GetPricePredictionResponse{
		Period:         period,
		PredictedPrice: predictedPrice,
		Confidence:     0.7, // 简化处理
		Algorithm:      "linear_regression",
		Predictions:    predictions,
		ValidUntil:     time.Now().AddDate(0, 0, 7),
	}
}

// buildPredictionResponse 构建预测响应
func (s *PriceService) buildPredictionResponse(prediction *model.PricePrediction) *GetPricePredictionResponse {
	response := &GetPricePredictionResponse{
		DeviceID:       prediction.DeviceID,
		Period:         prediction.PredictionType,
		PredictedPrice: prediction.PredictedPrice,
		Confidence:     prediction.Confidence,
		Algorithm:      prediction.Algorithm,
		ValidUntil:     prediction.ValidUntil,
	}

	// 解析影响因素
	if prediction.Factors != "" {
		var factors []*PredictionFactor
		json.Unmarshal([]byte(prediction.Factors), &factors)
		response.Factors = factors
	}

	return response
}

// fetchPriceFromSources 从数据源获取价格
func (s *PriceService) fetchPriceFromSources(deviceID int) (float64, error) {
	// 简化处理，返回模拟价格
	// 实际应该调用各种价格API
	return 1000.0 + float64(deviceID%100), nil
}

// determineTrendStatus 确定趋势状态
func (s *PriceService) determineTrendStatus(changeRate float64) string {
	if changeRate > 1 {
		return "rising"
	} else if changeRate < -1 {
		return "falling"
	}
	return "stable"
}

// checkAndTriggerAlerts 检查并触发价格预警
func (s *PriceService) checkAndTriggerAlerts(deviceID int, currentPrice float64) {
	alerts, err := s.priceRepo.CheckPriceAlerts(deviceID, currentPrice)
	if err != nil {
		return
	}

	for _, alert := range alerts {
		// 触发预警通知
		s.triggerAlert(alert, currentPrice)
	}
}

// triggerAlert 触发预警
func (s *PriceService) triggerAlert(alert *model.PriceAlert, currentPrice float64) {
	// 更新预警状态
	alert.Status = "triggered"
	alert.LastTriggeredAt = time.Now()
	alert.TriggerCount++
	s.priceRepo.UpdatePriceAlert(alert)

	// 发送通知（简化处理）
	fmt.Printf("价格预警触发: 设备%d, 当前价格%.2f\n", alert.DeviceID, currentPrice)
}

// validateAlertRequest 验证预警请求
func (s *PriceService) validateAlertRequest(req *CreatePriceAlertRequest) error {
	validAlertTypes := map[string]bool{
		"price_drop":   true,
		"price_rise":   true,
		"target_price": true,
	}
	if !validAlertTypes[req.AlertType] {
		return utils.NewBusinessError(utils.ERROR_PARAM, "无效的预警类型")
	}

	validThresholdTypes := map[string]bool{
		"percentage": true,
		"absolute":   true,
	}
	if !validThresholdTypes[req.ThresholdType] {
		return utils.NewBusinessError(utils.ERROR_PARAM, "无效的阈值类型")
	}

	if req.Threshold <= 0 {
		return utils.NewBusinessError(utils.ERROR_PARAM, "阈值必须大于0")
	}

	return nil
}

// calculateMarketSummary 计算市场价格汇总
func (s *PriceService) calculateMarketSummary(prices []float64) *MarketPriceSummary {
	if len(prices) == 0 {
		return &MarketPriceSummary{}
	}

	sort.Float64s(prices)

	summary := &MarketPriceSummary{
		MinPrice:    prices[0],
		MaxPrice:    prices[len(prices)-1],
		PriceRange:  prices[len(prices)-1] - prices[0],
		DataSources: len(prices),
	}

	// 计算平均价格
	sum := 0.0
	for _, price := range prices {
		sum += price
	}
	summary.AveragePrice = sum / float64(len(prices))

	// 计算中位数
	if len(prices)%2 == 0 {
		summary.MedianPrice = (prices[len(prices)/2-1] + prices[len(prices)/2]) / 2
	} else {
		summary.MedianPrice = prices[len(prices)/2]
	}

	return summary
}

// calculateVolatility 计算价格波动性
func (s *PriceService) calculateVolatility(prices []float64) float64 {
	if len(prices) < 2 {
		return 0
	}

	// 计算平均价格
	sum := 0.0
	for _, price := range prices {
		sum += price
	}
	avg := sum / float64(len(prices))

	// 计算标准差
	variance := 0.0
	for _, price := range prices {
		variance += math.Pow(price-avg, 2)
	}
	return math.Sqrt(variance / float64(len(prices)))
}