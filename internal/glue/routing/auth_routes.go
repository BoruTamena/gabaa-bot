package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/auth"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers all authentication-related routes.
func RegisterAuthRoutes(api *gin.RouterGroup, authHandler *auth.AuthHandler) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/telegram", authHandler.TelegramAuth)
		authGroup.POST("/admin/login", authHandler.AdminLogin)
		authGroup.POST("/telegram/session", authHandler.StartTelegramLoginSession)
		authGroup.GET("/telegram/session/:sessionId", authHandler.PollTelegramLoginSession)
	}
}
