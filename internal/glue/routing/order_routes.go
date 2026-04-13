package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/order"
	"github.com/gin-gonic/gin"
)

// RegisterOrderRoutes registers all order-related routes under the protected API group.
func RegisterOrderRoutes(api *gin.RouterGroup, orderHandler *order.OrderHandler) {
	api.POST("/order/cart/add", orderHandler.AddToCart)
	api.POST("/order/create", orderHandler.Checkout)
	api.GET("/store/:store_id/orders", orderHandler.ListOrders)
	api.GET("/user/orders", orderHandler.GetUserOrders)
	api.GET("/user/cart", orderHandler.GetUserCart)
}
