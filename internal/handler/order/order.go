package order

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderModule module.OrderModule
}

func NewOrderHandler(oModule module.OrderModule) *OrderHandler {
	return &OrderHandler{orderModule: oModule}
}

// AddToCart adds an item to the user's active cart
// @Summary Add to cart
// @Tags order
// @Param product_id query int true "Product ID"
// @Param quantity query int true "Quantity"
// @Router /order/cart/add [post]
func (h *OrderHandler) AddToCart(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
	quantity, _ := strconv.Atoi(c.Query("quantity"))
	userID := c.GetInt64("user_id")

	err := h.orderModule.AddToCart(c.Request.Context(), userID, productID, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added to cart"})
}

// Checkout creates an order from the user's cart
// @Summary Create order from cart
// @Tags order
// @Param store_id query int true "Store ID"
// @Router /order/create [post]
func (h *OrderHandler) Checkout(c *gin.Context) {
	storeID, _ := strconv.ParseInt(c.Query("store_id"), 10, 64)
	userID := c.GetInt64("user_id")

	order, err := h.orderModule.Checkout(c.Request.Context(), userID, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders returns all orders for a store with pagination
// @Summary List orders
// @Tags order
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Produce json
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

	response, err := h.orderModule.ListOrders(c.Request.Context(), storeID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

