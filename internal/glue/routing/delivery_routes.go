package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/delivery"
	"github.com/BoruTamena/gabaa-bot/internal/handler/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterDeliveryRoutes(api *gin.RouterGroup, handler *delivery.DeliveryHandler, authMiddleware *middleware.AuthMiddleware) {
	merchant := api.Group("/my-store/delivery")
	{
		merchant.GET("/area-presets", handler.ListAreaPresets)
		merchant.POST("/agents", handler.ConnectAgent)
		merchant.GET("/agents", handler.ListAgents)
		merchant.PUT("/agents/:id", handler.UpdateAgent)
		merchant.POST("/agents/:id/routes", handler.AddRoute)
		merchant.PUT("/routes/:route_id", handler.UpdateRoute)
		merchant.DELETE("/routes/:route_id", handler.DeleteRoute)
		merchant.DELETE("/agents/:id", handler.DisconnectAgent)
		merchant.GET("/shared-agents", handler.ListSharedAgents)
		merchant.POST("/agents/:id/adopt", handler.AdoptAgent)
	}

	api.GET("/my-store/orders/:order_id/delivery-suggestions", handler.GetDeliverySuggestions)

	deliveryGroup := api.Group("/delivery")
	deliveryGroup.Use(authMiddleware.DeliveryAuth())
	{
		deliveryGroup.GET("/profile", handler.GetProfile)
		deliveryGroup.GET("/orders", handler.ListAssignedOrders)
		deliveryGroup.GET("/orders/:order_id", handler.GetAssignedOrder)
		deliveryGroup.PUT("/orders/:order_id/status", handler.UpdateDeliveryOrderStatus)
	}
}
