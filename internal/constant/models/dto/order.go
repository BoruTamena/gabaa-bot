package dto

type Order struct {
	UserID     int64   `json:"user_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"total_price"`
}

type OrderList struct {
	Orders []Order `json:"orders"`
}
