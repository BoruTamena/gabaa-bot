package payment

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	orderModule  module.OrderModule
	walletModule module.WalletModule
}

func NewPaymentHandler(oModule module.OrderModule, wModule module.WalletModule) *PaymentHandler {
	return &PaymentHandler{orderModule: oModule, walletModule: wModule}
}

// VerifyPayment verifies a manual payment and credits the store wallet
// @Summary Verify payment (manual)
// @Tags payment
// @Accept json
// @Router /payment/verify [post]
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		OrderID int64 `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		c.JSON(appErr.Status, appErr)
		return
	}

	// 1. Mark order as completed/confirmed
	err := h.orderModule.UpdateOrderStatus(c.Request.Context(), req.OrderID, "completed")
	if err != nil {
		status, appErr := errorx.ErrorResponse(err)
		c.JSON(status, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment verified and wallet credited"})
}

// GetWallet returns the wallet balance for a store
// @Summary Get wallet balance
// @Tags wallet
// @Router /store/:store_id/wallet [get]
func (h *PaymentHandler) GetWallet(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, _ := strconv.ParseInt(storeIDStr, 10, 64)

	balance, err := h.walletModule.GetBalance(c.Request.Context(), storeID)
	if err != nil {
		status, appErr := errorx.ErrorResponse(err)
		c.JSON(status, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
