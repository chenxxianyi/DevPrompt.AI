package repository

import (
	"devprompt-ai/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据 ID 查询
func (r *UserRepository) FindByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// FindByEmail 根据邮箱查询
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// FindByUsername 根据用户名查询
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update 更新用户信息
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateMembership(userID uint64, level string, durationDays int) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return gorm.ErrRecordNotFound
	}

	base := time.Now()
	if user.MembershipExpiredAt != nil {
		if parsed, err := time.Parse("2006-01-02 15:04:05", *user.MembershipExpiredAt); err == nil && parsed.After(base) {
			base = parsed
		}
	}

	expiry := base.AddDate(0, 0, durationDays).Format("2006-01-02 15:04:05")
	updates := map[string]interface{}{
		"membership_level":      level,
		"membership_expired_at": expiry,
	}
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

// List 分页查询用户列表（管理后台用）
func (r *UserRepository) List(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}
