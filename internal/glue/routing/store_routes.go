package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/gin-gonic/gin"
)

// PublicStoreRoutes registers all public store-related routes.
func PublicStoreRoutes(
	api *gin.RouterGroup,
	storeHandler *store.StoreHandler,
	storyHandler *product.StoryHandler,
	productHandler *product.ProductHandler,
	orderHandler *order.OrderHandler,
) {
	// Exact /stores before /stores/:storeName
	api.GET("/stores", storeHandler.ListActiveStores)
	api.GET("/stores/:storeName", storeHandler.GetStoreDetailsByName)
	api.GET("/stores/:storeName/stories", storyHandler.PublicGetStoreStories)
	api.GET("/stores/:storeName/products", productHandler.GetStoreProducts)
	api.GET("/stores/:storeName/sales", orderHandler.GetStoreRecentSales)
}

// RegisterStoreRoutes registers all store-related routes under the protected API group.
func RegisterStoreRoutes(
	api *gin.RouterGroup,
	storeHandler *store.StoreHandler,
	analyticsHandler *store.AnalyticsHandler,
) {
	api.POST("/store/from-chat", storeHandler.CreateStore)
	api.GET("/store/:store_id", storeHandler.GetStore)
	api.GET("/store/:store_id/status", storeHandler.GetStoreStatus)
	api.PUT("/store/:store_id", storeHandler.UpdateStore)
	api.GET("/store/dashboard/:chat_id", storeHandler.GetDashboard)

	api.POST("/store/verification", storeHandler.SubmitStoreVerification)
	api.GET("/store/verification", storeHandler.GetStoreVerification)

	// Analytics routes
	analytics := api.Group("/store/analytics")
	{
		analytics.GET("/sales", analyticsHandler.GetSalesAnalytics)
		analytics.GET("/orders", analyticsHandler.GetOrderAnalytics)
		analytics.GET("/products", analyticsHandler.GetProductAnalytics)
		analytics.GET("/stories", analyticsHandler.GetStoryAnalytics)
	}
}
