package model

// PromptTemplate Prompt 模板
type PromptTemplate struct {
	BaseModel
	CategoryID    uint64   `gorm:"not null;index:idx_category_id" json:"categoryId"`
	Title         string   `gorm:"size:256;not null" json:"title"`
	Slug          string   `gorm:"size:256;not null;uniqueIndex:idx_slug" json:"slug"`
	Description   string   `gorm:"size:1024;default:''" json:"description"`
	Content       string   `gorm:"type:text" json:"content"`
	Tags          string   `gorm:"type:json" json:"-"`
	TagsJSON      []string `gorm:"-" json:"tags"`
	UseCount      int      `gorm:"default:0" json:"useCount"`
	LikeCount     int      `gorm:"default:0" json:"likeCount"`
	FavoriteCount int      `gorm:"default:0" json:"favoriteCount"`
	Status        string   `gorm:"type:enum('active','disabled');default:active" json:"status"`

	// 非数据库字段
	Category    *PromptCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	IsLiked     bool            `gorm:"-" json:"isLiked,omitempty"`
	IsFavorited bool            `gorm:"-" json:"isFavorited,omitempty"`
}

func (PromptTemplate) TableName() string {
	return "prompt_templates"
}
