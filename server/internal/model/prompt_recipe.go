package model

type PromptRecipe struct {
	BaseModel
	Type         string `gorm:"type:enum('project','cursor-rules','claude-code','optimize');not null;index:idx_type" json:"type"`
	TargetTool   string `gorm:"size:64;default:'common';index:idx_target_tool" json:"targetTool"`
	Version      string `gorm:"size:32;not null;default:'v1'" json:"version"`
	Name         string `gorm:"size:256;not null" json:"name"`
	Description  string `gorm:"size:1024;default:''" json:"description"`
	SystemPrompt string `gorm:"type:text;not null" json:"systemPrompt"`
	UserTemplate string `gorm:"type:text;default:''" json:"userTemplate"`
	OutputSchema string `gorm:"type:text;default:''" json:"outputSchema"`
	QualityRubric string `gorm:"type:text;default:''" json:"qualityRubric"`
	Status       string `gorm:"type:enum('draft','active','disabled');not null;default:'draft';index:idx_status" json:"status"`
	IsDefault    bool   `gorm:"default:false" json:"isDefault"`
	CreatedBy    uint64 `gorm:"default:0" json:"createdBy"`
}

func (PromptRecipe) TableName() string {
	return "prompt_recipes"
}