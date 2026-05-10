package repository

import (
	"devprompt-ai/internal/model"
	"errors"

	"gorm.io/gorm"
)

type AIModelRepository struct {
	db *gorm.DB
}

func NewAIModelRepository(db *gorm.DB) *AIModelRepository {
	return &AIModelRepository{db: db}
}

// FindActive 获取所有活跃的模型，按优先级排序
func (r *AIModelRepository) FindActive() ([]model.AIModel, error) {
	var models []model.AIModel
	err := r.db.Where("status = ?", "active").Order("priority ASC").Find(&models).Error
	return models, err
}

// FindDefault 获取默认模型
func (r *AIModelRepository) FindDefault() (*model.AIModel, error) {
	var mdl model.AIModel
	err := r.db.Where("is_default = ? AND status = ?", true, "active").First(&mdl).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &mdl, err
}

// List 获取所有模型
func (r *AIModelRepository) List() ([]model.AIModel, int64, error) {
	var models []model.AIModel
	var total int64
	if err := r.db.Model(&model.AIModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Order("priority ASC").Find(&models).Error
	return models, total, err
}

func (r *AIModelRepository) FindByID(id uint64) (*model.AIModel, error) {
	var mdl model.AIModel
	err := r.db.First(&mdl, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &mdl, err
}

func (r *AIModelRepository) Create(mdl *model.AIModel) error {
	return r.db.Create(mdl).Error
}

func (r *AIModelRepository) Update(mdl *model.AIModel) error {
	return r.db.Save(mdl).Error
}

func (r *AIModelRepository) Delete(id uint64) error {
	return r.db.Delete(&model.AIModel{}, id).Error
}
