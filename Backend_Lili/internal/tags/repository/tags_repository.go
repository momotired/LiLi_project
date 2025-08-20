package repository

import (
    "Backend_Lili/internal/user/model"
    "strings"

    "github.com/beego/beego/v2/client/orm"
)

type TagsRepository struct{}

func NewTagsRepository() *TagsRepository { return &TagsRepository{} }

func (r *TagsRepository) ListTags(params map[string]interface{}) ([]*model.Tag, error) {
    o := orm.NewOrm()
    qs := o.QueryTable("tags")
    if t, ok := params["type"].(string); ok && t != "" && t != "all" {
        qs = qs.Filter("type", t)
    }
    if c, ok := params["category"].(string); ok && c != "" {
        qs = qs.Filter("category", c)
    }
    if active, ok := params["active"].(bool); ok {
        qs = qs.Filter("active", active)
    }
    qs = qs.OrderBy("-usage_count", "name")
    var tags []*model.Tag
    _, err := qs.All(&tags)
    return tags, err
}

func (r *TagsRepository) GetTagByID(tagID int) (*model.Tag, error) {
    o := orm.NewOrm()
    tag := &model.Tag{ID: tagID}
    err := o.Read(tag)
    if err == orm.ErrNoRows {
        return nil, nil
    }
    return tag, err
}

func (r *TagsRepository) CountUserTagUsage(tagID int) (int, error) {
    o := orm.NewOrm()
    count, err := o.QueryTable("user_tags").Filter("tag_id", tagID).Count()
    return int(count), err
}

func (r *TagsRepository) CreateTag(tag *model.Tag) error {
    o := orm.NewOrm()
    _, err := o.Insert(tag)
    return err
}

func (r *TagsRepository) UpdateTag(tag *model.Tag) error {
    o := orm.NewOrm()
    _, err := o.Update(tag)
    return err
}

func (r *TagsRepository) DeleteTag(tagID int) error {
    o := orm.NewOrm()
    _, err := o.QueryTable("tags").Filter("id", tagID).Delete()
    return err
}

func (r *TagsRepository) GetSystemTags(category string) ([]*model.Tag, error) {
    params := map[string]interface{}{"type": "system"}
    if category != "" {
        params["category"] = category
    }
    return r.ListTags(params)
}

func (r *TagsRepository) GetCustomTags(userID int) ([]*model.Tag, error) {
    o := orm.NewOrm()
    var tags []*model.Tag
    _, err := o.QueryTable("tags").Filter("type", "custom").Filter("owner_id", userID).All(&tags)
    return tags, err
}

func (r *TagsRepository) GetPopularTags(limit int, category string) ([]*model.Tag, error) {
    o := orm.NewOrm()
    qs := o.QueryTable("tags").Filter("active", true)
    if category != "" { qs = qs.Filter("category", category) }
    qs = qs.OrderBy("-usage_count", "name")
    var tags []*model.Tag
    _, err := qs.Limit(limit).All(&tags)
    return tags, err
}

func (r *TagsRepository) SearchTags(keyword, category string) ([]*model.Tag, error) {
    o := orm.NewOrm()
    qs := o.QueryTable("tags")
    if category != "" { qs = qs.Filter("category", category) }
    if keyword != "" {
        keyword = strings.TrimSpace(keyword)
        cond := orm.NewCondition()
        cond = cond.Or("name__icontains", keyword).Or("description__icontains", keyword)
        qs = qs.SetCond(cond)
    }
    var tags []*model.Tag
    _, err := qs.OrderBy("name").All(&tags)
    return tags, err
}

func (r *TagsRepository) GetTagCategories() ([]*struct{ Category string; Count int }, error) {
    o := orm.NewOrm()
    var rows []struct{ Category string; Count int }
    _, err := o.Raw("SELECT COALESCE(category, '') as category, COUNT(*) as count FROM tags GROUP BY COALESCE(category, '')").QueryRows(&rows)
    
    // 转换为指针切片
    result := make([]*struct{ Category string; Count int }, len(rows))
    for i := range rows {
        result[i] = &rows[i]
    }
    return result, err
}

func (r *TagsRepository) IsTagInUse(tagID int) (bool, error) {
    o := orm.NewOrm()
    count, err := o.QueryTable("user_tags").Filter("tag_id", tagID).Count()
    return count > 0, err
}

func (r *TagsRepository) ExistsTagNameForUser(name string, userID int) (bool, error) {
    o := orm.NewOrm()
    count, err := o.QueryTable("tags").Filter("name", name).Filter("type", "custom").Filter("owner_id", userID).Count()
    return count > 0, err
}


