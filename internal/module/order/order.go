package order

import "github.com/BoruTamena/gabaa-bot/internal/module"

type orderModule struct {
}

func InitOrderModule() module.OrderModule {

	return &orderModule{}

}

func (order *orderModule) CreateOrder() error {

	return nil

}
