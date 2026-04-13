package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/gin-gonic/gin"
)

// RegisterPublicProductRoutes registers public (unauthenticated) product routes.
func RegisterPublicProductRoutes(r *gin.Engine, productHandler *product.ProductHandler) {
	r.GET("/products",
		productHandler.PublicListProducts)

	r.GET("/product/:id",
		productHandler.PublicGetProductByID)
}

// RegisterProductRoutes registers protected product routes under the API group.
func RegisterProductRoutes(api *gin.RouterGroup, productHandler *product.ProductHandler) {
	api.GET("/store/:store_id/products",
		productHandler.ListProducts)

	api.POST("/store/:store_id/product",
		productHandler.CreateProduct)

	api.PUT("/store/:store_id/product/:id",
		productHandler.UpdateProduct)

	api.DELETE("/store/:store_id/product/:id",
		productHandler.DeleteProduct)
}
