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
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {

	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(middleware.ErrorMiddleware())

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/telegram", authHandler.TelegramAuth)
	}

	// Protected routes
	api := r.Group("/")
	api.Use(authMiddleware.JWTAuth())
	{
		// Store routes
		api.POST("/store/from-chat", storeHandler.CreateStore)
		api.GET("/store/:store_id", storeHandler.GetStore)
		api.PUT("/store/:store_id", storeHandler.UpdateStore)
		api.GET("/store/dashboard/:chat_id", storeHandler.GetDashboard)

		// Product routes
		api.GET("/store/:store_id/products", productHandler.ListProducts)
		api.POST("/store/:store_id/product", productHandler.CreateProduct)
		api.PUT("/store/:store_id/product/:id", productHandler.UpdateProduct)
		api.DELETE("/store/:store_id/product/:id", productHandler.DeleteProduct)

		// Order routes
		api.POST("/order/cart/add", orderHandler.AddToCart)
		api.POST("/order/create", orderHandler.Checkout)
		api.GET("/store/:store_id/orders", orderHandler.ListOrders)

		// Payment routes
		api.POST("/payment/verify", paymentHandler.VerifyPayment)
		api.GET("/store/:store_id/wallet", paymentHandler.GetWallet)
	}

	return r
}
