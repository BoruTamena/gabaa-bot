package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/auth"
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/BoruTamena/gabaa-bot/internal/handler/payment"
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
)

type Handler struct {
	AuthHandler    *auth.AuthHandler
	StoreHandler   *store.StoreHandler
	ProductHandler *product.ProductHandler
	OrderHandler   *order.OrderHandler
	PaymentHandler *payment.PaymentHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func InitHandler(module Module, platform PlatFormLayer) Handler {
	return Handler{
		AuthHandler:    auth.NewAuthHandler(module.AuthModule),
		StoreHandler:   store.NewStoreHandler(module.StoreModule),
		ProductHandler: product.NewProductHandler(module.ProductModule),
		OrderHandler:   order.NewOrderHandler(module.OrderModule),
		PaymentHandler: payment.NewPaymentHandler(module.OrderModule, module.WalletModule),
		AuthMiddleware: middleware.NewAuthMiddleware(platform.tg),
	}
}
