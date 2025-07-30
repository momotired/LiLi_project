package model

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type UserPreferences struct {
	ID                      int       `orm:"column(id);auto;pk" json:"id"`
	UserID                  int       `orm:"column(user_id);unique" json:"user_id"`
	NotificationEnabled     bool      `orm:"column(notification_enabled);default(false)" json:"notification_enabled"`
	PriceAlertEnabled       bool      `orm:"column(price_alert_enabled);default(false)" json:"price_alert_enabled"`
	WarrantyReminderEnabled bool      `orm:"column(warranty_reminder_enabled);default(false)" json:"warranty_reminder_enabled"`
	Theme                   string    `orm:"column(theme);size(20);null" json:"theme"`
	Language                string    `orm:"column(language);size(20);null" json:"language"`
	CreatedAt               time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt               time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (u *UserPreferences) TableName() string {
	return "user_preferences"
}

func GetPreferencesByUserID(userID int) (*UserPreferences, error) {
	o := orm.NewOrm()
	prefs := &UserPreferences{}
	err := o.QueryTable("user_preferences").Filter("user_id", userID).One(prefs)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return prefs, err
}

func UpdatePreferences(prefs *UserPreferences) error {
	o := orm.NewOrm()
	prefs.UpdatedAt = time.Now()
	_, err := o.Update(prefs)
	return err
}

func CreatePreferences(prefs *UserPreferences) error {
	o := orm.NewOrm()
	prefs.CreatedAt = time.Now()
	prefs.UpdatedAt = time.Now()
	_, err := o.Insert(prefs)
	return err
}
