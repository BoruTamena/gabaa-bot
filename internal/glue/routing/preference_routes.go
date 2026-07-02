package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/preference"
	"github.com/gin-gonic/gin"
)

func RegisterPreferenceRoutes(api *gin.RouterGroup, h *preference.PreferenceHandler) {
	api.GET("/user/preferences", h.GetPreferences)
	api.PUT("/user/preferences", h.UpdatePreferences)
}
