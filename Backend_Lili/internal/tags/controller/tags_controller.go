package controller

import (
    "strconv"

    base "Backend_Lili/internal/auth/controller"
    "Backend_Lili/internal/tags/service"
    "Backend_Lili/pkg/utils"
)

type TagsController struct {
    base.BaseController
    svc *service.TagsService
}

func NewTagsController() *TagsController {
    return &TagsController{svc: service.NewTagsService()}
}

// GET /tags
func (c *TagsController) ListTags() {
    if _, err := c.GetCurrentUser(); err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    var req service.ListTagsRequest
    if err := c.ParseForm(&req); err != nil {
        c.WriteError(utils.ERROR_PARAM, "参数解析失败")
        return
    }
    resp, err := c.svc.ListTags(&req)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// GET /tags/:tagId
func (c *TagsController) GetTag() {
    if _, err := c.GetCurrentUser(); err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    idStr := c.Ctx.Input.Param(":tagId")
    tagID, err := strconv.Atoi(idStr)
    if err != nil {
        c.WriteError(utils.ERROR_PARAM, "标签ID格式错误")
        return
    }
    resp, err := c.svc.GetTag(tagID)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// POST /tags
func (c *TagsController) CreateTag() {
    claims, err := c.GetCurrentUser()
    if err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    var req service.CreateTagRequest
    if err := c.Ctx.Input.Bind(&req, ""); err != nil {
        c.WriteError(utils.ERROR_PARAM, "参数解析失败")
        return
    }
    resp, err := c.svc.CreateTag(claims.UserID, &req)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// PUT /tags/:tagId
func (c *TagsController) UpdateTag() {
    claims, err := c.GetCurrentUser()
    if err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    idStr := c.Ctx.Input.Param(":tagId")
    tagID, err := strconv.Atoi(idStr)
    if err != nil {
        c.WriteError(utils.ERROR_PARAM, "标签ID格式错误")
        return
    }
    var req service.UpdateTagRequest
    if err := c.Ctx.Input.Bind(&req, ""); err != nil {
        c.WriteError(utils.ERROR_PARAM, "参数解析失败")
        return
    }
    resp, err := c.svc.UpdateTag(claims.UserID, tagID, &req)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// DELETE /tags/:tagId
func (c *TagsController) DeleteTag() {
    claims, err := c.GetCurrentUser()
    if err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    idStr := c.Ctx.Input.Param(":tagId")
    tagID, err := strconv.Atoi(idStr)
    if err != nil {
        c.WriteError(utils.ERROR_PARAM, "标签ID格式错误")
        return
    }
    if err := c.svc.DeleteTag(claims.UserID, tagID); err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(map[string]any{"message": "删除成功"})
}

// GET /tags/system
func (c *TagsController) GetSystemTags() {
    if _, err := c.GetCurrentUser(); err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    var req service.ListTagsRequest
    if err := c.ParseForm(&req); err != nil {
        c.WriteError(utils.ERROR_PARAM, "参数解析失败")
        return
    }
    req.Type = "system"
    resp, err := c.svc.ListTags(&req)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// GET /tags/custom
func (c *TagsController) GetCustomTags() {
    claims, err := c.GetCurrentUser()
    if err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    resp, err := c.svc.GetCustomTags(claims.UserID)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// GET /tags/popular
func (c *TagsController) GetPopularTags() {
    if _, err := c.GetCurrentUser(); err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    var req service.GetPopularTagsRequest
    if err := c.ParseForm(&req); err != nil {
        c.WriteError(utils.ERROR_PARAM, "参数解析失败")
        return
    }
    resp, err := c.svc.GetPopularTags(&req)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// GET /tags/search
func (c *TagsController) SearchTags() {
    if _, err := c.GetCurrentUser(); err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    var req service.SearchTagsRequest
    if err := c.ParseForm(&req); err != nil {
        c.WriteError(utils.ERROR_PARAM, "参数解析失败")
        return
    }
    resp, err := c.svc.SearchTags(&req)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// GET /tags/categories
func (c *TagsController) GetTagCategories() {
    if _, err := c.GetCurrentUser(); err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    resp, err := c.svc.GetTagCategories()
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// GET /tags/recommendations
func (c *TagsController) GetRecommendations() {
    claims, err := c.GetCurrentUser()
    if err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    resp, err := c.svc.GetRecommendations(claims.UserID)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}

// GET /tags/statistics
func (c *TagsController) GetTagStatistics() {
    if _, err := c.GetCurrentUser(); err != nil {
        c.WriteError(utils.ERROR_AUTH, "认证失败")
        return
    }
    var req service.GetTagStatisticsRequest
    if err := c.ParseForm(&req); err != nil {
        c.WriteError(utils.ERROR_PARAM, "参数解析失败")
        return
    }
    resp, err := c.svc.GetTagStatistics(&req)
    if err != nil {
        utils.HandleBusinessError(c.Ctx, err)
        return
    }
    c.WriteJSON(resp)
}


