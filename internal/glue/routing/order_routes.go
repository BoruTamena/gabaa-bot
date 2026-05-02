package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/gin-gonic/gin"
)

// RegisterOrderRoutes registers all order-related routes under the protected API group.
func RegisterOrderRoutes(api *gin.RouterGroup, orderHandler *order.OrderHandler) {
	// Customer routes
	api.POST("/order/create", orderHandler.Checkout)
	api.GET("/orders/:order_id", orderHandler.GetOrder)
	api.PUT("/user/orders/:order_id/cancel", orderHandler.CancelOrder)
	api.GET("/user/orders", orderHandler.GetUserOrders)

	// Legacy admin routes (kept for backwards compat)
	api.PUT("/store/:store_id/orders/:order_id/status", orderHandler.UpdateOrderStatus)
	api.GET("/store/:store_id/orders", orderHandler.ListOrders)

	// Merchant-scoped routes (token-based)
	api.GET("/my-store/orders", orderHandler.MyStoreListOrders)
	api.GET("/my-store/orders/:order_id", orderHandler.MyStoreGetOrder)
	api.PUT("/my-store/orders/:order_id/status", orderHandler.MyStoreUpdateOrderStatus)
}
