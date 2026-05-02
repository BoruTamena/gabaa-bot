package order

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderModule module.OrderModule
}

func NewOrderHandler(oModule module.OrderModule) *OrderHandler {
	return &OrderHandler{orderModule: oModule}
}

// AddToCart and GetUserCart removed

// Checkout creates an order from the user's cart
// @Summary Create order from cart
// @Description Creates a new order based on the user's current cart items for a specific store
// @Tags Order
// @Accept json
// @Produce json
// @Param request body dto.CheckoutRequest true "Checkout Data"
// @Success 200 {object} response.BaseResponse{data=dto.Order}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /order/create [post]
func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Errorf("invalid request body: %v", err))
		return
	}

	order, err := h.orderModule.Checkout(c.Request.Context(), userID, req.StoreID, req.AddressID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, order)
}

// ListOrders returns all orders for a store with pagination
// @Summary List orders
// @Description Retrieve a paginated list of all orders for a given store
// @Tags Order
// @Produce json
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id/orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, _ := strconv.ParseInt(storeIDStr, 10, 64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dto.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.orderModule.ListOrders(c.Request.Context(), storeID, params)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// GetOrder returns an order by its ID
// @Summary Get order
// @Description Retrieve details of a specific order by ID
// @Tags Order
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.BaseResponse{data=dto.Order}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /orders/:order_id [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)

	order, err := h.orderModule.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, order)
}

// CancelOrder allows a user to cancel their own pending order
// @Summary Cancel order
// @Description Cancels a pending order and restores product stock
// @Tags Order
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/orders/:order_id/cancel [put]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)
	userID := c.GetInt64("user_id")

	err := h.orderModule.CancelOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "order cancelled"})
}

// UpdateOrderStatus updates the status of an order
// @Summary Update order status
// @Description Update the status of an existing order
// @Tags Order
// @Produce json
// @Param order_id path int true "Order ID"
// @Param status query string true "New Status"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id/orders/:order_id/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)
	status := c.Query("status")

	err := h.orderModule.UpdateOrderStatus(c.Request.Context(), orderID, status)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "order status updated"})
}

// GetUserOrders returns all orders for the authenticated user
// @Summary List user orders
// @Description Retrieve a paginated list of all orders belonging to the authenticated user
// @Tags Order
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/orders [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetInt64("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dto.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.orderModule.GetUserOrders(c.Request.Context(), userID, params)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// MyStoreListOrders lists all orders for the authenticated merchant's store
// @Summary List my store orders
// @Description Retrieve paginated orders scoped to the authenticated merchant's store with optional filtering
// @Tags Merchant Order
// @Produce json
// @Param order_id query int false "Search by exact Order ID"
// @Param status query string false "Filter by status (pending, shipped, delivered, cancelled)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/orders [get]
func (h *OrderHandler) MyStoreListOrders(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		response.Error(c, fmt.Errorf("store context missing"))
		return
	}

	var filter dto.OrderFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, fmt.Errorf("invalid query params"))
		return
	}

	// Inject store ID from token — never trust client
	filter.StoreID = storeID

	resp, err := h.orderModule.GetMyStoreOrders(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// MyStoreGetOrder returns full order detail for the authenticated merchant
// @Summary Get my store order detail
// @Description Retrieve full order detail including items, product info, and customer for a merchant-owned order
// @Tags Merchant Order
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.BaseResponse{data=dto.Order}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/orders/:order_id [get]
func (h *OrderHandler) MyStoreGetOrder(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		response.Error(c, fmt.Errorf("store context missing"))
		return
	}

	orderID, _ := strconv.ParseInt(c.Param("order_id"), 10, 64)

	order, err := h.orderModule.GetMyStoreOrder(c.Request.Context(), storeID, orderID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, order)
}

// MyStoreUpdateOrderStatus updates the status of a merchant's order
// @Summary Update my store order status
// @Description Merchant updates an order status to 'shipped' or 'delivered'. Wallet is credited on 'delivered'.
// @Tags Merchant Order
// @Produce json
// @Param order_id path int true "Order ID"
// @Param status body object true "Status body" SchemaExample({"status": "shipped"})
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/orders/:order_id/status [put]
func (h *OrderHandler) MyStoreUpdateOrderStatus(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		response.Error(c, fmt.Errorf("store context missing"))
		return
	}

	orderID, _ := strconv.ParseInt(c.Param("order_id"), 10, 64)

	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
		response.Error(c, fmt.Errorf("status is required"))
		return
	}

	if err := h.orderModule.UpdateMyStoreOrderStatus(c.Request.Context(), storeID, orderID, body.Status); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "order status updated to " + body.Status})
}
