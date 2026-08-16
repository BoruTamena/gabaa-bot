package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/gin-gonic/gin"
)

// RegisterPublicStoryRoutes registers public (unauthenticated) story ad routes.
func RegisterPublicStoryRoutes(api *gin.RouterGroup, h *product.StoryHandler) {
	// Static paths before /stories/:id so they are not captured as IDs.
	api.GET("/stories", h.PublicGetStoreStories)
	api.GET("/stories/active", h.PublicListActiveStories)
	api.GET("/stories/:id", h.PublicGetStory)
}

// RegisterStoryRoutes registers protected story ad routes under the API group.
func RegisterStoryRoutes(api *gin.RouterGroup, h *product.StoryHandler) {
	api.POST("/my-store/stories", h.CreateStory)
	api.GET("/my-store/stories", h.ListMyStories)
	api.GET("/my-store/stories/:id", h.GetMyStory)
	api.PUT("/my-store/stories/:id", h.UpdateStory)
	api.DELETE("/my-store/stories/:id", h.DeleteStory)
}
