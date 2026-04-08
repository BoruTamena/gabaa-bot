package db

import "time"

type Store struct {
	BaseModel
	SellerID int64  `gorm:"column:seller_id;not null" json:"seller_id"`
	ChatID   int64  `gorm:"column:chat_id;uniqueIndex;not null" json:"chat_id"`
	ChatType string `gorm:"column:chat_type;not null" json:"chat_type"`
	Name     string `gorm:"column:name;not null" json:"name"`
	Seller   User   `gorm:"foreignKey:SellerID;references:ID" json:"seller"`
}


type Wallet struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StoreID   int64     `gorm:"column:store_id;not null" json:"store_id"`
	Balance   float64   `gorm:"column:balance;type:numeric;default:0" json:"balance"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Store     Store     `gorm:"foreignKey:StoreID;references:ID" json:"store"`
}
