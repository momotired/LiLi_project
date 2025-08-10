package service

import (
	"Backend_Lili/internal/price/model"
	"time"
)

// ============= 价格信息相关类型 =============

// 获取设备价格响应
type GetDevicePriceResponse struct {
	DeviceID     int     `json:"device_id"`
	CurrentPrice float64 `json:"current_price"`
	MarketPrice  float64 `json:"market_price"`
	AveragePrice float64 `json:"average_price"`
	MinPrice     float64 `json:"min_price"`
	MaxPrice     float64 `json:"max_price"`
	PriceChange  float64 `json:"price_change"`
	ChangeRate   float64 `json:"change_rate"`
	TrendStatus  string  `json:"trend_status"` // rising/falling/stable
	LastUpdateAt string  `json:"last_update_at"`
}

// 价格历史请求
type GetPriceHistoryRequest struct {
	Period      string `json:"period" form:"period"`           // 7d/30d/90d/180d/1y
	Source      string `json:"source" form:"source"`           // 价格来源筛选
	Granularity string `json:"granularity" form:"granularity"` // day/week/month
}

// 价格历史响应
type GetPriceHistoryResponse struct {
	DeviceID   int                    `json:"device_id"`
	Period     string                 `json:"period"`
	Histories  []*model.PriceHistory  `json:"histories"`
	Statistics *PriceHistoryStatistics `json:"statistics"`
}

// 价格历史统计
type PriceHistoryStatistics struct {
	RecordCount  int     `json:"record_count"`
	MinPrice     float64 `json:"min_price"`
	MaxPrice     float64 `json:"max_price"`
	AveragePrice float64 `json:"average_price"`
	PriceRange   float64 `json:"price_range"`
	Volatility   float64 `json:"volatility"` // 价格波动性
}

// 价格趋势分析响应
type GetPriceTrendResponse struct {
	DeviceID      int                `json:"device_id"`
	TrendStatus   string             `json:"trend_status"`   // rising/falling/stable
	TrendStrength string             `json:"trend_strength"` // strong/moderate/weak
	ChangeRate    float64            `json:"change_rate"`    // 变化率
	Confidence    float64            `json:"confidence"`     // 置信度
	Analysis      *TrendAnalysisData `json:"analysis"`
}

// 趋势分析数据
type TrendAnalysisData struct {
	ShortTerm  *TrendInfo `json:"short_term"`  // 短期趋势（7天）
	MediumTerm *TrendInfo `json:"medium_term"` // 中期趋势（30天）
	LongTerm   *TrendInfo `json:"long_term"`   // 长期趋势（90天）
}

// 趋势信息
type TrendInfo struct {
	Direction    string  `json:"direction"`     // up/down/stable
	ChangeRate   float64 `json:"change_rate"`   // 变化率
	Volatility   float64 `json:"volatility"`    // 波动性
	Reliability  float64 `json:"reliability"`   // 可靠性
}

// 价格预测请求
type GetPricePredictionRequest struct {
	Period string `json:"period" form:"period"` // 30d/90d/180d
}

// 价格预测响应
type GetPricePredictionResponse struct {
	DeviceID       int                     `json:"device_id"`
	Period         string                  `json:"period"`
	PredictedPrice float64                 `json:"predicted_price"`
	Confidence     float64                 `json:"confidence"`
	Algorithm      string                  `json:"algorithm"`
	Factors        []*PredictionFactor     `json:"factors"`
	Predictions    []*PredictionDataPoint  `json:"predictions"`
	ValidUntil     time.Time               `json:"valid_until"`
}

// 预测数据点
type PredictionDataPoint struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
	Upper float64 `json:"upper"` // 置信区间上限
	Lower float64 `json:"lower"` // 置信区间下限
}

// 预测因素
type PredictionFactor struct {
	Factor string  `json:"factor"`
	Weight float64 `json:"weight"`
	Value  float64 `json:"value"`
	Impact string  `json:"impact"` // positive/negative/neutral
}

// 手动更新价格响应
type UpdatePriceResponse struct {
	DeviceID     int     `json:"device_id"`
	OldPrice     float64 `json:"old_price"`
	NewPrice     float64 `json:"new_price"`
	PriceChange  float64 `json:"price_change"`
	ChangeRate   float64 `json:"change_rate"`
	UpdatedAt    string  `json:"updated_at"`
	Source       string  `json:"source"`
}

