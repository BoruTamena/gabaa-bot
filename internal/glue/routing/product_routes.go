package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/gin-gonic/gin"
)

// RegisterPublicProductRoutes registers public (unauthenticated) product routes.
func RegisterPublicProductRoutes(api *gin.RouterGroup, productHandler *product.ProductHandler) {
	api.GET("/products", productHandler.PublicListProducts)
	api.GET("/product/:id", productHandler.PublicGetProductByID)
}

// RegisterProductRoutes registers protected product routes under the API group.
func RegisterProductRoutes(api *gin.RouterGroup, productHandler *product.ProductHandler) {
	api.POST("/products/:id/inquiries", productHandler.CreateInquiry)
	api.GET("/my/inquiries", productHandler.ListMyInquiries)
	api.GET("/my-store/inquiries", productHandler.ListStoreInquiries)
	api.GET("/my-store/inquiries/:inquiry_id", productHandler.GetStoreInquiry)
	api.POST("/my-store/inquiries/:inquiry_id/approve", productHandler.ApproveInquiry)
	api.POST("/my-store/inquiries/:inquiry_id/reject", productHandler.RejectInquiry)

	api.GET("/my-store/products", productHandler.ListProducts)
	api.GET("/my-store/product/:id", productHandler.GetMyProduct)

	api.POST("/my-store/product", productHandler.CreateProduct)

	api.PUT("/my-store/product/:id", productHandler.UpdateProduct)

	api.DELETE("/my-store/product/:id", productHandler.DeleteProduct)
}
