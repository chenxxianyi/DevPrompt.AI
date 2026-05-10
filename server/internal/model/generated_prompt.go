package model

// GeneratedPrompt 生成的 Prompt 记录
type GeneratedPrompt struct {
	BaseModel
	UserID   uint64 `gorm:"not null;index:idx_user_id" json:"userId"`
	Type     string `gorm:"type:enum('project','cursor-rules','claude-code','optimize');not null;index:idx_type" json:"type"`
	Title    string `gorm:"size:512;default:''" json:"title"`
	Input    string `gorm:"type:json" json:"input"`
	Output   string `gorm:"type:text" json:"output"`
	Model    string `gorm:"size:128;default:''" json:"model"`
	Provider string `gorm:"size:64;default:''" json:"provider"`
	Tokens   int    `gorm:"default:0" json:"tokens"`
}

func (GeneratedPrompt) TableName() string {
	return "generated_prompts"
}
