package cart

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartModule module.CartModule
}

func NewCartHandler(cModule module.CartModule) *CartHandler {
	return &CartHandler{cartModule: cModule}
}

// AddToCart adds an item to the user's active cart
// @Summary Add to cart
// @Description Add a product to the cart with a specific quantity
// @Tags Cart
// @Produce json
// @Param product_id query int true "Product ID"
// @Param quantity query int true "Quantity"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/cart/add [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
	quantity, _ := strconv.Atoi(c.Query("quantity"))
	userID := c.GetInt64("user_id")

	err := h.cartModule.AddToCart(c.Request.Context(), userID, productID, quantity)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "added to cart"})
}

// GetUserCart returns the user's active cart with product details
// @Summary Get user cart
// @Description Fetch the active user's cart including product details and total price
// @Tags Cart
// @Produce json
// @Success 200 {object} response.BaseResponse{data=dto.CartResponse}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/cart [get]
func (h *CartHandler) GetUserCart(c *gin.Context) {
	userID := c.GetInt64("user_id")

	cart, err := h.cartModule.GetUserCart(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, cart)
}

// UpdateCartItem updates the quantity of an item in the cart
// @Summary Update cart item
// @Description Update the specific quantity of a cart item
// @Tags Cart
// @Produce json
// @Param product_id query int true "Product ID"
// @Param action query string true "Action (increment or decrement)"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/cart/update [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
	action := c.Query("action")
	userID := c.GetInt64("user_id")

	if action != "increment" && action != "decrement" {
		response.Error(c, fmt.Errorf("invalid action: must be increment or decrement"))
		return
	}

	err := h.cartModule.UpdateCartItem(c.Request.Context(), userID, productID, action)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "cart item updated"})
}

// RemoveFromCart removes an item from the cart completely
// @Summary Remove from cart
// @Description Remove a specific item from the cart entirely
// @Tags Cart
// @Produce json
// @Param product_id query int true "Product ID"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/cart/remove [delete]
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
	userID := c.GetInt64("user_id")

	err := h.cartModule.RemoveFromCart(c.Request.Context(), userID, productID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "item removed from cart"})
}

// ClearCart empties the cart
// @Summary Clear cart
// @Description Empty the user's active cart completely
// @Tags Cart
// @Produce json
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/cart/clear [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetInt64("user_id")

	err := h.cartModule.ClearCart(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "cart cleared"})
}
