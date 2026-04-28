package db

type CartItem struct {
	UserID    int64 `gorm:"primaryKey;column:user_id"`
	ProductID int64 `gorm:"primaryKey;column:product_id"`
	Quantity  int   `gorm:"column:quantity;not null"`
}
