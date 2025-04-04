package db

type Order struct {
	BaseModel
	UserID     int64     ` gorm:"type:bigint" json:"user_id"`
	ProductID  string    `json:"product_id"`
	Quantity   int       `json:"quantity"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"total_price"`
	Product    []Product `json:"product" gorm:"foreignKey:TelID;references:ProductID"`
}
