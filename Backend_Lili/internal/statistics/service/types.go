package service

// 通用周期：month/quarter/year，或7d/30d/90d/180d/1y

type DashboardResponse struct {
    DeviceCount       int                    `json:"device_count"`
    TotalValue        float64                `json:"total_value"`
    Categories        map[string]int         `json:"categories"`
    UpcomingReminders int                    `json:"upcoming_reminders"`
    PriceAlerts       int                    `json:"price_alerts"`
}

type DevicesStatisticsRequest struct {
    Period  string `form:"period"`
    GroupBy string `form:"group_by"`
}
type DevicesStatisticsResponse struct {
    Period   string                   `json:"period"`
    GroupBy  string                   `json:"group_by"`
    Series   map[string]int           `json:"series"`
    Trend    []map[string]interface{} `json:"trend"`
}

type ValueAnalysisRequest struct { Period string `form:"period"` }
type ValueAnalysisResponse struct {
    TotalValueTrend   []map[string]interface{} `json:"total_value_trend"`
    CategoryBreakdown map[string]float64       `json:"category_breakdown"`
    AppreciationRate  float64                  `json:"appreciation_rate"`
}

type PriceTrendsRequest struct {
    Period    string `form:"period"`
    DeviceIDs string `form:"device_ids"` // 逗号分隔
}
type PriceTrendsResponse struct {
    Period string                     `json:"period"`
    Items  []map[string]interface{}   `json:"items"`
}

type BrandsStatisticsRequest struct { CategoryID int `form:"category_id"` }
type BrandsStatisticsResponse struct {
    BrandCounts map[string]int     `json:"brand_counts"`
    BrandValues map[string]float64 `json:"brand_values"`
}

type DeviceAgeStatisticsRequest struct { CategoryID int `form:"category_id"` }
type DeviceAgeStatisticsResponse struct {
    Buckets map[string]int `json:"buckets"`
    AvgDays int            `json:"avg_days"`
}

type DepreciationStatisticsRequest struct {
    CategoryID int    `form:"category_id"`
    Period     string `form:"period"`
}
type DepreciationStatisticsResponse struct {
    Distribution map[string]float64 `json:"distribution"`
    Trend        []map[string]interface{} `json:"trend"`
}

type SpendingStatisticsRequest struct {
    Period  string `form:"period"`
    GroupBy string `form:"group_by"`
}
type SpendingStatisticsResponse struct {
    Period string                     `json:"period"`
    GroupBy string                    `json:"group_by"`
    Trend  []map[string]interface{}   `json:"trend"`
}

type HeatmapRequest struct { Type string `form:"type"` }
type HeatmapResponse struct {
    Type string                   `json:"type"`
    Data []map[string]interface{} `json:"data"`
}

type InvestmentReturnRequest struct { IncludeSold *bool `form:"include_sold"` }
type InvestmentReturnResponse struct {
    SoldReturn      float64 `json:"sold_return"`
    HoldingEstimate float64 `json:"holding_estimate"`
    Distribution    []map[string]interface{} `json:"distribution"`
}

type CustomStatisticsRequest struct {
    Metrics []string               `json:"metrics"`
    Filters map[string]interface{} `json:"filters"`
    GroupBy string                 `json:"group_by"`
    Period  string                 `json:"period"`
}
type CustomStatisticsResponse struct { Result map[string]interface{} `json:"result"` }

type ExportReportRequest struct {
    ReportType    string `json:"report_type"`
    Format        string `json:"format"`
    Period        string `json:"period"`
    IncludeCharts *bool  `json:"include_charts"`
}
type ExportReportResponse struct { DownloadURL string `json:"download_url"` }

type InsightsResponse struct { Items []string `json:"items"` }

type ComparisonRequest struct {
    DeviceIDs []int    `json:"device_ids"`
    Metrics   []string `json:"metrics"`
}
type ComparisonResponse struct { Items []map[string]interface{} `json:"items"` }


