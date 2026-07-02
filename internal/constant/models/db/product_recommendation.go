package db

import "time"

type ProductRecommendation struct {
	BaseModel
	UserID    int64     `gorm:"column:user_id;not null;index;uniqueIndex:idx_user_product_rec" json:"user_id"`
	ProductID int64     `gorm:"column:product_id;not null;index;uniqueIndex:idx_user_product_rec" json:"product_id"`
	SentAt    time.Time `gorm:"column:sent_at;not null" json:"sent_at"`
}
