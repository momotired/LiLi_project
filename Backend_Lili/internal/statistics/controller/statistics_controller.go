package controller

import (
    "encoding/json"

    base "Backend_Lili/internal/auth/controller"
    "Backend_Lili/internal/statistics/service"
    "Backend_Lili/pkg/utils"
)

type StatisticsController struct {
    base.BaseController
    svc *service.StatisticsService
}

func NewStatisticsController() *StatisticsController {
    return &StatisticsController{}
}

func (c *StatisticsController) Prepare() {
    // 调用父类的Prepare方法
    c.BaseController.Prepare()
    
    // 初始化统计服务
    c.svc = service.NewStatisticsService()
}

// GET /statistics/dashboard
func (c *StatisticsController) GetDashboard() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    resp, err := c.svc.GetDashboard(claims.UserID)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/devices
func (c *StatisticsController) GetDevicesStatistics() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.DevicesStatisticsRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetDevicesStatistics(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/value-analysis
func (c *StatisticsController) GetValueAnalysis() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.ValueAnalysisRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetValueAnalysis(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/price-trends
func (c *StatisticsController) GetPriceTrends() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.PriceTrendsRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetPriceTrends(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/brands
func (c *StatisticsController) GetBrandsStatistics() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.BrandsStatisticsRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetBrandsStatistics(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/device-age
func (c *StatisticsController) GetDeviceAgeStatistics() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.DeviceAgeStatisticsRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetDeviceAgeStatistics(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/depreciation
func (c *StatisticsController) GetDepreciationStatistics() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.DepreciationStatisticsRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetDepreciationStatistics(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/spending
func (c *StatisticsController) GetSpendingStatistics() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.SpendingStatisticsRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetSpendingStatistics(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/heatmap
func (c *StatisticsController) GetHeatmap() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.HeatmapRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetHeatmap(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/investment-return
func (c *StatisticsController) GetInvestmentReturn() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.InvestmentReturnRequest
    if err := c.ParseForm(&req); err != nil { c.WriteError(utils.ERROR_PARAM, "参数解析失败"); return }
    resp, err := c.svc.GetInvestmentReturn(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// POST /statistics/custom
func (c *StatisticsController) PostCustomStatistics() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.CustomStatisticsRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil { c.WriteError(utils.ERROR_PARAM, "请求参数格式错误"); return }
    resp, err := c.svc.PostCustomStatistics(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// POST /statistics/export
func (c *StatisticsController) PostExportReport() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.ExportReportRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil { c.WriteError(utils.ERROR_PARAM, "请求参数格式错误"); return }
    resp, err := c.svc.PostExportReport(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// GET /statistics/insights
func (c *StatisticsController) GetInsights() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    resp, err := c.svc.GetInsights(claims.UserID)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}

// POST /statistics/comparison
func (c *StatisticsController) PostComparison() {
    claims, err := c.GetCurrentUser()
    if err != nil { c.WriteError(utils.ERROR_AUTH, "认证失败"); return }
    var req service.ComparisonRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil { c.WriteError(utils.ERROR_PARAM, "请求参数格式错误"); return }
    resp, err := c.svc.PostComparison(claims.UserID, &req)
    if err != nil { utils.HandleBusinessError(c.Ctx, err); return }
    c.WriteJSON(resp)
}


