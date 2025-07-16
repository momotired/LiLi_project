package repository

import (
	"lili-backend/lilibd/internal/auth/model"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type AuthRepository struct {
	o orm.Ormer
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		o: orm.NewOrm(),
	}
}

// 根据OpenID查询用户
func (r *AuthRepository) GetUserByOpenID(openID string) (*model.UserInfo, error) {
	user := &model.UserInfo{}
	err := r.o.QueryTable("users").Filter("openid", openID).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 根据用户ID查询用户
func (r *AuthRepository) GetUserByID(userID int) (*model.UserInfo, error) {
	user := &model.UserInfo{}
	err := r.o.QueryTable("users").Filter("id", userID).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 创建用户
func (r *AuthRepository) CreateUser(openID string) (*model.UserInfo, error) {
	user := &model.UserInfo{
		OpenID:   openID,
		NickName: "微信用户",
		Status:   1,
	}

	id, err := r.o.Insert(user)
	if err != nil {
		return nil, err
	}

	user.ID = int(id)
	return user, nil
}

// 更新用户信息（包括最后登录时间）
func (r *AuthRepository) UpdateUserInfo(userID int, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := r.o.QueryTable("users").Filter("id", userID).Update(orm.Params(updates))
	return err
}

// 更新用户最后登录时间
func (r *AuthRepository) UpdateLastLoginTime(userID int) error {
	_, err := r.o.QueryTable("users").Filter("id", userID).Update(orm.Params{
		"last_login_at": time.Now(),
		"updated_at":    time.Now(),
	})
	return err
}

// 保存RefreshToken
func (r *AuthRepository) SaveRefreshToken(userID int, refreshToken string) error {
	// 删除旧的 RefreshToken
	r.o.QueryTable("user_session").Filter("user_id", userID).Delete()

	// 创建新的会话记录
	session := &model.UserSession{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 7), // 7天过期
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := r.o.Insert(session)
	return err
}

// 检查RefreshToken是否存在
func (r *AuthRepository) CheckRefreshToken(userID int, refreshToken string) (bool, error) {
	count, err := r.o.QueryTable("user_session").
		Filter("user_id", userID).
		Filter("refresh_token", refreshToken).
		Filter("expires_at__gt", time.Now()).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 更新RefreshToken
func (r *AuthRepository) UpdateRefreshToken(userID int, oldRefreshToken, newRefreshToken string) error {
	_, err := r.o.QueryTable("user_session").
		Filter("user_id", userID).
		Filter("refresh_token", oldRefreshToken).
		Update(orm.Params{
			"refresh_token": newRefreshToken,
			"expires_at":    time.Now().Add(time.Hour * 24 * 7), // 7天过期
			"updated_at":    time.Now(),
		})
	return err
}

// 删除RefreshToken
func (r *AuthRepository) DeleteRefreshToken(userID int) error {
	_, err := r.o.QueryTable("user_session").Filter("user_id", userID).Delete()
	return err
}

// 添加Token到黑名单
func (r *AuthRepository) AddTokenToBlacklist(token string) error {
	// 解析Token获取过期时间
	claims, err := r.parseTokenClaims(token)
	if err != nil {
		return err
	}

	blacklist := &model.TokenBlacklist{
		Token:     token,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		CreatedAt: time.Now(),
	}

	_, err = r.o.Insert(blacklist)
	return err
}

// 检查Token是否在黑名单中
func (r *AuthRepository) IsTokenBlacklisted(token string) bool {
	count, err := r.o.QueryTable("token_blacklist").Filter("token", token).Count()
	if err != nil {
		// 导入 logs 包后使用 logs.Error
		return false
	}
	return count > 0
}

// 清理过期的黑名单Token
func (r *AuthRepository) CleanExpiredBlacklistTokens() error {
	_, err := r.o.QueryTable("token_blacklist").Filter("expires_at__lt", time.Now()).Delete()
	return err
}

// 解析Token获取Claims（内部使用）
func (r *AuthRepository) parseTokenClaims(token string) (*model.TokenClaims, error) {
	// 这里应该使用JWT工具解析Token，暂时返回一个默认值
	return &model.TokenClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}, nil
}
