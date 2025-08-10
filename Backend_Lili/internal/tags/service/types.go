package service

import "time"

type ListTagsRequest struct {
    Type     string `form:"type" json:"type"`           // system/custom/all
    Category string `form:"category" json:"category"`
    Active   *bool  `form:"active" json:"active"`
}

type TagInfo struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    Color       string    `json:"color"`
    Icon        string    `json:"icon"`
    Type        string    `json:"type"`
    Active      bool      `json:"active"`
    UsageCount  int       `json:"usage_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ListTagsResponse struct {
    Tags  []*TagInfo `json:"tags"`
    Total int        `json:"total"`
}

type GetTagResponse struct {
    Tag       *TagInfo `json:"tag"`
    UsedBy    int      `json:"used_by"`
}

type CreateTagRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Category    string `json:"category"`
    Color       string `json:"color"`
    Icon        string `json:"icon"`
}

type UpdateTagRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Category    string `json:"category"`
    Color       string `json:"color"`
    Icon        string `json:"icon"`
    Active      *bool  `json:"active"`
}

type GetPopularTagsRequest struct {
    Limit    int    `form:"limit"`
    Category string `form:"category"`
}

type SearchTagsRequest struct {
    Keyword  string `form:"keyword"`
    Category string `form:"category"`
}

type TagCategoryInfo struct {
    Category string `json:"category"`
    Count    int    `json:"count"`
}

type GetTagCategoriesResponse struct {
    Categories []*TagCategoryInfo `json:"categories"`
}

type RecommendationsResponse struct {
    Tags []*TagInfo `json:"tags"`
}

type GetTagStatisticsRequest struct {
    Period string `form:"period"` // month/quarter/year
}

type TagStatisticsResponse struct {
    UsageByTag       map[string]int `json:"usage_by_tag"`
    GrowthByCategory map[string]int `json:"growth_by_category"`
}


