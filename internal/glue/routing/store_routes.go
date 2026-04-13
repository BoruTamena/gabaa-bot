package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/store"
	"github.com/gin-gonic/gin"
)

// RegisterStoreRoutes registers all store-related routes under the protected API group.
func RegisterStoreRoutes(api *gin.RouterGroup, storeHandler *store.StoreHandler) {
	api.POST("/store/from-chat", storeHandler.CreateStore)
	api.GET("/store/:store_id", storeHandler.GetStore)
	api.PUT("/store/:store_id", storeHandler.UpdateStore)
	api.GET("/store/dashboard/:chat_id", storeHandler.GetDashboard)
}
