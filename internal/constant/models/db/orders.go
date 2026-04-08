package db

type Order struct {
	BaseModel
	UserID     int64   `gorm:"column:user_id;not null" json:"user_id"`
	StoreID    int64   `gorm:"column:store_id;not null" json:"store_id"`
	Status     string  `gorm:"column:status;not null;default:'pending'" json:"status"`
	TotalPrice float64 `gorm:"column:total_price;type:numeric;not null" json:"total_price"`
	User       User    `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Store      Store   `gorm:"foreignKey:StoreID;references:ID" json:"store"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   int64   `gorm:"column:order_id;not null" json:"order_id"`
	ProductID int64   `gorm:"column:product_id;not null" json:"product_id"`
	Quantity  int     `gorm:"column:quantity;not null" json:"quantity"`
	Price     float64 `gorm:"column:price;type:numeric;not null" json:"price"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}


type Payment struct {
	BaseModel
	OrderID int64  `gorm:"column:order_id;not null" json:"order_id"`
	Status  string `gorm:"column:status;not null;default:'pending'" json:"status"` // pending, confirmed
	Method  string `gorm:"column:method;not null" json:"method"`
	Order   Order  `gorm:"foreignKey:OrderID;references:ID" json:"order"`
}
