package model

import (
	"github.com/beego/beego/v2/client/orm"
)

type Tag struct {
	ID   int    `orm:"column(id);auto;pk" json:"id"`
	Name string `orm:"column(name);size(100);unique" json:"name"`
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
