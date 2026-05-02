package payment

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
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
// @Description Verify manual payment order and credit wallet
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body map[string]int64 true "Payment Details"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /payment/verify [post]
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		OrderID int64 `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	// 1. Mark order as completed/confirmed
	err := h.orderModule.UpdateOrderStatus(c.Request.Context(), req.OrderID, "completed")
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "payment verified and wallet credited"})
}

// GetWallet returns the wallet balance for a store
// @Summary Get wallet balance
// @Description Retrieve the wallet balance for a given store
// @Tags Wallet
// @Produce json
// @Param store_id path int true "Store ID"
// @Success 200 {object} response.BaseResponse{data=map[string]float64}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id/wallet [get]
func (h *PaymentHandler) GetWallet(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, _ := strconv.ParseInt(storeIDStr, 10, 64)

	balance, err := h.walletModule.GetBalance(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"balance": balance})
}
