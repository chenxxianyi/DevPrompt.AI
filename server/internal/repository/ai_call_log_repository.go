package repository

import (
	"devprompt-ai/internal/model"

	"gorm.io/gorm"
)

type AICallLogRepository struct {
	db *gorm.DB
}

func NewAICallLogRepository(db *gorm.DB) *AICallLogRepository {
	return &AICallLogRepository{db: db}
}

func (r *AICallLogRepository) Create(log *model.AICallLog) error {
	return r.db.Create(log).Error
}

func (r *AICallLogRepository) List(page, pageSize int) ([]model.AICallLog, int64, error) {
	var logs []model.AICallLog
	var total int64

	if err := r.db.Model(&model.AICallLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}
