package model

// PromptCategory Prompt 分类
type PromptCategory struct {
	BaseModel
	Name        string `gorm:"size:64;not null" json:"name"`
	Slug        string `gorm:"size:64;not null;uniqueIndex:idx_slug" json:"slug"`
	Description string `gorm:"size:256;default:''" json:"description"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      string `gorm:"type:enum('active','disabled');default:active" json:"status"`
}

func (PromptCategory) TableName() string {
	return "prompt_categories"
}
