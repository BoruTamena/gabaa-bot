package db

import "time"

const (
	TelegramLoginSessionStatusPending   = "pending"
	TelegramLoginSessionStatusCompleted = "completed"
)

type TelegramLoginSession struct {
	ID             string     `gorm:"column:id;primaryKey;size:64"`
	Status         string     `gorm:"column:status;not null;default:pending"`
	TelegramUserID *int64     `gorm:"column:telegram_user_id"`
	Username       string     `gorm:"column:username"`
	ExpiresAt      time.Time  `gorm:"column:expires_at;not null"`
	CompletedAt    *time.Time `gorm:"column:completed_at"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (TelegramLoginSession) TableName() string {
	return "telegram_login_sessions"
}
