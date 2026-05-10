package model

import "time"

// PromptLike 点赞关联
type PromptLike struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint64    `gorm:"not null;uniqueIndex:idx_user_prompt;index:idx_user_id" json:"userId"`
	PromptTemplateID uint64    `gorm:"not null;uniqueIndex:idx_user_prompt;index:idx_prompt_template_id" json:"promptTemplateId"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (PromptLike) TableName() string {
	return "prompt_likes"
}
