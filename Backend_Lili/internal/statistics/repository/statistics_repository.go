package repository

import (
    "fmt"
    "time"

    "github.com/beego/beego/v2/client/orm"
)

type StatisticsRepository struct{}

func NewStatisticsRepository() *StatisticsRepository { return &StatisticsRepository{} }

func (r *StatisticsRepository) GetDeviceCount(userID int) (int, error) {
    o := orm.NewOrm()
    cnt, err := o.QueryTable("devices").Filter("user_id", userID).Filter("deleted_at__isnull", true).Count()
    return int(cnt), err
}

func (r *StatisticsRepository) GetDevicesTotalValue(userID int) (float64, error) {
    type Row struct{ Total float64 }
    var row Row
    o := orm.NewOrm()
    err := o.Raw("SELECT COALESCE(SUM(purchase_price),0) AS total FROM devices WHERE user_id = ? AND deleted_at IS NULL", userID).QueryRow(&row)
    return row.Total, err
}

func (r *StatisticsRepository) GetDeviceCountByCategory(userID int) (map[string]int, error) {
    o := orm.NewOrm()
    type Row struct{ Name string; Count int }
    var rows []Row
    _, err := o.Raw(`SELECT COALESCE(c.name,'未分类') as name, COUNT(d.id) as count
        FROM categories c
        LEFT JOIN devices d ON c.id = d.category_id AND d.user_id = ? AND d.deleted_at IS NULL
        GROUP BY c.id, c.name`, userID).QueryRows(&rows)
    if err != nil { return nil, err }
    m := map[string]int{}
    for _, r := range rows { m[r.Name] = r.Count }
    return m, nil
}

func (r *StatisticsRepository) GetUpcomingRemindersCount(userID int) (int, error) {
    // 预留：如有提醒表，可在此统计；当前返回0
    return 0, nil
}

func (r *StatisticsRepository) GetActivePriceAlertsCount(userID int) (int, error) {
    o := orm.NewOrm()
    cnt, err := o.QueryTable("price_alerts").Filter("user_id", userID).Filter("enabled", true).Filter("status", "active").Count()
    return int(cnt), err
}

func (r *StatisticsRepository) GetDevicesStatistics(userID int, period, groupBy string) (map[string]int, []map[string]interface{}, error) {
    // 简化：返回当前分布与空趋势
    cats, err := r.GetDeviceCountByCategory(userID)
    if err != nil { return nil, nil, err }
    return cats, []map[string]interface{}{}, nil
}

func (r *StatisticsRepository) GetTotalValueTrend(userID int, period string) ([]map[string]interface{}, error) {
    // 简化：近6期的静态趋势
    trend := make([]map[string]interface{}, 0)
    for i := 5; i >= 0; i-- {
        trend = append(trend, map[string]interface{}{
            "period":  time.Now().AddDate(0, -i, 0).Format("2006-01"),
            "value":   1000.0 + float64(i*100),
        })
    }
    return trend, nil
}

func (r *StatisticsRepository) GetValueBreakdownByCategory(userID int) (map[string]float64, error) {
    o := orm.NewOrm()
    type Row struct{ Name string; Total float64 }
    var rows []Row
    _, err := o.Raw(`SELECT COALESCE(c.name,'未分类') as name, COALESCE(SUM(d.purchase_price),0) as total
        FROM categories c
        LEFT JOIN devices d ON c.id = d.category_id AND d.user_id = ? AND d.deleted_at IS NULL
        GROUP BY c.id, c.name`, userID).QueryRows(&rows)
    if err != nil { return nil, err }
    m := map[string]float64{}
    for _, r := range rows { m[r.Name] = r.Total }
    return m, nil
}

func (r *StatisticsRepository) GetAppreciationRate(userID int, period string) (float64, error) {
    // 简化：返回固定值
    return 0.05, nil
}

func (r *StatisticsRepository) GetPriceTrends(userID int, period string, deviceIDs []int) ([]map[string]interface{}, error) {
    // 简化：返回空数组，后续可基于 price_histories 汇总
    return []map[string]interface{}{}, nil
}