// ============= 价格预警相关类型 =============

// 创建价格预警请求
type CreatePriceAlertRequest struct {
	AlertType           string   `json:"alert_type" valid:"Required"`           // price_drop/price_rise/target_price
	Threshold           float64  `json:"threshold" valid:"Required"`            // 预警阈值
	ThresholdType       string   `json:"threshold_type" valid:"Required"`       // percentage/absolute
	Enabled             bool     `json:"enabled"`                               // 是否启用，默认true
	NotificationMethods []string `json:"notification_methods"`                 // 通知方式数组
}

// 更新价格预警请求
type UpdatePriceAlertRequest struct {
	AlertType           string   `json:"alert_type"`
	Threshold           float64  `json:"threshold"`
	ThresholdType       string   `json:"threshold_type"`
	Enabled             *bool    `json:"enabled"` // 使用指针以区分false和未设置
	NotificationMethods []string `json:"notification_methods"`
}

// 价格预警列表请求
type GetPriceAlertsRequest struct {
	Status   string `json:"status" form:"status"`     // active/triggered/disabled
	DeviceID int    `json:"device_id" form:"device_id"` // 设备ID筛选
}

// 价格预警列表响应
type GetPriceAlertsResponse struct {
	Alerts []*PriceAlertInfo `json:"alerts"`
	Total  int               `json:"total"`
}

// 价格预警信息
type PriceAlertInfo struct {
	ID                  int       `json:"id"`
	DeviceID            int       `json:"device_id"`
	DeviceName          string    `json:"device_name"`
	AlertType           string    `json:"alert_type"`
	Threshold           float64   `json:"threshold"`
	ThresholdType       string    `json:"threshold_type"`
	Enabled             bool      `json:"enabled"`
	NotificationMethods []string  `json:"notification_methods"`
	LastTriggeredAt     *string   `json:"last_triggered_at"`
	TriggerCount        int       `json:"trigger_count"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}

// ============= 市场价格对比相关类型 =============

// 市场价格对比响应
type GetMarketComparisonResponse struct {
	DeviceID      int                      `json:"device_id"`
	CurrentPrice  float64                  `json:"current_price"`
	Comparisons   []*MarketPriceComparison `json:"comparisons"`
	BestPrice     *MarketPriceComparison   `json:"best_price"`
	PriceSummary  *MarketPriceSummary      `json:"price_summary"`
}

// 市场价格对比
type MarketPriceComparison struct {
	Platform    string  `json:"platform"`
	Price       float64 `json:"price"`
	Condition   string  `json:"condition"`
	URL         string  `json:"url"`
	RecordDate  string  `json:"record_date"`
	Reliability float64 `json:"reliability"`
}

// 市场价格汇总
type MarketPriceSummary struct {
	MinPrice     float64 `json:"min_price"`
	MaxPrice     float64 `json:"max_price"`
	AveragePrice float64 `json:"average_price"`
	MedianPrice  float64 `json:"median_price"`
	PriceRange   float64 `json:"price_range"`
	DataSources  int     `json:"data_sources"`
}

// ============= 价格数据源相关类型 =============

// 价格数据源响应
type GetPriceSourcesResponse struct {
	Sources []*PriceSourceInfo `json:"sources"`
}

// 价格数据源信息
type PriceSourceInfo struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Platform    string  `json:"platform"`
	Status      string  `json:"status"`
	Reliability float64 `json:"reliability"`
	UpdateFreq  int     `json:"update_freq"`
	LastSync    *string `json:"last_sync"`
}

// ============= 批量更新相关类型 =============

// 批量更新价格请求
type BatchUpdatePricesRequest struct {
	DeviceIDs []int `json:"device_ids" valid:"Required"`
}

// 批量更新价格响应
type BatchUpdatePricesResponse struct {
	TotalCount   int                    `json:"total_count"`
	SuccessCount int                    `json:"success_count"`
	FailCount    int                    `json:"fail_count"`
	Results      []*BatchUpdateResult   `json:"results"`
	Errors       []string               `json:"errors"`
}

// 批量更新结果
type BatchUpdateResult struct {
	DeviceID    int     `json:"device_id"`
	Success     bool    `json:"success"`
	OldPrice    float64 `json:"old_price"`
	NewPrice    float64 `json:"new_price"`
	PriceChange float64 `json:"price_change"`
	Error       string  `json:"error,omitempty"`
}