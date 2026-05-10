package model

// AICallLog AI 调用日志
type AICallLog struct {
	ID               uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint64 `gorm:"not null;index:idx_user_id" json:"userId"`
	Provider         string `gorm:"size:64;default:''" json:"provider"`
	Model            string `gorm:"size:128;default:''" json:"model"`
	RequestType      string `gorm:"size:64;default:''" json:"requestType"`
	PromptTokens     int    `gorm:"default:0" json:"promptTokens"`
	CompletionTokens int    `gorm:"default:0" json:"completionTokens"`
	TotalTokens      int    `gorm:"default:0" json:"totalTokens"`
	Status           string `gorm:"type:enum('success','failed');default:success" json:"status"`
	ErrorMessage     string `gorm:"size:1024;default:''" json:"errorMessage"`
	LatencyMs        int    `gorm:"default:0" json:"latencyMs"`
	CreatedAt        string `gorm:"not null" json:"createdAt"`
}

func (AICallLog) TableName() string {
	return "ai_call_logs"
}
