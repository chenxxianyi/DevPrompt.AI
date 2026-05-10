package model

// MembershipPlan 会员套餐
type MembershipPlan struct {
	BaseModel
	Name         string   `gorm:"size:64;not null" json:"name"`
	Code         string   `gorm:"type:enum('free','pro','team','enterprise');not null;uniqueIndex:idx_code" json:"code"`
	Price        float64  `gorm:"type:decimal(10,2);default:0" json:"price"`
	DurationDays int      `gorm:"default:0" json:"durationDays"`
	DailyLimit   int      `gorm:"default:5" json:"dailyLimit"`
	Features     string   `gorm:"type:json" json:"-"`
	FeaturesJSON []string `gorm:"-" json:"features"`
	Status       string   `gorm:"type:enum('active','disabled');default:active" json:"status"`
}

func (MembershipPlan) TableName() string {
	return "membership_plans"
}
