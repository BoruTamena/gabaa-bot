package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/product"
	"github.com/gin-gonic/gin"
)

// RegisterPublicStoryRoutes registers public (unauthenticated) story ad routes.
func RegisterPublicStoryRoutes(r *gin.Engine, h *product.StoryHandler) {
	r.GET("/stories", h.PublicListActiveStories)
	r.GET("/stories/:id", h.PublicGetStory)
}

// RegisterStoryRoutes registers protected story ad routes under the API group.
func RegisterStoryRoutes(api *gin.RouterGroup, h *product.StoryHandler) {
	api.POST("/my-store/stories", h.CreateStory)
	api.GET("/my-store/stories", h.ListMyStories)
	api.GET("/my-store/stories/:id", h.GetMyStory)
	api.PUT("/my-store/stories/:id", h.UpdateStory)
	api.DELETE("/my-store/stories/:id", h.DeleteStory)
}
