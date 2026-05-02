package dto

import (
	"time"
)

type OrderCustomer struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type OrderProductDetail struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Images   []string `json:"images"`
	Category string   `json:"category"`
}

type Order struct {
	ID              int64          `json:"id"`
	StoreID         int64          `json:"store_id"`
	UserID          int64          `json:"user_id"`
	Status          string         `json:"status"`
	TotalPrice      float64        `json:"total_price"`
	CreatedAt       time.Time      `json:"created_at"`
	Customer        *OrderCustomer `json:"customer,omitempty"`
	ShippingAddress *Address       `json:"shipping_address,omitempty"`
	OrderItems      []OrderItem    `json:"order_items"`
}

type OrderItem struct {
	ID        int64               `json:"id"`
	OrderID   int64               `json:"order_id"`
	ProductID int64               `json:"product_id"`
	Quantity  int                 `json:"quantity"`
	Price     float64             `json:"price"`
	Product   *OrderProductDetail `json:"product,omitempty"`
}

type OrderFilterParams struct {
	PaginationParams
	StoreID int64  `form:"-"`       // injected server-side from token
	OrderID *int64 `form:"order_id"` // search by specific order ID
	Status  string `form:"status"`   // filter: created, shipped, delivered, cancelled
}

type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
}

type CheckoutRequest struct {
	StoreID   int64 `json:"store_id" binding:"required"`
	AddressID int64 `json:"address_id" binding:"required"`
}
