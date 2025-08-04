package service

import (
	"Backend_Lili/internal/device/model"
	"time"
)

// 设备列表请求
type GetDevicesListRequest struct {
	Page       int    `json:"page" form:"page"`               // 页码
	Limit      int    `json:"limit" form:"limit"`             // 每页数量
	CategoryID int    `json:"category_id" form:"category_id"` // 分类ID
	Status     string `json:"status" form:"status"`           // 设备状态
	Sort       string `json:"sort" form:"sort"`               // 排序字段
	Order      string `json:"order" form:"order"`             // 排序方式
	Search     string `json:"search" form:"search"`           // 搜索关键词
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
	TemplateID     int                    `json:"template_id" valid:"Required"`
	Name           string                 `json:"name" valid:"Required"`
	Brand          string                 `json:"brand" valid:"Required"`
	Model          string                 `json:"model" valid:"Required"`
	CategoryID     int                    `json:"category_id" valid:"Required"`
	PurchasePrice  float64                `json:"purchase_price" valid:"Required"`
	PurchaseDate   string                 `json:"purchase_date" valid:"Required"` // YYYY-MM-DD格式
	WarrantyDate   string                 `json:"warranty_date"`                  // YYYY-MM-DD格式
	SerialNumber   string                 `json:"serial_number"`
	Color          string                 `json:"color"`
	Storage        string                 `json:"storage"`
	Memory         string                 `json:"memory"`
	Processor      string                 `json:"processor"`
	ScreenSize     string                 `json:"screen_size"`
	Condition      string                 `json:"condition"` // new/good/fair/poor
	Notes          string                 `json:"notes"`
	Images         []string               `json:"images"`         // 图片URL数组
	Specifications map[string]interface{} `json:"specifications"` // 其他规格参数
}

// 更新设备请求
type UpdateDeviceRequest struct {
	Name           string                 `json:"name"`
	Brand          string                 `json:"brand"`
	Model          string                 `json:"model"`
	SerialNumber   string                 `json:"serial_number"`
	Color          string                 `json:"color"`
	Storage        string                 `json:"storage"`
	Memory         string                 `json:"memory"`
	Processor      string                 `json:"processor"`
	ScreenSize     string                 `json:"screen_size"`
	PurchasePrice  float64                `json:"purchase_price"`
	PurchaseDate   string                 `json:"purchase_date"` // YYYY-MM-DD格式
	WarrantyDate   string                 `json:"warranty_date"` // YYYY-MM-DD格式
	Condition      string                 `json:"condition"`     // new/good/fair/poor
	Notes          string                 `json:"notes"`
	Specifications map[string]interface{} `json:"specifications"` // 其他规格参数
}

// 更新设备状态请求
type UpdateDeviceStatusRequest struct {
	Status    string  `json:"status" valid:"Required"` // active/sold/broken/lost
	SalePrice float64 `json:"sale_price"`              // 出售价格，当status为sold时必填
	SaleDate  string  `json:"sale_date"`               // 出售日期，当status为sold时必填
	Notes     string  `json:"notes"`                   // 状态变更备注
}

// 设备价值评估响应
type DeviceValuationResponse struct {
	DeviceID          int                   `json:"device_id"`
	PurchasePrice     float64               `json:"purchase_price"`
	CurrentValue      float64               `json:"current_value"`
	Depreciation      float64               `json:"depreciation"`       // 贬值金额
	DepreciationRate  float64               `json:"depreciation_rate"`  // 贬值率(%)
	HoldingDays       int                   `json:"holding_days"`       // 持有天数
	DailyDepreciation float64               `json:"daily_depreciation"` // 日均贬值
	LastUpdateTime    time.Time             `json:"last_update_time"`
	PriceHistories    []*model.PriceHistory `json:"price_histories"`
}

// 批量导入设备请求
type BatchImportDevicesRequest struct {
	Devices          []CreateDeviceRequest `json:"devices" valid:"Required"`
	IgnoreDuplicates bool                  `json:"ignore_duplicates"` // 是否忽略重复设备
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
	Algorithm        string             `json:"algorithm"`         // 使用的算法名称
	Accuracy         float64            `json:"accuracy"`          // 预测准确度评估 (0-1)
	PredictionPoints []*PredictionPoint `json:"prediction_points"` // 预测数据点
	TrendAnalysis    *TrendAnalysis     `json:"trend_analysis"`    // 趋势分析
	CreatedAt        time.Time          `json:"created_at"`
}

// 趋势分析
type TrendAnalysis struct {
	Trend           string  `json:"trend"`             // 趋势方向：rising/falling/stable
	TrendStrength   string  `json:"trend_strength"`    // 趋势强度：strong/moderate/weak
	VolatilityLevel string  `json:"volatility_level"`  // 波动水平：high/medium/low
	DailyChangeRate float64 `json:"daily_change_rate"` // 日均变化率 (%)
	Confidence      float64 `json:"confidence"`        // 预测置信度 (%)
}

// =============模板相关类型=============

