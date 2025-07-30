package repository

import (
	"Backend_Lili/internal/user/model"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// 用户相关操作
func (r *UserRepository) GetUserByOpenID(openID string) (*model.User, error) {
	o := orm.NewOrm()
	if openID == "" {
		return nil, errors.New("openID is required")
	}
	user := &model.User{}
	err := o.QueryTable("users").Filter("open_id", openID).Filter("deleted_at__isnull", true).One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetUserByID(userID int) (*model.User, error) {
	o := orm.NewOrm()
	user := &model.User{}
	err := o.QueryTable("users").Filter("id", userID).Filter("deleted_at__isnull", true).One(user)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) CreateUser(user *model.User) error {
	o := orm.NewOrm()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := o.Insert(user)
	return err
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	o := orm.NewOrm()
	user.UpdatedAt = time.Now()
	_, err := o.Update(user)
	return err
}

func (r *UserRepository) UpdateUserLastLogin(userID int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("users").Filter("id", userID).Update(orm.Params{
		"last_login_at": time.Now(),
		"updated_at":    time.Now(),
	})
	return err
}

func (r *UserRepository) SoftDeleteUser(userID int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("users").Filter("id", userID).Update(orm.Params{
		"deleted_at": time.Now(),
		"updated_at": time.Now(),
		"status":     0,
	})
	return err
}

// 用户偏好设置相关操作
func (r *UserRepository) GetPreferencesByUserID(userID int) (*model.UserPreferences, error) {
	o := orm.NewOrm()
	prefs := &model.UserPreferences{}
	err := o.QueryTable("user_preferences").Filter("user_id", userID).One(prefs)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return prefs, err
}

func (r *UserRepository) CreatePreferences(prefs *model.UserPreferences) error {
	o := orm.NewOrm()
	prefs.CreatedAt = time.Now()
	prefs.UpdatedAt = time.Now()
	_, err := o.Insert(prefs)
	return err
}

func (r *UserRepository) UpdatePreferences(prefs *model.UserPreferences) error {
	o := orm.NewOrm()
	prefs.UpdatedAt = time.Now()
	_, err := o.Update(prefs)
	return err
}

func (r *UserRepository) CreateOrUpdatePreferences(prefs *model.UserPreferences) error {
	existing, err := r.GetPreferencesByUserID(prefs.UserID)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.CreatePreferences(prefs)
	} else {
		prefs.ID = existing.ID
		return r.UpdatePreferences(prefs)
	}
}

// 用户标签相关操作
func (r *UserRepository) GetTagsByUserID(userID int) ([]model.Tag, error) {
	o := orm.NewOrm()
	var tags []model.Tag
	_, err := o.Raw(`SELECT t.id, t.name FROM tags t INNER JOIN user_tags ut ON t.id = ut.tag_id WHERE ut.user_id = ?`, userID).QueryRows(&tags)
	return tags, err
}

func (r *UserRepository) GetAllTags() ([]model.Tag, error) {
	o := orm.NewOrm()
	var tags []model.Tag
	_, err := o.QueryTable("tags").All(&tags)
	return tags, err
}

func (r *UserRepository) ValidateTagIDs(tagIDs []int) (bool, error) {
	if len(tagIDs) == 0 {
		return true, nil
	}

	o := orm.NewOrm()
	count, err := o.QueryTable("tags").Filter("id__in", tagIDs).Count()
	if err != nil {
		return false, err
	}

	return int(count) == len(tagIDs), nil
}

func (r *UserRepository) UpdateUserTags(userID int, tagIDs []int) error {
	o := orm.NewOrm()

	// 先删除原有标签
	_, err := o.QueryTable("user_tags").Filter("user_id", userID).Delete()
	if err != nil {
		return err
	}

	// 添加新标签
	for _, tagID := range tagIDs {
		ut := &model.UserTag{UserID: userID, TagID: tagID}
		_, err = o.Insert(ut)
		if err != nil {
			return err
		}
	}

	return nil
}

// 用户统计相关操作
func (r *UserRepository) GetUserDeviceCount(userID int) (int, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable("devices").Filter("user_id", userID).Filter("deleted_at__isnull", true).Count()
	return int(count), err
}

func (r *UserRepository) GetUserDeviceCountByCategory(userID int) (map[string]int, error) {
	o := orm.NewOrm()
	type CategoryCount struct {
		CategoryName string `json:"category_name"`
		Count        int    `json:"count"`
	}

	var results []CategoryCount
	_, err := o.Raw(`
		SELECT c.name as category_name, COUNT(d.id) as count 
		FROM categories c 
		LEFT JOIN devices d ON c.id = d.category_id AND d.user_id = ? AND d.deleted_at IS NULL
		GROUP BY c.id, c.name
	`, userID).QueryRows(&results)

	if err != nil {
		return nil, err
	}

	categoryCount := make(map[string]int)
	for _, result := range results {
		categoryCount[result.CategoryName] = result.Count
	}

	return categoryCount, nil
}

func (r *UserRepository) GetUserDevicesTotalValue(userID int) (float64, error) {
	o := orm.NewOrm()
	type ValueResult struct {
		TotalValue float64 `json:"total_value"`
	}

	var result ValueResult
	err := o.Raw(`
		SELECT COALESCE(SUM(d.purchase_price), 0) as total_value 
		FROM devices d 
		WHERE d.user_id = ? AND d.deleted_at IS NULL
	`, userID).QueryRow(&result)

	return result.TotalValue, err
}

func (r *UserRepository) GetUserDeviceAverageHoldTime(userID int) (int, error) {
	o := orm.NewOrm()
	type TimeResult struct {
		AvgDays int `json:"avg_days"`
	}

	var result TimeResult
	err := o.Raw(`
		SELECT COALESCE(AVG(DATEDIFF(COALESCE(d.sold_at, NOW()), d.purchase_date)), 0) as avg_days
		FROM devices d 
		WHERE d.user_id = ? AND d.deleted_at IS NULL AND d.purchase_date IS NOT NULL
	`, userID).QueryRow(&result)

	return result.AvgDays, err
}
