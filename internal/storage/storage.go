package storage

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
)

// define U storage interface here

type ProductStorage interface {
	CreateProduct(ctx context.Context, product dto.Product) error
}
