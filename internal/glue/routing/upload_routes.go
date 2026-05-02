package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/upload"
	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(api *gin.RouterGroup, uploadHandler *upload.UploadHandler) {
	api.POST("/upload/images", uploadHandler.UploadImages)
}
