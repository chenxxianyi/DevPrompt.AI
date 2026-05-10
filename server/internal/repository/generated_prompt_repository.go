package repository

import (
	"devprompt-ai/internal/model"

	"gorm.io/gorm"
)

type GeneratedPromptRepository struct {
	db *gorm.DB
}

func NewGeneratedPromptRepository(db *gorm.DB) *GeneratedPromptRepository {
	return &GeneratedPromptRepository{db: db}
}

func (r *GeneratedPromptRepository) Create(gp *model.GeneratedPrompt) error {
	return r.db.Create(gp).Error
}

func (r *GeneratedPromptRepository) FindByUserID(userID uint64, page, pageSize int) ([]model.GeneratedPrompt, int64, error) {
	var prompts []model.GeneratedPrompt
	var total int64

	query := r.db.Model(&model.GeneratedPrompt{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&prompts).Error
	return prompts, total, err
}
