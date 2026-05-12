package repository

import (
	"devprompt-ai/internal/model"
	"errors"

	"gorm.io/gorm"
)

type PromptFavoriteRepository struct {
	db *gorm.DB
}

func NewPromptFavoriteRepository(db *gorm.DB) *PromptFavoriteRepository {
	return &PromptFavoriteRepository{db: db}
}

// Exists 检查是否已收藏
func (r *PromptFavoriteRepository) Exists(userID, promptID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.PromptFavorite{}).
		Where("user_id = ? AND prompt_template_id = ?", userID, promptID).
		Count(&count).Error
	return count > 0, err
}

// Create 添加收藏
func (r *PromptFavoriteRepository) Create(userID, promptID uint64) error {
	fav := &model.PromptFavorite{
		UserID:           userID,
		PromptTemplateID: promptID,
	}
	return r.db.Create(fav).Error
}

// Delete 取消收藏
func (r *PromptFavoriteRepository) Delete(userID, promptID uint64) error {
	return r.db.Where("user_id = ? AND prompt_template_id = ?", userID, promptID).
		Delete(&model.PromptFavorite{}).Error
}

// FindUserFavorites 查询用户收藏的模板 ID 列表
func (r *PromptFavoriteRepository) FindUserFavorites(userID uint64, promptIDs []uint64) (map[uint64]bool, error) {
	var favs []model.PromptFavorite
	err := r.db.Where("user_id = ? AND prompt_template_id IN ?", userID, promptIDs).Find(&favs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	result := make(map[uint64]bool)
	for _, f := range favs {
		result[f.PromptTemplateID] = true
	}
	return result, nil
}

func (r *PromptFavoriteRepository) ListUserFavorites(userID uint64, page, pageSize int) ([]uint64, int64, error) {
	var total int64
	var ids []uint64

	base := r.db.Table("prompt_favorites pf").
		Joins("JOIN prompt_templates pt ON pt.id = pf.prompt_template_id AND pt.status = ?", "active").
		Where("pf.user_id = ?", userID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := base.Order("pf.id DESC").Offset(offset).Limit(pageSize).Pluck("pf.prompt_template_id", &ids).Error; err != nil {
		return nil, 0, err
	}

	return ids, total, nil
}
