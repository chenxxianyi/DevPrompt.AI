package model

// AIModel AI 模型配置（API Key 从环境变量/配置读取，不存数据库）
type AIModel struct {
	BaseModel
	Provider       string `gorm:"size:64;not null;uniqueIndex:idx_provider_model" json:"provider"`
	ModelName      string `gorm:"size:128;not null;uniqueIndex:idx_provider_model" json:"modelName"`
	DisplayName    string `gorm:"size:256;not null" json:"displayName"`
	APIConfigKey   string `gorm:"size:512;default:''" json:"-"` // 不暴露给前端
	IsDefault      bool   `gorm:"default:false" json:"isDefault"`
	Status         string `gorm:"type:enum('active','disabled');default:active" json:"status"`
	Priority       int    `gorm:"default:0" json:"priority"`
	TimeoutSeconds int    `gorm:"default:60" json:"timeoutSeconds"`
}

func (AIModel) TableName() string {
	return "ai_models"
}
