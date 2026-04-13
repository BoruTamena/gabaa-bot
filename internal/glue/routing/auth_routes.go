package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/auth"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers all authentication-related routes.
func RegisterAuthRoutes(r *gin.Engine, authHandler *auth.AuthHandler) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/telegram", authHandler.TelegramAuth)
	}
}
