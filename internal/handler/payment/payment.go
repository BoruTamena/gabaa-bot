package payment

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentModule module.PaymentModule
	walletModule  module.WalletModule
}

func NewPaymentHandler(pModule module.PaymentModule, wModule module.WalletModule) *PaymentHandler {
	return &PaymentHandler{paymentModule: pModule, walletModule: wModule}
}

// HandleLakiPayWebhook processes LakiPay payment status webhooks
// @Summary LakiPay webhook
// @Description Receives and processes LakiPay deposit and withdrawal webhooks
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/webhook/lakipay [post]
func (h *PaymentHandler) HandleLakiPayWebhook(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	result := h.paymentModule.HandleWebhook(c.Request.Context(), rawBody)
	c.JSON(result.StatusCode, gin.H{"message": result.Message})
}

// GetWallet returns the wallet summary for a store
// @Summary Get wallet summary
// @Description Retrieve pending, available, and locked balances for a store
// @Tags Wallet
// @Produce json
// @Param store_id path int true "Store ID"
// @Success 200 {object} response.BaseResponse{data=dto.Wallet}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id/wallet [get]
func (h *PaymentHandler) GetWallet(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, _ := strconv.ParseInt(storeIDStr, 10, 64)

	wallet, err := h.walletModule.GetWalletSummary(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, wallet)
}

// GetMyStoreWallet returns the wallet summary for the authenticated merchant store
// @Summary Get my store wallet
// @Description Retrieve wallet balances for the merchant store from JWT
// @Tags Wallet
// @Produce json
// @Success 200 {object} response.BaseResponse{data=dto.Wallet}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/wallet [get]
func (h *PaymentHandler) GetMyStoreWallet(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		response.Error(c, fmt.Errorf("store context missing"))
		return
	}

	wallet, err := h.walletModule.GetWalletSummary(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, wallet)
}

// RequestWithdrawal initiates a payout from available balance
// @Summary Request withdrawal
// @Description Merchant withdraws available balance via LakiPay
// @Tags Wallet
// @Accept json
// @Produce json
// @Param request body dto.WithdrawalRequest true "Withdrawal details"
// @Success 200 {object} response.BaseResponse{data=dto.Withdrawal}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/wallet/withdraw [post]
func (h *PaymentHandler) RequestWithdrawal(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		response.Error(c, fmt.Errorf("store context missing"))
		return
	}

	var req dto.WithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Errorf("invalid request body: %v", err))
		return
	}

	withdrawal, err := h.walletModule.RequestWithdrawal(c.Request.Context(), storeID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, withdrawal)
}

// ListWithdrawals returns paginated withdrawal history for the merchant store
// @Summary List withdrawals
// @Description Retrieve withdrawal history for the authenticated merchant store
// @Tags Wallet
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/wallet/withdrawals [get]
func (h *PaymentHandler) ListWithdrawals(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		response.Error(c, fmt.Errorf("store context missing"))
		return
	}

	var params dto.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, fmt.Errorf("invalid query params"))
		return
	}

	resp, err := h.walletModule.ListWithdrawals(c.Request.Context(), storeID, params)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}
