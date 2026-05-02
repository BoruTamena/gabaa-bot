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
	api.GET("/my-store/products", productHandler.ListProducts)
	api.GET("/my-store/product/:id", productHandler.GetMyProduct)

	api.POST("/my-store/product", productHandler.CreateProduct)

	api.PUT("/my-store/product/:id", productHandler.UpdateProduct)

	api.DELETE("/my-store/product/:id", productHandler.DeleteProduct)
}
