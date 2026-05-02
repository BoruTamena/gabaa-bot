package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/address"
	"github.com/gin-gonic/gin"
)

// RegisterAddressRoutes registers all address-related routes under the protected API group.
func RegisterAddressRoutes(api *gin.RouterGroup, addressHandler *address.AddressHandler) {
	api.POST("/user/addresses", addressHandler.CreateAddress)
	api.GET("/user/addresses", addressHandler.GetAddresses)
	api.PUT("/user/addresses/:id", addressHandler.UpdateAddress)
	api.DELETE("/user/addresses/:id", addressHandler.DeleteAddress)
	api.PUT("/user/addresses/:id/default", addressHandler.SetDefaultAddress)
}
