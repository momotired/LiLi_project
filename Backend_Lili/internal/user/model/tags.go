package model

import (
    "time"

    "github.com/beego/beego/v2/client/orm"
)

type Tag struct {
    ID          int       `orm:"column(id);auto;pk" json:"id"`
    Name        string    `orm:"column(name);size(100)" json:"name"`
    Description string    `orm:"column(description);size(500);null" json:"description"`
    Category    string    `orm:"column(category);size(100);null" json:"category"`
    Color       string    `orm:"column(color);size(20);null" json:"color"`
    Icon        string    `orm:"column(icon);size(200);null" json:"icon"`
    Type        string    `orm:"column(type);size(20);default(system)" json:"type"` // system/custom
    Active      bool      `orm:"column(active);default(true)" json:"active"`
    UsageCount  int       `orm:"column(usage_count);default(0)" json:"usage_count"`
    OwnerID     int       `orm:"column(owner_id);null" json:"owner_id"` // 自定义标签所属用户
    CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
    UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

type UserTag struct {
	ID     int `orm:"column(id);auto;pk" json:"id"`
	UserID int `orm:"column(user_id);index" json:"user_id"`
	TagID  int `orm:"column(tag_id);index" json:"tag_id"`
}

func (t *Tag) TableName() string {
	return "tags"
}

func (ut *UserTag) TableName() string {
	return "user_tags"
}

func GetTagsByUserID(userID int) ([]Tag, error) {
	o := orm.NewOrm()
	var tags []Tag
	_, err := o.Raw(`SELECT t.id, t.name FROM tags t INNER JOIN user_tags ut ON t.id = ut.tag_id WHERE ut.user_id = ?`, userID).QueryRows(&tags)
	return tags, err
}

func UpdateUserTags(userID int, tagIDs []int) error {
	o := orm.NewOrm()
	// 先删除原有标签
	o.QueryTable("user_tags").Filter("user_id", userID).Delete()
	// 添加新标签
	for _, tagID := range tagIDs {
		ut := &UserTag{UserID: userID, TagID: tagID}
		o.Insert(ut)
	}
	return nil
}
