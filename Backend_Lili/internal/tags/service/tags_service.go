package service

import (
    "Backend_Lili/internal/tags/repository"
    umodel "Backend_Lili/internal/user/model"
    "Backend_Lili/pkg/utils"
)

type TagsService struct {
    repo *repository.TagsRepository
}

func NewTagsService() *TagsService { return &TagsService{repo: repository.NewTagsRepository()} }

func (s *TagsService) toInfo(t *umodel.Tag) *TagInfo {
    return &TagInfo{
        ID:          t.ID,
        Name:        t.Name,
        Description: t.Description,
        Category:    t.Category,
        Color:       t.Color,
        Icon:        t.Icon,
        Type:        t.Type,
        Active:      t.Active,
        UsageCount:  t.UsageCount,
        CreatedAt:   t.CreatedAt,
        UpdatedAt:   t.UpdatedAt,
    }
}

func (s *TagsService) ListTags(req *ListTagsRequest) (*ListTagsResponse, error) {
    params := map[string]interface{}{}
    if req != nil {
        if req.Type != "" { params["type"] = req.Type }
        if req.Category != "" { params["category"] = req.Category }
        if req.Active != nil { params["active"] = *req.Active }
    }
    tags, err := s.repo.ListTags(params)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取标签失败") }
    resp := &ListTagsResponse{Tags: make([]*TagInfo, 0, len(tags)), Total: len(tags)}
    for _, t := range tags { resp.Tags = append(resp.Tags, s.toInfo(t)) }
    return resp, nil
}

func (s *TagsService) GetTag(tagID int) (*GetTagResponse, error) {
    if tagID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效") }
    tag, err := s.repo.GetTagByID(tagID)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取标签失败") }
    if tag == nil { return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "标签不存在") }
    usedBy, _ := s.repo.CountUserTagUsage(tagID)
    return &GetTagResponse{Tag: s.toInfo(tag), UsedBy: usedBy}, nil
}

func (s *TagsService) CreateTag(userID int, req *CreateTagRequest) (*TagInfo, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    if req.Name == "" || req.Category == "" { return nil, utils.NewBusinessError(utils.ERROR_PARAM, "名称和分类必填") }
    exists, err := s.repo.ExistsTagNameForUser(req.Name, userID)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "校验失败") }
    if exists { return nil, utils.NewBusinessError(utils.ERROR_BUSINESS, "标签名称已存在") }
    tag := &umodel.Tag{
        Name:        req.Name,
        Description: req.Description,
        Category:    req.Category,
        Color:       req.Color,
        Icon:        req.Icon,
        Type:        "custom",
        Active:      true,
        OwnerID:     userID,
    }
    if err := s.repo.CreateTag(tag); err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "创建失败") }
    return s.toInfo(tag), nil
}

func (s *TagsService) UpdateTag(userID, tagID int, req *UpdateTagRequest) (*TagInfo, error) {
    if userID <= 0 || tagID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_PARAM, "参数无效") }
    tag, err := s.repo.GetTagByID(tagID)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取失败") }
    if tag == nil { return nil, utils.NewBusinessError(utils.ERROR_NOT_FOUND, "标签不存在") }
    if tag.Type == "system" { return nil, utils.NewBusinessError(utils.ERROR_FORBIDDEN, "系统标签不可修改") }
    if tag.OwnerID != userID { return nil, utils.NewBusinessError(utils.ERROR_FORBIDDEN, "无权操作该标签") }
    if req.Name != "" { tag.Name = req.Name }
    if req.Description != "" { tag.Description = req.Description }
    if req.Category != "" { tag.Category = req.Category }
    if req.Color != "" { tag.Color = req.Color }
    if req.Icon != "" { tag.Icon = req.Icon }
    if req.Active != nil { tag.Active = *req.Active }
    if err := s.repo.UpdateTag(tag); err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "更新失败") }
    return s.toInfo(tag), nil
}

