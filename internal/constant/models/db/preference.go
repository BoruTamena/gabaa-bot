package db

type UserCategoryPreference struct {
	BaseModel
	UserID     int64  `gorm:"column:user_id;not null;uniqueIndex" json:"user_id"`
	Categories string `gorm:"column:categories;type:jsonb;not null;default:'[]'" json:"categories"`
}
