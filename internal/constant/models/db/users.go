package db

import (
)


type User struct {
	BaseModel
	TelegramUserID *int64 `gorm:"column:telegram_user_id;uniqueIndex" json:"telegram_user_id"`
	Email          string `gorm:"column:email;uniqueIndex" json:"email"`
	PasswordHash   string `gorm:"column:password_hash" json:"-"`
	Username       string `gorm:"column:username" json:"username"`
	Role           string `gorm:"column:role;not null" json:"role"` // 'merchant' or 'customer'
	PendingStoreID *int64 `gorm:"column:pending_store_id" json:"pending_store_id"`
}

