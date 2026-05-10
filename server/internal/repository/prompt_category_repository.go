package repository

import (
	"devprompt-ai/internal/model"
	"errors"

	"gorm.io/gorm"
)

type PromptCategoryRepository struct {
	db *gorm.DB
}

func NewPromptCategoryRepository(db *gorm.DB) *PromptCategoryRepository {
	return &PromptCategoryRepository{db: db}
}

func (r *PromptCategoryRepository) FindAll() ([]model.PromptCategory, error) {
	var categories []model.PromptCategory
	err := r.db.Where("status = ?", "active").Order("sort ASC").Find(&categories).Error
	return categories, err
}

func (r *PromptCategoryRepository) FindBySlug(slug string) (*model.PromptCategory, error) {
	var cat model.PromptCategory
	err := r.db.Where("slug = ?", slug).First(&cat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cat, err
}

func (r *PromptCategoryRepository) Create(cat *model.PromptCategory) error {
	return r.db.Create(cat).Error
}

func (r *PromptCategoryRepository) Update(cat *model.PromptCategory) error {
	return r.db.Save(cat).Error
}

func (r *PromptCategoryRepository) Delete(id uint64) error {
	return r.db.Delete(&model.PromptCategory{}, id).Error
}

func (r *PromptCategoryRepository) ListAll() ([]model.PromptCategory, int64, error) {
	var categories []model.PromptCategory
	var total int64
	if err := r.db.Model(&model.PromptCategory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Order("sort ASC").Find(&categories).Error
	return categories, total, err
}
