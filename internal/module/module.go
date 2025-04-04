package module

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"gopkg.in/telebot.v4"
)

type UserModule interface {
	CreateUser(ctx context.Context, userDto dto.User) error
}

type ProductModule interface {
	CreateProduct(c telebot.Context, product dto.Product) error
}

type OrderModule interface {
	AddToCart(cxt context.Context, user_id, productId string) error
	CreateOrder(ctx context.Context, orderRequest dto.Order) error
}
