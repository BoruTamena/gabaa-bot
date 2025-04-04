package order

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/module"
)

type orderModule struct {
}

func InitOrderModule() module.OrderModule {

	return &orderModule{}

}

func (order *orderModule) AddToCart(cxt context.Context, productId string) error {

	return nil

}
func (order *orderModule) CreateOrder() error {

	return nil

}
