package storage

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/google/uuid"
)

// product interface
type ProductStorage interface {
	CreateProduct(ctx context.Context, product dto.Product) (error, uuid.UUID)
}
