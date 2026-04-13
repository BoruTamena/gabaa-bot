package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/gin-gonic/gin"
)

// RegisterPublicCategoryRoutes registers public (unauthenticated) category routes.
func RegisterPublicCategoryRoutes(r *gin.Engine, categoryHandler *product.CategoryHandler) {
	r.GET("/categories",
		categoryHandler.ListAllCategories)
}

// RegisterCategoryRoutes registers protected category routes under the API group.
func RegisterCategoryRoutes(api *gin.RouterGroup, categoryHandler *product.CategoryHandler) {

	api.GET("/store/:store_id/categories",
		categoryHandler.ListStoreCategories)

	api.POST("/store/:store_id/category",
		categoryHandler.CreateCategory)
}
