package model

// ProjectType 项目类型
type ProjectType struct {
	BaseModel
	Name        string `gorm:"size:64;not null" json:"name"`
	Value       string `gorm:"size:64;not null;uniqueIndex" json:"value"`
	Description string `gorm:"size:256;default:''" json:"description"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      string `gorm:"type:enum('active','disabled');default:active" json:"status"`
}

func (ProjectType) TableName() string {
	return "project_types"
}
