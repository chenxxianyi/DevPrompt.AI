package model

// TrialRequest 试用/购买线索
type TrialRequest struct {
	BaseModel
	UserID   uint64 `gorm:"not null;index:idx_user_plan" json:"userId"`
	PlanCode string `gorm:"size:32;not null;index:idx_user_plan" json:"planCode"`
	Contact  string `gorm:"size:255" json:"contact"`
	Message  string `gorm:"size:1000" json:"message"`
	Status   string `gorm:"type:enum('pending','contacted','approved','rejected');default:pending" json:"status"`
}

func (TrialRequest) TableName() string {
	return "trial_requests"
}
