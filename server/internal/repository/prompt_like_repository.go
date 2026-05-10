package repository

import (
	"devprompt-ai/internal/model"
	"errors"

	"gorm.io/gorm"
)

type PromptLikeRepository struct {
	db *gorm.DB
}

func NewPromptLikeRepository(db *gorm.DB) *PromptLikeRepository {
	return &PromptLikeRepository{db: db}
}

// Exists 检查是否已点赞
func (r *PromptLikeRepository) Exists(userID, promptID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.PromptLike{}).
		Where("user_id = ? AND prompt_template_id = ?", userID, promptID).
		Count(&count).Error
	return count > 0, err
}

// Create 添加点赞
func (r *PromptLikeRepository) Create(userID, promptID uint64) error {
	like := &model.PromptLike{
		UserID:           userID,
		PromptTemplateID: promptID,
	}
	return r.db.Create(like).Error
}

// Delete 取消点赞
func (r *PromptLikeRepository) Delete(userID, promptID uint64) error {
	return r.db.Where("user_id = ? AND prompt_template_id = ?", userID, promptID).
		Delete(&model.PromptLike{}).Error
}

// FindUserLikes 查询用户点赞的模板 ID 列表
func (r *PromptLikeRepository) FindUserLikes(userID uint64, promptIDs []uint64) (map[uint64]bool, error) {
	var likes []model.PromptLike
	err := r.db.Where("user_id = ? AND prompt_template_id IN ?", userID, promptIDs).Find(&likes).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	result := make(map[uint64]bool)
	for _, l := range likes {
		result[l.PromptTemplateID] = true
	}
	return result, nil
}
