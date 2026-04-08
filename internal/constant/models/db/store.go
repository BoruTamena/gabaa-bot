package db

import "time"

type Store struct {
	BaseModel
	SellerID       int64  `gorm:"column:seller_id;not null" json:"seller_id"`
	TelegramChatID int64  `gorm:"column:telegram_chat_id;uniqueIndex;not null" json:"telegram_chat_id"`
	Name           string `gorm:"column:name;not null" json:"name"`
	Category       string `gorm:"column:category;not null" json:"category"`
	Description    string `gorm:"column:description" json:"description"`
	LogoImage      string `gorm:"column:logo_image" json:"logo_image"`
	CoverImage     string `gorm:"column:cover_image" json:"cover_image"`
	Phone          string `gorm:"column:phone" json:"phone"`
	Email          string `gorm:"column:email" json:"email"`
	Location       string `gorm:"column:location" json:"location"`
	Seller         User   `gorm:"foreignKey:SellerID;references:ID" json:"seller"`
}


type Wallet struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StoreID   int64     `gorm:"column:store_id;not null" json:"store_id"`
	Balance   float64   `gorm:"column:balance;type:numeric;default:0" json:"balance"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Store     Store     `gorm:"foreignKey:StoreID;references:ID" json:"store"`
}
