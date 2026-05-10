package repository

import (
	"devprompt-ai/internal/model"
	"errors"

	"gorm.io/gorm"
)

type ProjectTypeRepository struct {
	db *gorm.DB
}

func NewProjectTypeRepository(db *gorm.DB) *ProjectTypeRepository {
	return &ProjectTypeRepository{db: db}
}

func (r *ProjectTypeRepository) FindActive() ([]model.ProjectType, error) {
	var types []model.ProjectType
	err := r.db.Where("status = ?", "active").Order("sort ASC").Find(&types).Error
	return types, err
}

func (r *ProjectTypeRepository) FindAll() ([]model.ProjectType, int64, error) {
	var types []model.ProjectType
	var total int64
	if err := r.db.Model(&model.ProjectType{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Order("sort ASC").Find(&types).Error
	return types, total, err
}

func (r *ProjectTypeRepository) FindByID(id uint64) (*model.ProjectType, error) {
	var pt model.ProjectType
	err := r.db.First(&pt, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pt, err
}

func (r *ProjectTypeRepository) Create(pt *model.ProjectType) error {
	return r.db.Create(pt).Error
}

func (r *ProjectTypeRepository) Update(pt *model.ProjectType) error {
	return r.db.Save(pt).Error
}

func (r *ProjectTypeRepository) Delete(id uint64) error {
	return r.db.Delete(&model.ProjectType{}, id).Error
}
