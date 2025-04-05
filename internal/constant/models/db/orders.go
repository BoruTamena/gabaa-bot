package db

import "github.com/google/uuid"

type Order struct {
	BaseModel
	BuyerID    int64       `gorm:"type:bigint;not null;constraint:OnDelete:CASCADE" json:"buyer_id"`
	Status     string      `gorm:"type:text;default:'pending';check:status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')" json:"status"`
	TotalPrice float64     `gorm:"type:decimal(10,2);not null" json:"total_price"`
	OrderItems []OrderItem `gorm:"foreignkey:OrderId" `
}

type OrderItem struct {
	BaseModel
	OrderId   uuid.UUID `gorm:"type:uuid;not null"`
	ProductId uuid.UUID `gorm:"type:uuid;not null"`
	Price     float64
	Quantity  int
	Product   Product `gorm:"foreignKey:ProductId;"`
}
