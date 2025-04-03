package module

import (
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"gopkg.in/telebot.v4"
)

// define module interface here

type ProductModule interface {
	CreateProduct(c telebot.Context, product dto.Product) error
}

type OrderModule interface {
	CreateOrder() error
}