func (r *StatisticsRepository) GetBrandsStatistics(userID, categoryID int) (map[string]int, map[string]float64, error) {
    o := orm.NewOrm()
    cond := "WHERE d.user_id = ? AND d.deleted_at IS NULL"
    args := []interface{}{userID}
    if categoryID > 0 { cond += " AND d.category_id = ?"; args = append(args, categoryID) }

    type C struct{ Brand string; Cnt int }
    type V struct{ Brand string; Total float64 }
    var cc []C
    var vv []V
    _, err := o.Raw(fmt.Sprintf("SELECT COALESCE(d.brand,'未知') as brand, COUNT(1) as cnt FROM devices d %s GROUP BY d.brand", cond), args...).QueryRows(&cc)
    if err != nil { return nil, nil, err }
    _, err = o.Raw(fmt.Sprintf("SELECT COALESCE(d.brand,'未知') as brand, COALESCE(SUM(d.purchase_price),0) as total FROM devices d %s GROUP BY d.brand", cond), args...).QueryRows(&vv)
    if err != nil { return nil, nil, err }

    counts := map[string]int{}
    values := map[string]float64{}
    for _, r := range cc { counts[r.Brand] = r.Cnt }
    for _, r := range vv { values[r.Brand] = r.Total }
    return counts, values, nil
}

func (r *StatisticsRepository) GetDeviceAgeStatistics(userID, categoryID int) (map[string]int, int, error) {
    o := orm.NewOrm()
    cond := "WHERE user_id = ? AND deleted_at IS NULL AND purchase_date IS NOT NULL"
    args := []interface{}{userID}
    if categoryID > 0 { cond += " AND category_id = ?"; args = append(args, categoryID) }
    type Row struct{ Days int }
    var rows []Row
    _, err := o.Raw("SELECT DATEDIFF(COALESCE(sold_at, NOW()), purchase_date) as days FROM devices "+cond, args...).QueryRows(&rows)
    if err != nil { return nil, 0, err }
    buckets := map[string]int{"<90":0, "90-180":0, "180-365":0, ">365":0}
    sum := 0
    for _, r := range rows { sum += r.Days; switch {
        case r.Days < 90: buckets["<90"]++
        case r.Days < 180: buckets["90-180"]++
        case r.Days < 365: buckets["180-365"]++
        default: buckets[">365"]++
    }}
    avg := 0
    if len(rows) > 0 { avg = sum/len(rows) }
    return buckets, avg, nil
}

func (r *StatisticsRepository) GetDepreciationStatistics(userID, categoryID int, period string) (map[string]float64, []map[string]interface{}, error) {
    // 简化：返回空分布与空趋势
    return map[string]float64{}, []map[string]interface{}{}, nil
}

func (r *StatisticsRepository) GetSpendingTrend(userID int, period, groupBy string) ([]map[string]interface{}, error) {
    // 简化：返回近12期月份的静态趋势
    trend := make([]map[string]interface{}, 0)
    for i := 11; i >= 0; i-- {
        trend = append(trend, map[string]interface{}{
            "period": time.Now().AddDate(0, -i, 0).Format("2006-01"),
            "amount": float64(500 + i*20),
        })
    }
    return trend, nil
}

func (r *StatisticsRepository) GetHeatmap(userID int, heatmapType string) ([]map[string]interface{}, error) {
    // 简化：返回示例点
    return []map[string]interface{}{
        {"x": "手机", "y": "苹果", "value": 12},
        {"x": "相机", "y": "索尼", "value": 8},
    }, nil
}

func (r *StatisticsRepository) GetInvestmentReturn(userID int, includeSold bool) (float64, float64, []map[string]interface{}, error) {
    // 简化：示例数据
    return 1234.56, 789.01, []map[string]interface{}{{"category":"手机","roi":0.12},{"category":"相机","roi":0.08}}, nil
}

func (r *StatisticsRepository) GetCustomStatistics(userID int, req interface{}) (map[string]interface{}, error) {
    // 简化：回显
    return map[string]interface{}{"ok": true}, nil
}

func (r *StatisticsRepository) GetComparison(userID int, deviceIDs []int, metrics []string) ([]map[string]interface{}, error) {
    // 简化：返回设备与指标占位
    items := make([]map[string]interface{}, 0)
    for _, id := range deviceIDs {
        row := map[string]interface{}{"device_id": id}
        for _, m := range metrics { row[m] = 0 }
        items = append(items, row)
    }
    return items, nil
}


