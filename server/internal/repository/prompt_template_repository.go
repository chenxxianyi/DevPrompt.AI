package repository

import (
	"devprompt-ai/internal/model"
	"errors"

	"gorm.io/gorm"
)

type PromptTemplateRepository struct {
	db *gorm.DB
}

func NewPromptTemplateRepository(db *gorm.DB) *PromptTemplateRepository {
	return &PromptTemplateRepository{db: db}
}

// List 分页搜索查询
func (r *PromptTemplateRepository) List(keyword, categorySlug, sort string, page, pageSize int) ([]model.PromptTemplate, int64, error) {
	var templates []model.PromptTemplate
	var total int64

	query := r.db.Model(&model.PromptTemplate{}).Where("prompt_templates.status = ?", "active")

	// 关键词搜索
	if keyword != "" {
		query = query.Where("prompt_templates.title LIKE ? OR prompt_templates.description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 分类筛选 (通过 slug 关联)
	if categorySlug != "" {
		query = query.Joins("JOIN prompt_categories ON prompt_templates.category_id = prompt_categories.id").
			Where("prompt_categories.slug = ?", categorySlug)
	}

	// 计数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	switch sort {
	case "new":
		query = query.Order("prompt_templates.created_at DESC")
	case "likes":
		query = query.Order("prompt_templates.like_count DESC")
	default: // "hot"
		query = query.Order("prompt_templates.use_count DESC")
	}

	// 分页
	offset := (page - 1) * pageSize
	err := query.Preload("Category").Offset(offset).Limit(pageSize).Find(&templates).Error
	return templates, total, err
}

// FindByID 根据 ID 查询
func (r *PromptTemplateRepository) FindByID(id uint64) (*model.PromptTemplate, error) {
	var t model.PromptTemplate
	err := r.db.Preload("Category").First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

// FindBySlug 根据 slug 查询
func (r *PromptTemplateRepository) FindBySlug(slug string) (*model.PromptTemplate, error) {
	var t model.PromptTemplate
	err := r.db.Preload("Category").Where("slug = ?", slug).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

// IncrementUseCount 增加使用计数
func (r *PromptTemplateRepository) IncrementUseCount(id uint64) error {
	return r.db.Model(&model.PromptTemplate{}).Where("id = ?", id).
		UpdateColumn("use_count", gorm.Expr("use_count + 1")).Error
}

// IncrementLikeCount 增加点赞计数
func (r *PromptTemplateRepository) IncrementLikeCount(id uint64, delta int) error {
	return r.db.Model(&model.PromptTemplate{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

// IncrementFavoriteCount 增加收藏计数
func (r *PromptTemplateRepository) IncrementFavoriteCount(id uint64, delta int) error {
	return r.db.Model(&model.PromptTemplate{}).Where("id = ?", id).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", delta)).Error
}

// Create 创建模板
func (r *PromptTemplateRepository) Create(t *model.PromptTemplate) error {
	return r.db.Create(t).Error
}

// Update 更新模板
func (r *PromptTemplateRepository) Update(t *model.PromptTemplate) error {
	return r.db.Save(t).Error
}

// Delete 软删除模板
func (r *PromptTemplateRepository) Delete(id uint64) error {
	return r.db.Delete(&model.PromptTemplate{}, id).Error
}

// AdminList 管理后台列表（包含禁用项）
func (r *PromptTemplateRepository) AdminList(page, pageSize int) ([]model.PromptTemplate, int64, error) {
	var templates []model.PromptTemplate
	var total int64

	if err := r.db.Model(&model.PromptTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Preload("Category").Order("id DESC").Offset(offset).Limit(pageSize).Find(&templates).Error
	return templates, total, err
}