func (s *TagsService) DeleteTag(userID, tagID int) error {
    if userID <= 0 || tagID <= 0 { return utils.NewBusinessError(utils.ERROR_PARAM, "参数无效") }
    tag, err := s.repo.GetTagByID(tagID)
    if err != nil { return utils.NewBusinessError(utils.ERROR_DATABASE, "获取失败") }
    if tag == nil { return utils.NewBusinessError(utils.ERROR_NOT_FOUND, "标签不存在") }
    if tag.Type == "system" { return utils.NewBusinessError(utils.ERROR_FORBIDDEN, "系统标签不可删除") }
    if tag.OwnerID != userID { return utils.NewBusinessError(utils.ERROR_FORBIDDEN, "无权操作该标签") }
    inUse, err := s.repo.IsTagInUse(tagID)
    if err != nil { return utils.NewBusinessError(utils.ERROR_DATABASE, "校验失败") }
    if inUse { return utils.NewBusinessError(utils.ERROR_BUSINESS, "标签已被使用，无法删除") }
    if err := s.repo.DeleteTag(tagID); err != nil { return utils.NewBusinessError(utils.ERROR_DATABASE, "删除失败") }
    return nil
}

func (s *TagsService) GetCustomTags(userID int) (*ListTagsResponse, error) {
    if userID <= 0 { return nil, utils.NewBusinessError(utils.ERROR_AUTH, "认证失败") }
    tags, err := s.repo.GetCustomTags(userID)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取失败") }
    resp := &ListTagsResponse{Tags: make([]*TagInfo, 0, len(tags)), Total: len(tags)}
    for _, t := range tags { resp.Tags = append(resp.Tags, s.toInfo(t)) }
    return resp, nil
}

func (s *TagsService) GetPopularTags(req *GetPopularTagsRequest) (*ListTagsResponse, error) {
    limit := 10
    if req != nil && req.Limit > 0 { limit = req.Limit }
    category := ""
    if req != nil { category = req.Category }
    tags, err := s.repo.GetPopularTags(limit, category)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取失败") }
    resp := &ListTagsResponse{Tags: make([]*TagInfo, 0, len(tags)), Total: len(tags)}
    for _, t := range tags { resp.Tags = append(resp.Tags, s.toInfo(t)) }
    return resp, nil
}

func (s *TagsService) SearchTags(req *SearchTagsRequest) (*ListTagsResponse, error) {
    if req == nil || req.Keyword == "" { return &ListTagsResponse{Tags: []*TagInfo{}, Total: 0}, nil }
    tags, err := s.repo.SearchTags(req.Keyword, req.Category)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "搜索失败") }
    resp := &ListTagsResponse{Tags: make([]*TagInfo, 0, len(tags)), Total: len(tags)}
    for _, t := range tags { resp.Tags = append(resp.Tags, s.toInfo(t)) }
    return resp, nil
}

func (s *TagsService) GetTagCategories() (*GetTagCategoriesResponse, error) {
    rows, err := s.repo.GetTagCategories()
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取失败") }
    resp := &GetTagCategoriesResponse{Categories: make([]*TagCategoryInfo, 0, len(rows))}
    for _, r := range rows { resp.Categories = append(resp.Categories, &TagCategoryInfo{Category: r.Category, Count: r.Count}) }
    return resp, nil
}

func (s *TagsService) GetRecommendations(userID int) (*RecommendationsResponse, error) {
    // 简化：推荐热门系统标签
    tags, err := s.repo.GetPopularTags(10, "")
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取失败") }
    resp := &RecommendationsResponse{Tags: make([]*TagInfo, 0, len(tags))}
    for _, t := range tags { resp.Tags = append(resp.Tags, s.toInfo(t)) }
    return resp, nil
}

func (s *TagsService) GetTagStatistics(req *GetTagStatisticsRequest) (*TagStatisticsResponse, error) {
    // 简化：统计按标签使用次数与分类增长（此处返回静态或基于当前总数）
    params := map[string]interface{}{}
    tags, err := s.repo.ListTags(params)
    if err != nil { return nil, utils.NewBusinessError(utils.ERROR_DATABASE, "获取失败") }
    usage := map[string]int{}
    growth := map[string]int{}
    for _, t := range tags {
        usage[t.Name] = t.UsageCount
        growth[t.Category] += 1
    }
    return &TagStatisticsResponse{UsageByTag: usage, GrowthByCategory: growth}, nil
}


