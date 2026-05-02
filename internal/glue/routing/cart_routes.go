package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/cart"
	"github.com/gin-gonic/gin"
)

// RegisterCartRoutes registers all cart-related routes under the protected API group.
func RegisterCartRoutes(api *gin.RouterGroup, cartHandler *cart.CartHandler) {
	api.POST("/user/cart/add", cartHandler.AddToCart)
	api.PUT("/user/cart/update", cartHandler.UpdateCartItem)
	api.DELETE("/user/cart/remove", cartHandler.RemoveFromCart)
	api.DELETE("/user/cart/clear", cartHandler.ClearCart)
	api.GET("/user/cart", cartHandler.GetUserCart)
}
