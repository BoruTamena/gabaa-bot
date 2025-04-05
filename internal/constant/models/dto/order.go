package dto

import "github.com/google/uuid"

type Order struct {
	BuyerID    int64       `json:"buyer_id"`
	Status     string      `json:"status"`
	TotalPrice float64     ` json:"total_price"`
	OrderItems []OrderItem `gorm:"foreignkey:OrderId" `
}

type OrderItem struct {
	OrderId   uuid.UUID
	ProductId uuid.UUID
	Price     float64
	Quantity  int
}
