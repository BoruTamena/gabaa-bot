package db

import (
)


type User struct {
	BaseModel
	TelegramUserID int64  `gorm:"column:telegram_user_id;uniqueIndex;not null" json:"telegram_user_id"`
	Username       string `gorm:"column:username" json:"username"`
	Role           string `gorm:"column:role;not null" json:"role"` // 'admin' or 'customer'
}

