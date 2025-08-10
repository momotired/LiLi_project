package model

import (
	"time"
)

// Price 设备当前价格信息表
type Price struct {
	ID           int       `orm:"column(id);auto;pk" json:"id"`
	DeviceID     int       `orm:"column(device_id)" json:"device_id"`
	UserID       int       `orm:"column(user_id)" json:"user_id"`
	CurrentPrice float64   `orm:"column(current_price);digits(10);decimals(2)" json:"current_price"`
	MarketPrice  float64   `orm:"column(market_price);digits(10);decimals(2);null" json:"market_price"`
	AveragePrice float64   `orm:"column(average_price);digits(10);decimals(2);null" json:"average_price"`
	MinPrice     float64   `orm:"column(min_price);digits(10);decimals(2);null" json:"min_price"`
	MaxPrice     float64   `orm:"column(max_price);digits(10);decimals(2);null" json:"max_price"`
	PriceChange  float64   `orm:"column(price_change);digits(10);decimals(2);default(0)" json:"price_change"` // 价格变化金额
	ChangeRate   float64   `orm:"column(change_rate);digits(5);decimals(2);default(0)" json:"change_rate"`    // 变化百分比
	TrendStatus  string    `orm:"column(trend_status);size(20);default(stable)" json:"trend_status"`          // rising/falling/stable
	LastUpdateAt time.Time `orm:"column(last_update_at);type(datetime);null" json:"last_update_at"`           // 最后更新时间
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (p *Price) TableName() string {
	return "prices"
}

// PriceHistory 价格历史记录表（扩展版）
type PriceHistory struct {
	ID          int       `orm:"column(id);auto;pk" json:"id"`
	DeviceID    int       `orm:"column(device_id)" json:"device_id"`
	UserID      int       `orm:"column(user_id)" json:"user_id"`
	Source      string    `orm:"column(source);size(50)" json:"source"`             // 价格来源：manual/api/web_scrape
	SourceID    string    `orm:"column(source_id);size(100);null" json:"source_id"` // 来源平台的商品ID
	Platform    string    `orm:"column(platform);size(50);null" json:"platform"`    // 平台名称：如闲鱼、转转、京东等
	Price       float64   `orm:"column(price);digits(10);decimals(2)" json:"price"`
	Condition   string    `orm:"column(condition);size(20)" json:"condition"` // new/good/fair/poor
	Description string    `orm:"column(description);size(500);null" json:"description"`
	URL         string    `orm:"column(url);size(1000);null" json:"url"` // 商品链接
	RecordDate  time.Time `orm:"column(record_date);type(date)" json:"record_date"`
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
}

func (ph *PriceHistory) TableName() string {
	return "price_histories"
}

// PriceAlert 价格预警表
type PriceAlert struct {
	ID                  int       `orm:"column(id);auto;pk" json:"id"`
	DeviceID            int       `orm:"column(device_id)" json:"device_id"`
	UserID              int       `orm:"column(user_id)" json:"user_id"`
	AlertType           string    `orm:"column(alert_type);size(20)" json:"alert_type"`                           // price_drop/price_rise/target_price
	Threshold           float64   `orm:"column(threshold);digits(10);decimals(2)" json:"threshold"`               // 预警阈值
	ThresholdType       string    `orm:"column(threshold_type);size(20)" json:"threshold_type"`                   // percentage/absolute
	Enabled             bool      `orm:"column(enabled);default(true)" json:"enabled"`                            // 是否启用
	NotificationMethods string    `orm:"column(notification_methods);size(200);null" json:"notification_methods"` // JSON格式存储通知方式
	LastTriggeredAt     time.Time `orm:"column(last_triggered_at);type(datetime);null" json:"last_triggered_at"`
	TriggerCount        int       `orm:"column(trigger_count);default(0)" json:"trigger_count"` // 触发次数
	Status              string    `orm:"column(status);size(20);default(active)" json:"status"` // active/triggered/disabled/expired
	CreatedAt           time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt           time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (pa *PriceAlert) TableName() string {
	return "price_alerts"
}

// PriceSource 价格数据源表
type PriceSource struct {
	ID          int       `orm:"column(id);auto;pk" json:"id"`
	Name        string    `orm:"column(name);size(100)" json:"name"`              // 数据源名称
	Platform    string    `orm:"column(platform);size(50)" json:"platform"`       // 平台标识
	BaseURL     string    `orm:"column(base_url);size(500);null" json:"base_url"` // 基础URL
	ApiEndpoint string    `orm:"column(api_endpoint);size(500);null" json:"api_endpoint"`
	Status      string    `orm:"column(status);size(20);default(active)" json:"status"`                      // active/inactive/error
	Reliability float64   `orm:"column(reliability);digits(3);decimals(2);default(1.00)" json:"reliability"` // 可靠性评分 0-1
	UpdateFreq  int       `orm:"column(update_freq);default(24)" json:"update_freq"`                         // 更新频率（小时）
	LastSync    time.Time `orm:"column(last_sync);type(datetime);null" json:"last_sync"`
	Config      string    `orm:"column(config);type(json);null" json:"config"` // JSON格式配置信息
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (ps *PriceSource) TableName() string {
	return "price_sources"
}

// PricePrediction 价格预测表
type PricePrediction struct {
	ID             int       `orm:"column(id);auto;pk" json:"id"`
	DeviceID       int       `orm:"column(device_id)" json:"device_id"`
	UserID         int       `orm:"column(user_id)" json:"user_id"`
	PredictionType string    `orm:"column(prediction_type);size(20)" json:"prediction_type"` // 30d/90d/180d
	PredictedPrice float64   `orm:"column(predicted_price);digits(10);decimals(2)" json:"predicted_price"`
	Confidence     float64   `orm:"column(confidence);digits(3);decimals(2)" json:"confidence"`           // 置信度 0-1
	Algorithm      string    `orm:"column(algorithm);size(50)" json:"algorithm"`                          // 使用的算法
	Factors        string    `orm:"column(factors);type(json);null" json:"factors"`                       // 影响因素JSON
	ValidUntil     time.Time `orm:"column(valid_until);type(datetime)" json:"valid_until"`                // 预测有效期
	ActualPrice    float64   `orm:"column(actual_price);digits(10);decimals(2);null" json:"actual_price"` // 实际价格（用于验证）
	Accuracy       float64   `orm:"column(accuracy);digits(3);decimals(2);null" json:"accuracy"`          // 准确度
	CreatedAt      time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
}

func (pp *PricePrediction) TableName() string {
	return "price_predictions"
}

// NotificationMethod 通知方式结构
type NotificationMethod struct {
	Type    string `json:"type"`    // push/email/sms
	Address string `json:"address"` // 通知地址
	Enabled bool   `json:"enabled"` // 是否启用
}

// PredictionFactor 价格预测影响因素
type PredictionFactor struct {
	Factor string  `json:"factor"` // 因素名称
	Weight float64 `json:"weight"` // 权重
	Value  float64 `json:"value"`  // 当前值
	Impact string  `json:"impact"` // 影响方向: positive/negative/neutral
}
