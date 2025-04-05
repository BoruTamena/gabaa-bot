package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	TelID     int64          `json:" tel_id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Phone     string         `gorm:"default:null" json:"phone"`
	Role      string         `json:"role"`
	CreatedAt time.Time      ` gorm:"default:CURRENT_TIMESTAMP"  json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
