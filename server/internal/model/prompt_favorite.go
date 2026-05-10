package model

import "time"

// PromptFavorite 收藏关联
type PromptFavorite struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint64    `gorm:"not null;uniqueIndex:idx_user_prompt;index:idx_user_id" json:"userId"`
	PromptTemplateID uint64    `gorm:"not null;uniqueIndex:idx_user_prompt;index:idx_prompt_template_id" json:"promptTemplateId"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (PromptFavorite) TableName() string {
	return "prompt_favorites"
}
