package service

import (
	"Backend_Lili/internal/device/model"
	"time"
)

// 设备列表请求
type GetDevicesListRequest struct {
	Page       int    `json:"page" form:"page"`             // 页码
	Limit      int    `json:"limit" form:"limit"`           // 每页数量
	CategoryID int    `json:"category_id" form:"category_id"` // 分类ID
	Status     string `json:"status" form:"status"`         // 设备状态
	Sort       string `json:"sort" form:"sort"`             // 排序字段
	Order      string `json:"order" form:"order"`           // 排序方式
	Search     string `json:"search" form:"search"`         // 搜索关键词
}

// 设备列表响应
type GetDevicesListResponse struct {
	Devices    []*model.Device `json:"devices"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// 创建设备请求
type CreateDeviceRequest struct {
	TemplateID      int                    `json:"template_id" valid:"Required"`
	Name            string                 `json:"name" valid:"Required"`
	Brand           string                 `json:"brand" valid:"Required"`
	Model           string                 `json:"model" valid:"Required"`
	CategoryID      int                    `json:"category_id" valid:"Required"`
	PurchasePrice   float64                `json:"purchase_price" valid:"Required"`
	PurchaseDate    string                 `json:"purchase_date" valid:"Required"` // YYYY-MM-DD格式
	WarrantyDate    string                 `json:"warranty_date"`                  // YYYY-MM-DD格式
	SerialNumber    string                 `json:"serial_number"`
	Color           string                 `json:"color"`
	Storage         string                 `json:"storage"`
	Memory          string                 `json:"memory"`
	Processor       string                 `json:"processor"`
	ScreenSize      string                 `json:"screen_size"`
	Condition       string                 `json:"condition"`       // new/good/fair/poor
	Notes           string                 `json:"notes"`
	Images          []string               `json:"images"`          // 图片URL数组
	Specifications  map[string]interface{} `json:"specifications"`  // 其他规格参数
}

// 更新设备请求
type UpdateDeviceRequest struct {
	Name            string                 `json:"name"`
	Brand           string                 `json:"brand"`
	Model           string                 `json:"model"`
	SerialNumber    string                 `json:"serial_number"`
	Color           string                 `json:"color"`
	Storage         string                 `json:"storage"`
	Memory          string                 `json:"memory"`
	Processor       string                 `json:"processor"`
	ScreenSize      string                 `json:"screen_size"`
	PurchasePrice   float64                `json:"purchase_price"`
	PurchaseDate    string                 `json:"purchase_date"`   // YYYY-MM-DD格式
	WarrantyDate    string                 `json:"warranty_date"`   // YYYY-MM-DD格式
	Condition       string                 `json:"condition"`       // new/good/fair/poor
	Notes           string                 `json:"notes"`
	Specifications  map[string]interface{} `json:"specifications"`  // 其他规格参数
}

// 更新设备状态请求
type UpdateDeviceStatusRequest struct {
	Status    string  `json:"status" valid:"Required"`    // active/sold/broken/lost
	SalePrice float64 `json:"sale_price"`                 // 出售价格，当status为sold时必填
	SaleDate  string  `json:"sale_date"`                  // 出售日期，当status为sold时必填
	Notes     string  `json:"notes"`                      // 状态变更备注
}

// 设备价值评估响应
type DeviceValuationResponse struct {
	DeviceID          int                    `json:"device_id"`
	PurchasePrice     float64                `json:"purchase_price"`
	CurrentValue      float64                `json:"current_value"`
	Depreciation      float64                `json:"depreciation"`       // 贬值金额
	DepreciationRate  float64                `json:"depreciation_rate"`  // 贬值率(%)
	HoldingDays       int                    `json:"holding_days"`       // 持有天数
	DailyDepreciation float64                `json:"daily_depreciation"` // 日均贬值
	LastUpdateTime    time.Time              `json:"last_update_time"`
	PriceHistories    []*model.PriceHistory  `json:"price_histories"`
}

// 批量导入设备请求
type BatchImportDevicesRequest struct {
	Devices           []CreateDeviceRequest `json:"devices" valid:"Required"`
	IgnoreDuplicates  bool                  `json:"ignore_duplicates"` // 是否忽略重复设备
}

// 批量导入设备响应
type BatchImportDevicesResponse struct {
	TotalCount   int      `json:"total_count"`
	SuccessCount int      `json:"success_count"`
	FailCount    int      `json:"fail_count"`
	Errors       []string `json:"errors"`
}

// 上传设备图片请求
type UploadDeviceImageRequest struct {
	ImageType string `json:"image_type"` // normal/cover
	SortOrder int    `json:"sort_order"`
}

// 上传设备图片响应
type UploadDeviceImageResponse struct {
	ImageID   int    `json:"image_id"`
	ImageURL  string `json:"image_url"`
	ImageType string `json:"image_type"`
	SortOrder int    `json:"sort_order"`
}

// 价格预测请求
type PricePredictionRequest struct {
	Days int `json:"days" form:"days"` // 预测未来多少天的价格，默认30天
}

// 价格预测数据点
type PredictionPoint struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

// 价格预测响应
type PricePredictionResponse struct {
	DeviceID         int                `json:"device_id"`
	CurrentValue     float64            `json:"current_value"`
	PredictionDays   int                `json:"prediction_days"`
	Algorithm        string             `json:"algorithm"`        // 使用的算法名称
	Accuracy         float64            `json:"accuracy"`         // 预测准确度评估 (0-1)
	PredictionPoints []*PredictionPoint `json:"prediction_points"` // 预测数据点
	TrendAnalysis    *TrendAnalysis     `json:"trend_analysis"`   // 趋势分析
	CreatedAt        time.Time          `json:"created_at"`
}

// 趋势分析
type TrendAnalysis struct {
	Trend           string  `json:"trend"`            // 趋势方向：rising/falling/stable
	TrendStrength   string  `json:"trend_strength"`   // 趋势强度：strong/moderate/weak
	VolatilityLevel string  `json:"volatility_level"` // 波动水平：high/medium/low
	DailyChangeRate float64 `json:"daily_change_rate"` // 日均变化率 (%)
	Confidence      float64 `json:"confidence"`       // 预测置信度 (%)
} 