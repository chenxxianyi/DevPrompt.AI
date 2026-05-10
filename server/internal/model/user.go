package model

// User 用户模型，对应前端 User 接口
type User struct {
	BaseModel
	Username           string  `gorm:"size:64;not null;uniqueIndex:idx_username" json:"username"`
	Email              string  `gorm:"size:128;not null;uniqueIndex:idx_email" json:"email"`
	PasswordHash       string  `gorm:"size:256;not null" json:"-"`
	Avatar             string  `gorm:"size:512;default:''" json:"avatar"`
	Role               string  `gorm:"type:enum('user','admin');default:user" json:"role"`
	MembershipLevel    string  `gorm:"type:enum('free','pro','team','enterprise');default:free" json:"membershipLevel"`
	MembershipExpiredAt *string `gorm:"null" json:"membershipExpiredAt"`
	DailyGenerateCount int     `gorm:"default:0" json:"dailyGenerateCount"`
	LastGenerateDate   *string `gorm:"null" json:"lastGenerateDate"`
	Status             string  `gorm:"type:enum('active','disabled');default:active" json:"status"`
}

func (User) TableName() string {
	return "users"
}
