package routing

import (
	_ "github.com/BoruTamena/gabaa-bot/docs"
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
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewGinRouter(
	authHandler *auth.AuthHandler,
	storeHandler *store.StoreHandler,
	analyticsHandler *store.AnalyticsHandler,
	productHandler *product.ProductHandler,
	orderHandler *order.OrderHandler,
	cartHandler *cart.CartHandler,
	paymentHandler *payment.PaymentHandler,
	categoryHandler *product.CategoryHandler,
	authMiddleware *middleware.AuthMiddleware,
	webhookHandler *telegram.WebhookHandler,
	uploadHandler *upload.UploadHandler,
	addressHandler *address.AddressHandler,
	storyHandler *product.StoryHandler,
	favoriteHandler *product.FavoriteHandler,
	preferenceHandler *preference.PreferenceHandler,
	deliveryHandler *delivery.DeliveryHandler,
) *gin.Engine {

	engine := gin.Default()

	// Global middleware — must be first
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.ErrorMiddleware())

	// Swagger stays at root (outside versioned API)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := engine.Group("/api/v1")

	// ── Public routes (no auth required) ──────────────────────────────
	RegisterAuthRoutes(apiV1, authHandler)
	RegisterPublicProductRoutes(apiV1, productHandler)
	RegisterPublicCategoryRoutes(apiV1, categoryHandler)
	RegisterPublicStoryRoutes(apiV1, storyHandler)

	// Public store routes
	PublicStoreRoutes(apiV1, storeHandler, storyHandler, productHandler, orderHandler)

	// Upload routes
	RegisterUploadRoutes(apiV1, uploadHandler)

	apiV1.POST("/webhook/telegram", authMiddleware.TelegramWebhookSecret(), webhookHandler.HandleUpdate)
	RegisterLakiPayWebhook(apiV1, paymentHandler)

	// ── Protected routes (JWT auth required) ──────────────────────────
	protected := apiV1.Group("/")
	protected.Use(authMiddleware.JWTAuth())
	{
		RegisterStoreRoutes(protected, storeHandler, analyticsHandler, authMiddleware)
		RegisterProductRoutes(protected, productHandler)
		RegisterCategoryRoutes(protected, categoryHandler)
		RegisterOrderRoutes(protected, orderHandler)
		RegisterCartRoutes(protected, cartHandler)
		RegisterPaymentRoutes(protected, paymentHandler)
		RegisterAddressRoutes(protected, addressHandler)
		RegisterStoryRoutes(protected, storyHandler)
		RegisterFavoriteRoutes(protected, favoriteHandler)
		RegisterPreferenceRoutes(protected, preferenceHandler)
		RegisterProtectedUploadRoutes(protected, uploadHandler)
		RegisterDeliveryRoutes(protected, deliveryHandler, authMiddleware)
	}

	return engine
}
