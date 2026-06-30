package db

type Favorite struct {
	BaseModel
	UserID    int64 `gorm:"column:user_id;not null;index;uniqueIndex:idx_user_product" json:"user_id"`
	ProductID int64 `gorm:"column:product_id;not null;index;uniqueIndex:idx_user_product" json:"product_id"`

	// Associations
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
