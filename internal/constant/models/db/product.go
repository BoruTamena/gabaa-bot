package db

type Product struct {

	BaseModel
	StoreID     int64   `gorm:"column:store_id;not null" json:"store_id"`
	Name        string  `gorm:"column:name;not null" json:"name"`
	Description string  `gorm:"column:description" json:"description"`
	Price       float64 `gorm:"column:price;type:numeric;not null" json:"price"`
	Stock       int     `gorm:"column:stock;not null" json:"stock"`
	Category    string  `gorm:"column:category" json:"category"`
	Images      string  `gorm:"column:images;type:jsonb" json:"images"`
}



