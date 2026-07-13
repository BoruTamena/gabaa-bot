package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	deliveryModule module.DeliveryModule
	orderModule    module.OrderModule
}

func NewDeliveryHandler(dModule module.DeliveryModule, oModule module.OrderModule) *DeliveryHandler {
	return &DeliveryHandler{deliveryModule: dModule, orderModule: oModule}
}

func (h *DeliveryHandler) ListAreaPresets(c *gin.Context) {
	presets, err := h.deliveryModule.ListAreaPresets(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, presets)
}

func (h *DeliveryHandler) ConnectAgent(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	userID := c.GetInt64("user_id")
	var req dto.ConnectDeliveryAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	agent, err := h.deliveryModule.ConnectAgent(c.Request.Context(), storeID, userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, agent)
}

func (h *DeliveryHandler) ListAgents(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	agents, err := h.deliveryModule.ListAgents(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, agents)
}

func (h *DeliveryHandler) UpdateAgent(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	agentID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.UpdateDeliveryAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	agent, err := h.deliveryModule.UpdateAgent(c.Request.Context(), storeID, agentID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, agent)
}

func (h *DeliveryHandler) AddRoute(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	agentID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.AddDeliveryRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	route, err := h.deliveryModule.AddRoute(c.Request.Context(), storeID, agentID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, route)
}

func (h *DeliveryHandler) UpdateRoute(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	routeID, _ := strconv.ParseInt(c.Param("route_id"), 10, 64)
	var req dto.UpdateDeliveryRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	route, err := h.deliveryModule.UpdateRoute(c.Request.Context(), storeID, routeID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, route)
}

func (h *DeliveryHandler) DeleteRoute(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	routeID, _ := strconv.ParseInt(c.Param("route_id"), 10, 64)
	if err := h.deliveryModule.DeleteRoute(c.Request.Context(), storeID, routeID); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "route deleted"})
}

func (h *DeliveryHandler) DisconnectAgent(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	agentID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.deliveryModule.DisconnectAgent(c.Request.Context(), storeID, agentID); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "agent disconnected"})
}

func (h *DeliveryHandler) ListSharedAgents(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	agents, err := h.deliveryModule.ListSharedAgents(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, agents)
}

func (h *DeliveryHandler) AdoptAgent(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	userID := c.GetInt64("user_id")
	agentID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	agent, err := h.deliveryModule.AdoptAgent(c.Request.Context(), storeID, userID, agentID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, agent)
}

func (h *DeliveryHandler) GetDeliverySuggestions(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	orderID, _ := strconv.ParseInt(c.Param("order_id"), 10, 64)
	suggestions, err := h.deliveryModule.GetDeliverySuggestions(c.Request.Context(), storeID, orderID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, suggestions)
}

func (h *DeliveryHandler) GetProfile(c *gin.Context) {
	agentID := c.GetInt64("delivery_agent_id")
	profile, err := h.deliveryModule.GetProfile(c.Request.Context(), agentID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, profile)
}

func (h *DeliveryHandler) ListAssignedOrders(c *gin.Context) {
	agentID := c.GetInt64("delivery_agent_id")
	status := c.Query("status")
	var params dto.PaginationParams
	_ = c.ShouldBindQuery(&params)
	orders, err := h.deliveryModule.ListAssignedOrders(c.Request.Context(), agentID, status, params)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, orders)
}

func (h *DeliveryHandler) GetAssignedOrder(c *gin.Context) {
	agentID := c.GetInt64("delivery_agent_id")
	orderID, _ := strconv.ParseInt(c.Param("order_id"), 10, 64)
	order, err := h.deliveryModule.GetAssignedOrder(c.Request.Context(), agentID, orderID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, order)
}

func (h *DeliveryHandler) UpdateDeliveryOrderStatus(c *gin.Context) {
	agentID := c.GetInt64("delivery_agent_id")
	orderID, _ := strconv.ParseInt(c.Param("order_id"), 10, 64)
	var req dto.DeliveryOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		response.Error(c, fmt.Errorf("status is required"))
		return
	}
	if err := h.deliveryModule.UpdateDeliveryOrderStatus(c.Request.Context(), agentID, orderID, req.Status); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "order status updated to " + req.Status})
}
