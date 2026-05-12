package model

// TrialRequest 试用/购买线索
type TrialRequest struct {
	BaseModel
	UserID    uint64 `gorm:"not null;index:idx_user_plan" json:"userId"`
	PlanCode  string `gorm:"size:32;not null;index:idx_user_plan" json:"planCode"`
	Contact   string `gorm:"size:255" json:"contact"`
	Company   string `gorm:"size:255" json:"company"`
	TeamSize  string `gorm:"size:64" json:"teamSize"`
	UseCase   string `gorm:"size:255" json:"useCase"`
	Message   string `gorm:"size:1000" json:"message"`
	Status    string `gorm:"type:enum('pending','contacted','approved','rejected');default:pending" json:"status"`
	AdminNote string `gorm:"size:1000" json:"adminNote"`
}

func (TrialRequest) TableName() string {
	return "trial_requests"
}
