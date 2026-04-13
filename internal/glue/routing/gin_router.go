package routing

import (
	_ "github.com/BoruTamena/gabaa-bot/docs"
	"github.com/BoruTamena/gabaa-bot/internal/handler/auth"
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/BoruTamena/gabaa-bot/internal/handler/payment"
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewGinRouter(
	authHandler *auth.AuthHandler,
	storeHandler *store.StoreHandler,
	productHandler *product.ProductHandler,
	orderHandler *order.OrderHandler,
	paymentHandler *payment.PaymentHandler,
	categoryHandler *product.CategoryHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {

	r := gin.Default()

	// Global middleware — must be first
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorMiddleware())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ── Public routes (no auth required) ──────────────────────────────
	RegisterAuthRoutes(r, authHandler)
	RegisterPublicProductRoutes(r, productHandler)
	RegisterPublicCategoryRoutes(r, categoryHandler)

	// ── Protected routes (JWT auth required) ──────────────────────────
	api := r.Group("/")
	api.Use(authMiddleware.JWTAuth())
	{
		RegisterStoreRoutes(api, storeHandler)
		RegisterProductRoutes(api, productHandler)
		RegisterCategoryRoutes(api, categoryHandler)
		RegisterOrderRoutes(api, orderHandler)
		RegisterPaymentRoutes(api, paymentHandler)
	}

	return r
}
