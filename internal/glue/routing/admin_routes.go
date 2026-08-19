package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/admin"
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(
	api *gin.RouterGroup,
	adminHandler *admin.Handler,
	storeHandler *store.StoreHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	group := api.Group("/admin")
	group.Use(authMiddleware.PlatformAdminAuth())
	{
		group.GET("/stores", adminHandler.ListStores)
		group.GET("/stores/:store_id", adminHandler.GetStore)
		group.PATCH("/stores/:store_id/status", adminHandler.UpdateStoreStatus)
		group.POST("/stores/:store_id/kyc", adminHandler.UpsertStoreKYC)

		group.GET("/stores/:store_id/orders", adminHandler.ListStoreOrders)
		group.GET("/stores/:store_id/wallet", adminHandler.GetStoreWallet)
		group.GET("/stores/:store_id/wallet/withdrawals", adminHandler.ListStoreWithdrawals)
		group.GET("/stores/:store_id/transactions", adminHandler.ListStoreTransactions)
		group.GET("/stores/:store_id/stories", adminHandler.ListStoreStories)
		group.GET("/stores/:store_id/deliveries", adminHandler.ListStoreDeliveries)

		group.GET("/store-verifications", adminHandler.ListStoreVerifications)
		group.POST("/store-verifications/:store_id/approve", storeHandler.ApproveStoreVerification)
		group.POST("/store-verifications/:store_id/reject", storeHandler.RejectStoreVerification)

		group.GET("/orders", adminHandler.ListOrders)
		group.GET("/users", adminHandler.ListUsers)
	}
}