// 获取模板列表请求
type GetTemplatesListRequest struct {
	CategoryID int  `json:"category_id" form:"category_id"` // 分类ID
	Active     bool `json:"active" form:"active"`           // 是否只显示启用的模板
	Page       int  `json:"page" form:"page"`               // 页码
	Limit      int  `json:"limit" form:"limit"`             // 每页数量
}

// 获取模板列表响应
type GetTemplatesListResponse struct {
	Templates  []*model.DeviceTemplate `json:"templates"`
	Total      int                     `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}

// 创建模板请求
type CreateTemplateRequest struct {
	Name        string                 `json:"name" valid:"Required"`
	CategoryID  int                    `json:"category_id" valid:"Required"`
	Description string                 `json:"description"`
	Icon        string                 `json:"icon"`
	Fields      []*model.TemplateField `json:"fields" valid:"Required"`
	Active      bool                   `json:"active"`
}

// 更新模板请求
type UpdateTemplateRequest struct {
	Name        string                 `json:"name"`
	CategoryID  int                    `json:"category_id"`
	Description string                 `json:"description"`
	Icon        string                 `json:"icon"`
	Fields      []*model.TemplateField `json:"fields"`
	Active      bool                   `json:"active"`
}

// 验证设备数据请求
type ValidateDeviceDataRequest struct {
	DeviceData map[string]interface{} `json:"device_data" valid:"Required"`
}

// 验证设备数据响应
type ValidateDeviceDataResponse struct {
	Valid   bool                   `json:"valid"`
	Errors  map[string]string      `json:"errors,omitempty"`
	Missing []string               `json:"missing,omitempty"` // 缺失的必填字段
	Invalid map[string]interface{} `json:"invalid,omitempty"` // 无效的字段值
}

// 模板统计响应
type TemplateStatisticsResponse struct {
	TemplateID    int                 `json:"template_id"`
	UseCount      int                 `json:"use_count"`      // 使用次数
	DeviceCount   int                 `json:"device_count"`   // 基于该模板创建的设备数量
	UserCount     int                 `json:"user_count"`     // 使用该模板的用户数量
	RecentDevices []*model.Device     `json:"recent_devices"` // 最近创建的设备
	PopularFields []*PopularFieldInfo `json:"popular_fields"` // 热门字段统计
}

// 热门字段信息
type PopularFieldInfo struct {
	FieldName string  `json:"field_name"`
	UseCount  int     `json:"use_count"`
	UseRate   float64 `json:"use_rate"` // 使用率
}

// =============分类相关类型=============

// 获取分类列表请求
type GetCategoriesListRequest struct {
	Type         string `json:"type" form:"type"`                   // system/custom/all
	ParentID     int    `json:"parent_id" form:"parent_id"`         // 父分类ID
	IncludeCount bool   `json:"include_count" form:"include_count"` // 是否包含设备数量
}

// 获取分类列表响应
type GetCategoriesListResponse struct {
	Categories []*model.Category `json:"categories"`
}

// 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" valid:"Required"`
	Description string `json:"description"`
	ParentID    int    `json:"parent_id"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
}

// 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    int    `json:"parent_id"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
}

// 分类排序请求
type SortCategoriesRequest struct {
	CategoryOrders []*CategoryOrderItem `json:"category_orders" valid:"Required"`
}

// 分类排序项
type CategoryOrderItem struct {
	CategoryID int `json:"category_id" valid:"Required"`
	SortOrder  int `json:"sort_order" valid:"Required"`
}

// 分类统计请求
type GetCategoryStatisticsRequest struct {
	Period string `json:"period" form:"period"` // month/quarter/year
}

// 分类统计响应
type CategoryStatisticsResponse struct {
	Period            string                       `json:"period"`
	Categories        []*CategoryStatInfo          `json:"categories"`
	TotalDevices      int                          `json:"total_devices"`
	TotalValue        float64                      `json:"total_value"`
	TrendData         []*CategoryTrendData         `json:"trend_data"`         // 趋势数据
	ValueDistribution []*CategoryValueDistribution `json:"value_distribution"` // 价值分布
}

// 分类统计信息
type CategoryStatInfo struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	DeviceCount  int     `json:"device_count"`
	TotalValue   float64 `json:"total_value"`
	AvgValue     float64 `json:"avg_value"`
	Percentage   float64 `json:"percentage"` // 占比
}

// 分类趋势数据
type CategoryTrendData struct {
	Date        string  `json:"date"`
	CategoryID  int     `json:"category_id"`
	DeviceCount int     `json:"device_count"`
	TotalValue  float64 `json:"total_value"`
}

// 分类价值分布
type CategoryValueDistribution struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	ValueRange   string  `json:"value_range"` // 价值范围：0-1000, 1000-5000等
	DeviceCount  int     `json:"device_count"`
	Percentage   float64 `json:"percentage"`
}

// 搜索分类请求
type SearchCategoriesRequest struct {
	Keyword string `json:"keyword" form:"keyword" valid:"Required"`
	Type    string `json:"type" form:"type"` // system/custom/all
}
