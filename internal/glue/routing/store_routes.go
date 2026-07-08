package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/gin-gonic/gin"
)

// RegisterStoreRoutes registers all store-related routes under the protected API group.
func RegisterStoreRoutes(
	api *gin.RouterGroup,
	storeHandler *store.StoreHandler,
	analyticsHandler *store.AnalyticsHandler,
	authMiddleware *middleware.AuthMiddleware,
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

	admin := api.Group("/admin")
	admin.Use(authMiddleware.PlatformAdminAuth())
	{
		admin.GET("/store-verifications", storeHandler.ListStoreVerifications)
		admin.POST("/store-verifications/:store_id/approve", storeHandler.ApproveStoreVerification)
		admin.POST("/store-verifications/:store_id/reject", storeHandler.RejectStoreVerification)
	}
}

