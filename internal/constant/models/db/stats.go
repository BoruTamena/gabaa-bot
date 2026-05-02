package db

import "time"

type StoreStat struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StoreID   int64     `gorm:"column:store_id;uniqueIndex;not null" json:"store_id"`
	Views     int64     `gorm:"column:views;default:0" json:"views"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Store     Store     `gorm:"foreignKey:StoreID;references:ID" json:"store"`
}
