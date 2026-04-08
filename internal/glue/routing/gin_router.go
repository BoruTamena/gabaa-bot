package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/BoruTamena/gabaa-bot/internal/handler/payment"
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Telegram E-Commerce Bot API
// @version 1.0
// @description Backend API for the Telegram E-Commerce Bot with MiniApp support.
// @host localhost:8080
// @BasePath /

func NewGinRouter(
	storeHandler *store.StoreHandler,
	productHandler *product.ProductHandler,
	orderHandler *order.OrderHandler,
	paymentHandler *payment.PaymentHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {

	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := authMiddleware

	// Store routes
	store := r.Group("/store")
	store.Use(auth.TelegramAuth()) // Use Telegram auth for MiniApp
	{
		store.POST("/from-chat", storeHandler.CreateStoreFromChat)
		store.GET("/dashboard/:chat_id", storeHandler.GetDashboard)
		store.GET("/:store_id/products", productHandler.ListProducts)
		store.POST("/:store_id/product", productHandler.CreateProduct)
		store.PUT("/:store_id/product/:id", productHandler.UpdateProduct)
		store.DELETE("/:store_id/product/:id", productHandler.DeleteProduct)
		store.GET("/:store_id/orders", orderHandler.ListOrders)
		store.GET("/:store_id/wallet", paymentHandler.GetWallet)
	}

	// Order routes
	order := r.Group("/order")
	order.Use(auth.TelegramAuth())
	{
		order.POST("/cart/add", orderHandler.AddToCart)
		order.POST("/create", orderHandler.Checkout)
	}

	// Payment routes
	payment := r.Group("/payment")
	payment.Use(auth.TelegramAuth())
	{
		payment.POST("/verify", paymentHandler.VerifyPayment)
	}


	return r
}
