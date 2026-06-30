package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/gin-gonic/gin"
)

// RegisterFavoriteRoutes registers protected favorite routes under the API group.
func RegisterFavoriteRoutes(api *gin.RouterGroup, h *product.FavoriteHandler) {
	api.POST("/favorites/:product_id", h.AddFavorite)
	api.DELETE("/favorites/:product_id", h.RemoveFavorite)
	api.GET("/favorites", h.ListUserFavorites)
}
