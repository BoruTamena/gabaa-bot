package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/BoruTamena/gabaa-bot/internal/handler/payment"
)

type Handler struct {
	StoreHandler   *store.StoreHandler
	ProductHandler *product.ProductHandler
	OrderHandler   *order.OrderHandler
	PaymentHandler *payment.PaymentHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func InitHandler(module Module, platform PlatFormLayer) Handler {
	return Handler{
		StoreHandler:   store.NewStoreHandler(module.StoreModule),
		ProductHandler: product.NewProductHandler(module.ProductModule),
		OrderHandler:   order.NewOrderHandler(module.OrderModule),
		PaymentHandler: payment.NewPaymentHandler(module.OrderModule, module.WalletModule),
		AuthMiddleware: middleware.NewAuthMiddleware(platform.tg),
	}
}

