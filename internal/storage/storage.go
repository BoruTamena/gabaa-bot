package storage

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/google/uuid"
)

type UserStorage interface {
	CreateUser(ctx context.Context, userDto dto.User) error
}
type ProductStorage interface {
	CreateProduct(ctx context.Context, product dto.Product) (error, uuid.UUID)
	GetProductByID(ctx context.Context, id string) (dto.Product, error)
	// GetProductsBySellerId

}

type OrderStorage interface {
	CreateOrder(ctx context.Context, order dto.Order) (error, uuid.UUID)
	// GetOrderByID(ctx context.Context, id string) (dto.Order, error)
	// GetOrdersByUserID(ctx context.Context, userID string) ([]dto.Order, error)
	// GetOrdersByProductID(ctx context.Context, productID string) ([]dto.Order, error)
	// GetOrdersByStatus(ctx context.Context, status string) ([]dto.Order, error)
	// UpdateOrderStatus(ctx context.Context, orderID string, status string) error
	// DeleteOrder(ctx context.Context, orderID string) error
}
