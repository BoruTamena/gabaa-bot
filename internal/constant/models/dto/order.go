package dto

import (
	"time"
)

type Order struct {
	ID         int64       `json:"id"`
	StoreID    int64       `json:"store_id"`
	UserID     int64       `json:"user_id"`
	Status     string      `json:"status"`
	TotalPrice float64     `json:"total_price"`
	CreatedAt  time.Time   `json:"created_at"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

