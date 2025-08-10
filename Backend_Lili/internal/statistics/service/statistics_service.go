package service

import (
    "strconv"
    "strings"
    "time"

    "Backend_Lili/internal/statistics/repository"
    "Backend_Lili/pkg/utils"
)

type StatisticsService struct { repo *repository.StatisticsRepository }

func NewStatisticsService() *StatisticsService { return &StatisticsService{repo: repository.NewStatisticsRepository()} }

func (s *StatisticsService) GetDashboard(userID int) (*DashboardResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    deviceCount, _ := s.repo.GetDeviceCount(userID)
    totalValue, _ := s.repo.GetDevicesTotalValue(userID)
    categories, _ := s.repo.GetDeviceCountByCategory(userID)
    reminders, _ := s.repo.GetUpcomingRemindersCount(userID)
    alerts, _ := s.repo.GetActivePriceAlertsCount(userID)
    return &DashboardResponse{
        DeviceCount: deviceCount,
        TotalValue: totalValue,
        Categories: categories,
        UpcomingReminders: reminders,
        PriceAlerts: alerts,
    }, nil
}

func (s *StatisticsService) GetDevicesStatistics(userID int, req *DevicesStatisticsRequest) (*DevicesStatisticsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    if req.Period == "" { req.Period = "month" }
    if req.GroupBy == "" { req.GroupBy = "category" }
    series, trend, err := s.repo.GetDevicesStatistics(userID, req.Period, req.GroupBy)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取统计失败") }
    return &DevicesStatisticsResponse{ Period: req.Period, GroupBy: req.GroupBy, Series: series, Trend: trend }, nil
}

func (s *StatisticsService) GetValueAnalysis(userID int, req *ValueAnalysisRequest) (*ValueAnalysisResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    if req.Period == "" { req.Period = "month" }
    trend, _ := s.repo.GetTotalValueTrend(userID, req.Period)
    breakdown, _ := s.repo.GetValueBreakdownByCategory(userID)
    rate, _ := s.repo.GetAppreciationRate(userID, req.Period)
    return &ValueAnalysisResponse{ TotalValueTrend: trend, CategoryBreakdown: breakdown, AppreciationRate: rate }, nil
}

func (s *StatisticsService) GetPriceTrends(userID int, req *PriceTrendsRequest) (*PriceTrendsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    if req.Period == "" { req.Period = "30d" }
    var ids []int
    if strings.TrimSpace(req.DeviceIDs) != "" {
        for _, p := range strings.Split(req.DeviceIDs, ",") {
            if v, err := strconv.Atoi(strings.TrimSpace(p)); err == nil { ids = append(ids, v) }
        }
    }
    items, err := s.repo.GetPriceTrends(userID, req.Period, ids)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取价格趋势失败") }
    return &PriceTrendsResponse{ Period: req.Period, Items: items }, nil
}

func (s *StatisticsService) GetBrandsStatistics(userID int, req *BrandsStatisticsRequest) (*BrandsStatisticsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    counts, values, err := s.repo.GetBrandsStatistics(userID, req.CategoryID)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取品牌统计失败") }
    return &BrandsStatisticsResponse{ BrandCounts: counts, BrandValues: values }, nil
}

func (s *StatisticsService) GetDeviceAgeStatistics(userID int, req *DeviceAgeStatisticsRequest) (*DeviceAgeStatisticsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    buckets, avgDays, err := s.repo.GetDeviceAgeStatistics(userID, req.CategoryID)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取使用年限统计失败") }
    return &DeviceAgeStatisticsResponse{ Buckets: buckets, AvgDays: avgDays }, nil
}

func (s *StatisticsService) GetDepreciationStatistics(userID int, req *DepreciationStatisticsRequest) (*DepreciationStatisticsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    if req.Period == "" { req.Period = "month" }
    dist, trend, err := s.repo.GetDepreciationStatistics(userID, req.CategoryID, req.Period)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取贬值率统计失败") }
    return &DepreciationStatisticsResponse{ Distribution: dist, Trend: trend }, nil
}

func (s *StatisticsService) GetSpendingStatistics(userID int, req *SpendingStatisticsRequest) (*SpendingStatisticsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    if req.Period == "" { req.Period = "year" }
    if req.GroupBy == "" { req.GroupBy = "month" }
    trend, err := s.repo.GetSpendingTrend(userID, req.Period, req.GroupBy)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取支出统计失败") }
    return &SpendingStatisticsResponse{ Period: req.Period, GroupBy: req.GroupBy, Trend: trend }, nil
}

func (s *StatisticsService) GetHeatmap(userID int, req *HeatmapRequest) (*HeatmapResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    if req.Type == "" { return nil, utils.NewBusinessError(utils.ERROR_PARAM, "type 必填") }
    data, err := s.repo.GetHeatmap(userID, req.Type)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取热力图失败") }
    return &HeatmapResponse{ Type: req.Type, Data: data }, nil
}

func (s *StatisticsService) GetInvestmentReturn(userID int, req *InvestmentReturnRequest) (*InvestmentReturnResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    includeSold := true
    if req.IncludeSold != nil { includeSold = *req.IncludeSold }
    sold, holding, dist, err := s.repo.GetInvestmentReturn(userID, includeSold)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取投资收益失败") }
    return &InvestmentReturnResponse{ SoldReturn: sold, HoldingEstimate: holding, Distribution: dist }, nil
}

func (s *StatisticsService) PostCustomStatistics(userID int, req *CustomStatisticsRequest) (*CustomStatisticsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    result, err := s.repo.GetCustomStatistics(userID, req)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取自定义统计失败") }
    return &CustomStatisticsResponse{ Result: result }, nil
}

func (s *StatisticsService) PostExportReport(userID int, req *ExportReportRequest) (*ExportReportResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    if req.ReportType == "" || req.Format == "" { return nil, utils.NewBusinessError(utils.ERROR_PARAM, "report_type/format 必填") }
    // 简化：生成一个临时下载链接
    name := "stats_" + strconv.FormatInt(time.Now().Unix(), 10) + "." + strings.ToLower(req.Format)
    return &ExportReportResponse{ DownloadURL: "/downloads/" + name }, nil
}

func (s *StatisticsService) GetInsights(userID int) (*InsightsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    items := []string{
        "近30天你的设备购买支出较上期下降",
        "手机类设备占比最高，可考虑优化持仓结构",
        "部分设备价格波动较大，建议设置价格预警",
    }
    return &InsightsResponse{ Items: items }, nil
}

func (s *StatisticsService) PostComparison(userID int, req *ComparisonRequest) (*ComparisonResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    items, err := s.repo.GetComparison(userID, req.DeviceIDs, req.Metrics)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取对比分析失败") }
    return &ComparisonResponse{ Items: items }, nil
}


