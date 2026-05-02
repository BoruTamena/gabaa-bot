package db

import "time"

type Address struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64      `gorm:"column:user_id;not null" json:"user_id"`
	Label         string     `gorm:"column:label;not null;default:'home'" json:"label"` // home, work, other
	RecipientName string     `gorm:"column:recipient_name;not null" json:"recipient_name"`
	Phone         string     `gorm:"column:phone;not null" json:"phone"`
	Street        string     `gorm:"column:street;not null" json:"street"`
	City          string     `gorm:"column:city;not null" json:"city"`
	Region        string     `gorm:"column:region" json:"region"`
	Country       string     `gorm:"column:country;not null;default:'Ethiopia'" json:"country"`
	IsDefault     bool       `gorm:"column:is_default;not null;default:false" json:"is_default"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"-"`
}
