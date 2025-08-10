package repository

import (
	"Backend_Lili/internal/price/model"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type PriceRepository struct{}

func NewPriceRepository() *PriceRepository {
	return &PriceRepository{}
}

// GetDevicePrice 获取设备当前价格信息
func (r *PriceRepository) GetDevicePrice(deviceID, userID int) (*model.Price, error) {
	o := orm.NewOrm()
	price := &model.Price{}
	err := o.QueryTable("prices").
		Filter("device_id", deviceID).
		Filter("user_id", userID).
		One(price)

	if err == orm.ErrNoRows {
		return nil, nil
	}
	return price, err
}

// UpdateDevicePrice 更新或创建设备价格信息
func (r *PriceRepository) UpdateDevicePrice(price *model.Price) error {
	o := orm.NewOrm()

	// 检查是否已存在
	existing := &model.Price{}
	err := o.QueryTable("prices").
		Filter("device_id", price.DeviceID).
		Filter("user_id", price.UserID).
		One(existing)

	if err == orm.ErrNoRows {
		// 不存在，创建新记录
		_, err = o.Insert(price)
		return err
	} else if err != nil {
		return err
	}

	// 存在，更新记录
	price.ID = existing.ID
	price.CreatedAt = existing.CreatedAt
	_, err = o.Update(price)
	return err
}

// GetPriceHistory 获取价格历史记录
func (r *PriceRepository) GetPriceHistory(deviceID, userID int, params map[string]interface{}) ([]*model.PriceHistory, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("price_histories").
		Filter("device_id", deviceID).
		Filter("user_id", userID)

	// 时间范围筛选
	if period, ok := params["period"].(string); ok && period != "" {
		var startDate time.Time
		now := time.Now()
		switch period {
		case "7d":
			startDate = now.AddDate(0, 0, -7)
		case "30d":
			startDate = now.AddDate(0, -1, 0)
		case "90d":
			startDate = now.AddDate(0, -3, 0)
		case "180d":
			startDate = now.AddDate(0, -6, 0)
		case "1y":
			startDate = now.AddDate(-1, 0, 0)
		}
		if !startDate.IsZero() {
			qs = qs.Filter("record_date__gte", startDate.Format("2006-01-02"))
		}
	}

	// 数据源筛选
	if source, ok := params["source"].(string); ok && source != "" {
		qs = qs.Filter("source", source)
	}

	qs = qs.OrderBy("-record_date", "-created_at")

	var histories []*model.PriceHistory
	_, err := qs.All(&histories)
	return histories, err
}

// CreatePriceHistory 创建价格历史记录
func (r *PriceRepository) CreatePriceHistory(history *model.PriceHistory) error {
	o := orm.NewOrm()
	_, err := o.Insert(history)
	return err
}

// GetPriceAlerts 获取价格预警列表
func (r *PriceRepository) GetPriceAlerts(userID int, params map[string]interface{}) ([]*model.PriceAlert, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("price_alerts").Filter("user_id", userID)

	// 状态筛选
	if status, ok := params["status"].(string); ok && status != "" {
		qs = qs.Filter("status", status)
	}

	// 设备筛选
	if deviceID, ok := params["device_id"].(int); ok && deviceID > 0 {
		qs = qs.Filter("device_id", deviceID)
	}

	qs = qs.OrderBy("-created_at")

	var alerts []*model.PriceAlert
	_, err := qs.All(&alerts)
	return alerts, err
}

// CreatePriceAlert 创建价格预警
func (r *PriceRepository) CreatePriceAlert(alert *model.PriceAlert) error {
	o := orm.NewOrm()
	_, err := o.Insert(alert)
	return err
}

// UpdatePriceAlert 更新价格预警
func (r *PriceRepository) UpdatePriceAlert(alert *model.PriceAlert) error {
	o := orm.NewOrm()
	_, err := o.Update(alert)
	return err
}

// DeletePriceAlert 删除价格预警
func (r *PriceRepository) DeletePriceAlert(alertID, userID int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("price_alerts").
		Filter("id", alertID).
		Filter("user_id", userID).
		Delete()
	return err
}

// GetPriceAlertByID 根据ID获取价格预警
func (r *PriceRepository) GetPriceAlertByID(alertID, userID int) (*model.PriceAlert, error) {
	o := orm.NewOrm()
	alert := &model.PriceAlert{}
	err := o.QueryTable("price_alerts").
		Filter("id", alertID).
		Filter("user_id", userID).
		One(alert)

	if err == orm.ErrNoRows {
		return nil, nil
	}
	return alert, err
}

// GetPriceSources 获取价格数据源列表
func (r *PriceRepository) GetPriceSources() ([]*model.PriceSource, error) {
	o := orm.NewOrm()
	var sources []*model.PriceSource
	_, err := o.QueryTable("price_sources").
		Filter("status", "active").
		OrderBy("-reliability", "name").
		All(&sources)
	return sources, err
}

// BatchUpdatePrices 批量更新价格
func (r *PriceRepository) BatchUpdatePrices(prices []*model.Price) error {
	o := orm.NewOrm()
	// 开启事务
	err := o.Begin()
	if err != nil {
		return err
	}

	for _, price := range prices {
		// 检查是否已存在
		existing := &model.Price{}
		err := o.QueryTable("prices").
			Filter("device_id", price.DeviceID).
			Filter("user_id", price.UserID).
			One(existing)

		if err == orm.ErrNoRows {
			// 不存在，创建新记录
			_, err = o.Insert(price)
			if err != nil {
				o.Rollback()
				return err
			}
		} else if err != nil {
			o.Rollback()
			return err
		} else {
			// 存在，更新记录
			price.ID = existing.ID
			price.CreatedAt = existing.CreatedAt
			_, err = o.Update(price)
			if err != nil {
				o.Rollback()
				return err
			}
		}
	}

	return o.Commit()
}

// GetPricePrediction 获取价格预测
func (r *PriceRepository) GetPricePrediction(deviceID, userID int, predictionType string) (*model.PricePrediction, error) {
	o := orm.NewOrm()
	prediction := &model.PricePrediction{}
	err := o.QueryTable("price_predictions").
		Filter("device_id", deviceID).
		Filter("user_id", userID).
		Filter("prediction_type", predictionType).
		Filter("valid_until__gt", time.Now()).
		OrderBy("-created_at").
		One(prediction)

	if err == orm.ErrNoRows {
		return nil, nil
	}
	return prediction, err
}

// CreatePricePrediction 创建价格预测
func (r *PriceRepository) CreatePricePrediction(prediction *model.PricePrediction) error {
	o := orm.NewOrm()
	_, err := o.Insert(prediction)
	return err
}

// GetMarketComparison 获取市场价格对比数据
func (r *PriceRepository) GetMarketComparison(deviceID, userID int) ([]*model.PriceHistory, error) {
	o := orm.NewOrm()

	// 获取最近30天各平台的最新价格
	var histories []*model.PriceHistory

	sql := `
		SELECT ph1.* FROM price_histories ph1
		INNER JOIN (
			SELECT platform, MAX(record_date) as max_date
			FROM price_histories 
			WHERE device_id = ? AND user_id = ? 
			AND record_date >= DATE_SUB(NOW(), INTERVAL 30 DAY)
			AND platform IS NOT NULL
			GROUP BY platform
		) ph2 ON ph1.platform = ph2.platform AND ph1.record_date = ph2.max_date
		WHERE ph1.device_id = ? AND ph1.user_id = ?
		ORDER BY ph1.price ASC
	`

	_, err := o.Raw(sql, deviceID, userID, deviceID, userID).QueryRows(&histories)
	return histories, err
}

// GetPriceStatistics 获取价格统计信息
func (r *PriceRepository) GetPriceStatistics(deviceID, userID int, days int) (map[string]interface{}, error) {
	o := orm.NewOrm()

	startDate := time.Now().AddDate(0, 0, -days)

	var result []orm.Params
	sql := `
		SELECT 
			COUNT(*) as record_count,
			MIN(price) as min_price,
			MAX(price) as max_price,
			AVG(price) as avg_price,
			STDDEV(price) as price_stddev
		FROM price_histories 
		WHERE device_id = ? AND user_id = ? AND record_date >= ?
	`

	_, err := o.Raw(sql, deviceID, userID, startDate.Format("2006-01-02")).Values(&result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return map[string]interface{}{
			"record_count": 0,
			"min_price":    0,
			"max_price":    0,
			"avg_price":    0,
			"price_stddev": 0,
		}, nil
	}

	return map[string]interface{}(result[0]), nil
}

// CheckPriceAlerts 检查需要触发的价格预警
func (r *PriceRepository) CheckPriceAlerts(deviceID int, currentPrice float64) ([]*model.PriceAlert, error) {
	o := orm.NewOrm()

	var alerts []*model.PriceAlert

	// 获取该设备的所有启用的预警
	_, err := o.QueryTable("price_alerts").
		Filter("device_id", deviceID).
		Filter("enabled", true).
		Filter("status", "active").
		All(&alerts)

	if err != nil {
		return nil, err
	}

	var triggeredAlerts []*model.PriceAlert

	// 检查每个预警是否需要触发
	for _, alert := range alerts {
		shouldTrigger := false

		switch alert.AlertType {
		case "price_drop":
			if alert.ThresholdType == "percentage" {
				// 假设有上次价格进行比较
				// 这里简化处理，实际应该获取历史价格进行比较
				shouldTrigger = true // 需要实际逻辑
			} else {
				shouldTrigger = currentPrice <= alert.Threshold
			}
		case "price_rise":
			if alert.ThresholdType == "percentage" {
				shouldTrigger = true // 需要实际逻辑
			} else {
				shouldTrigger = currentPrice >= alert.Threshold
			}
		case "target_price":
			shouldTrigger = currentPrice <= alert.Threshold
		}

		if shouldTrigger {
			triggeredAlerts = append(triggeredAlerts, alert)
		}
	}

	return triggeredAlerts, nil
}
