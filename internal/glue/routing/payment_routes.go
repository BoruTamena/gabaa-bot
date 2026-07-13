package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/payment"
	"github.com/gin-gonic/gin"
)

// RegisterPaymentRoutes registers all payment-related routes under the protected API group.
func RegisterPaymentRoutes(api *gin.RouterGroup, paymentHandler *payment.PaymentHandler) {
	api.GET("/store/:store_id/wallet", paymentHandler.GetWallet)

	api.GET("/my-store/wallet", paymentHandler.GetMyStoreWallet)
	api.POST("/my-store/wallet/withdraw", paymentHandler.RequestWithdrawal)
	api.GET("/my-store/wallet/withdrawals", paymentHandler.ListWithdrawals)
	api.GET("/my-store/wallet/withdrawals/:withdrawal_id", paymentHandler.GetMyStoreWithdrawal)
	api.GET("/my-store/transactions", paymentHandler.ListMyStoreTransactions)
}

// RegisterLakiPayWebhook registers the public LakiPay webhook endpoint.
func RegisterLakiPayWebhook(r *gin.Engine, paymentHandler *payment.PaymentHandler) {
	r.POST("/api/v1/webhook/lakipay", paymentHandler.HandleLakiPayWebhook)
}
