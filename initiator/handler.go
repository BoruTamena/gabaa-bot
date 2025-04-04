package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler"
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
)

type Handler struct {
	productHandler handler.Product
	orderHandler   handler.Order
	userHandler    handler.User
}

func InitHandler(module Module) Handler {

	return Handler{

		productHandler: product.InitProductHandler(module.productModule),
		orderHandler:   order.InitOrderHandler(module.orderModule),
	}

}
