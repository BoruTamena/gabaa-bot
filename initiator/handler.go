package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/address"
	"github.com/BoruTamena/gabaa-bot/internal/handler/auth"
	"github.com/BoruTamena/gabaa-bot/internal/handler/cart"
	"github.com/BoruTamena/gabaa-bot/internal/handler/delivery"
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/BoruTamena/gabaa-bot/internal/handler/payment"
	"github.com/BoruTamena/gabaa-bot/internal/handler/preference"
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/BoruTamena/gabaa-bot/internal/handler/telegram"
	"github.com/BoruTamena/gabaa-bot/internal/handler/upload"
)

type Handler struct {
	AuthHandler      *auth.AuthHandler
	StoreHandler     *store.StoreHandler
	AnalyticsHandler *store.AnalyticsHandler
	ProductHandler   *product.ProductHandler
	OrderHandler     *order.OrderHandler
	CartHandler      *cart.CartHandler
	PaymentHandler   *payment.PaymentHandler
	CategoryHandler  *product.CategoryHandler
	AuthMiddleware   *middleware.AuthMiddleware
	WebhookHandler   *telegram.WebhookHandler
	UploadHandler    *upload.UploadHandler
	AddressHandler   *address.AddressHandler
	StoryHandler     *product.StoryHandler
	FavoriteHandler  *product.FavoriteHandler
	PreferenceHandler *preference.PreferenceHandler
	DeliveryHandler   *delivery.DeliveryHandler
}

func InitHandler(module Module, platform PlatFormLayer) Handler {
	return Handler{
		AuthHandler:      auth.NewAuthHandler(module.AuthModule),
		StoreHandler:     store.NewStoreHandler(module.StoreModule),
		AnalyticsHandler:  store.NewAnalyticsHandler(module.AnalyticsModule),
		ProductHandler:   product.NewProductHandler(module.ProductModule),
		OrderHandler:     order.NewOrderHandler(module.OrderModule),
		CartHandler:      cart.NewCartHandler(module.CartModule),
		PaymentHandler:   payment.NewPaymentHandler(module.PaymentModule, module.WalletModule),
		CategoryHandler:  product.NewCategoryHandler(module.CategoryModule),
		AuthMiddleware:   middleware.NewAuthMiddleware(platform.tg, module.StoreModule),
		WebhookHandler:   telegram.NewWebhookHandler(platform.tg),
		UploadHandler:    upload.NewUploadHandler(module.UploadModule),
		AddressHandler:   address.NewAddressHandler(module.AddressModule),
		StoryHandler:     product.NewStoryHandler(module.StoryModule),
		FavoriteHandler:  product.NewFavoriteHandler(module.FavoriteModule),
		PreferenceHandler: preference.NewPreferenceHandler(module.RecommendationModule),
		DeliveryHandler:   delivery.NewDeliveryHandler(module.DeliveryModule, module.OrderModule),
	}
}

