package repository

import (
	"devprompt-ai/internal/model"
	"errors"

	"gorm.io/gorm"
)

type PromptRecipeRepository struct {
	db *gorm.DB
}

func NewPromptRecipeRepository(db *gorm.DB) *PromptRecipeRepository {
	return &PromptRecipeRepository{db: db}
}

func (r *PromptRecipeRepository) FindByID(id uint64) (*model.PromptRecipe, error) {
	var recipe model.PromptRecipe
	err := r.db.First(&recipe, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &recipe, err
}

func (r *PromptRecipeRepository) List(recipeType, targetTool, status string) ([]model.PromptRecipe, int64, error) {
	var recipes []model.PromptRecipe
	var total int64

	q := r.db.Model(&model.PromptRecipe{})
	if recipeType != "" {
		q = q.Where("type = ?", recipeType)
	}
	if targetTool != "" {
		q = q.Where("target_tool = ?", targetTool)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("type ASC, target_tool ASC, version DESC").Find(&recipes).Error
	return recipes, total, err
}

func (r *PromptRecipeRepository) FindActiveByType(recipeType string) ([]model.PromptRecipe, error) {
	var recipes []model.PromptRecipe
	err := r.db.Where("type = ? AND status = ?", recipeType, "active").
		Order("is_default DESC, version DESC").
		Find(&recipes).Error
	return recipes, err
}

func (r *PromptRecipeRepository) FindDefaultByType(recipeType string) (*model.PromptRecipe, error) {
	var recipe model.PromptRecipe
	err := r.db.Where("type = ? AND status = ? AND is_default = ?", recipeType, "active", true).First(&recipe).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &recipe, err
}

func (r *PromptRecipeRepository) Create(recipe *model.PromptRecipe) error {
	return r.db.Create(recipe).Error
}

func (r *PromptRecipeRepository) Update(recipe *model.PromptRecipe) error {
	return r.db.Save(recipe).Error
}

func (r *PromptRecipeRepository) Delete(id uint64) error {
	return r.db.Delete(&model.PromptRecipe{}, id).Error
}

func (r *PromptRecipeRepository) SetDefault(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var recipe model.PromptRecipe
		if err := tx.First(&recipe, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.PromptRecipe{}).
			Where("type = ? AND id != ?", recipe.Type, id).
			Update("is_default", false).Error; err != nil {
			return err
		}

		return tx.Model(&recipe).Update("is_default", true).Error
	})
}
