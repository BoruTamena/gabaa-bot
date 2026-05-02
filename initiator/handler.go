package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/auth"
	"github.com/BoruTamena/gabaa-bot/internal/handler/cart"
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/BoruTamena/gabaa-bot/internal/handler/payment"
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/BoruTamena/gabaa-bot/internal/handler/telegram"
	"github.com/BoruTamena/gabaa-bot/internal/handler/upload"
)

type Handler struct {
	AuthHandler    *auth.AuthHandler
	StoreHandler   *store.StoreHandler
	ProductHandler *product.ProductHandler
	OrderHandler    *order.OrderHandler
	CartHandler     *cart.CartHandler
	PaymentHandler  *payment.PaymentHandler
	CategoryHandler *product.CategoryHandler
	AuthMiddleware  *middleware.AuthMiddleware
	WebhookHandler  *telegram.WebhookHandler
	UploadHandler   *upload.UploadHandler
}

func InitHandler(module Module, platform PlatFormLayer) Handler {
	return Handler{
		AuthHandler:    auth.NewAuthHandler(module.AuthModule),
		StoreHandler:   store.NewStoreHandler(module.StoreModule),
		ProductHandler: product.NewProductHandler(module.ProductModule),
		OrderHandler:    order.NewOrderHandler(module.OrderModule),
		CartHandler:     cart.NewCartHandler(module.CartModule),
		PaymentHandler:  payment.NewPaymentHandler(module.OrderModule, module.WalletModule),
		CategoryHandler: product.NewCategoryHandler(module.CategoryModule),
		AuthMiddleware:  middleware.NewAuthMiddleware(platform.tg),
		WebhookHandler:  telegram.NewWebhookHandler(platform.tg),
		UploadHandler:   upload.NewUploadHandler(module.UploadModule),
	}
}
