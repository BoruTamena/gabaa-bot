package module

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
)

// define module interface here

type ProductModule interface {
	CreateProduct(ctx context.Context, product dto.Product) error
}
