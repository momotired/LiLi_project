package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	ID          int       `orm:"auto;pk" json:"id"`
	OpenID      string    `orm:"size(100);unique" json:"openid"`
	UnionID     string    `orm:"size(100);null" json:"unionid"`
	SessionKey  string    `orm:"size(100);null" json:"-"`
	Nickname    string    `orm:"size(100);null" json:"nickname"`
	Avatar      string    `orm:"size(500);null" json:"avatar"`
	Gender      int       `orm:"default(0)" json:"gender"` // 0:未知 1:男 2:女
	City        string    `orm:"size(100);null" json:"city"`
	Province    string    `orm:"size(100);null" json:"province"`
	Country     string    `orm:"size(100);null" json:"country"`
	Language    string    `orm:"size(50);null" json:"language"`
	Status      int       `orm:"default(1)" json:"status"` // 1:正常 0:禁用
	LastLoginAt time.Time `orm:"null" json:"last_login_at"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt   time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
	DeletedAt   time.Time `orm:"null;type(datetime)" json:"-"`
}

func (u *User) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(User))
}

// 根据OpenID获取用户
func GetUserByOpenID(openID string) (*User, error) {
	o := orm.NewOrm()
	if openID == "" {
		return nil, errors.New("openID is required")
	}
	user := &User{}
	err := o.QueryTable("users").Filter("open_id", openID).Filter("deleted_at__isnull", true).One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// 创建用户
func CreateUser(user *User) error {
	o := orm.NewOrm()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := o.Insert(user)
	return err
}

// 更新用户
func UpdateUser(user *User) error {
	o := orm.NewOrm()
	user.UpdatedAt = time.Now()
	_, err := o.Update(user)
	return err
}

// 更新用户最后登录时间
func UpdateUserLastLogin(userID int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("users").Filter("id", userID).Update(orm.Params{
		"last_login_at": time.Now(),
		"updated_at":    time.Now(),
	})
	return err
}
