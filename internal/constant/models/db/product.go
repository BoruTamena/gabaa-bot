package db

type Product struct {

	BaseModel
	SellerID    int64   `gorm:"column:seller_id;not null" json:"seller_id"`
	StoreID     *int64  `gorm:"column:store_id" json:"store_id"`
	Name        string  `gorm:"column:name;not null" json:"name"`
	Description string  `gorm:"column:description" json:"description"`
	Price       float64 `gorm:"column:price;type:numeric;not null" json:"price"`
	Stock       int     `gorm:"column:stock;not null" json:"stock"`
	Category    string  `gorm:"column:category" json:"category"`
	Images      string  `gorm:"column:images;type:jsonb" json:"images"`
	Status      string  `gorm:"column:status;default:draft" json:"status"` // draft, published, archived
	IsPosted    bool    `gorm:"column:is_posted;default:false" json:"is_posted"`
	IsBoosted   bool    `gorm:"column:is_boosted;default:false" json:"is_boosted"`
}



