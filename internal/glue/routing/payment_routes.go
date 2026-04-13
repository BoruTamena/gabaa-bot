package routing

import (
	"github.com/BoruTamena/gabaa-bot/internal/handler/payment"
	"github.com/gin-gonic/gin"
)

// RegisterPaymentRoutes registers all payment-related routes under the protected API group.
func RegisterPaymentRoutes(api *gin.RouterGroup, paymentHandler *payment.PaymentHandler) {
	api.POST("/payment/verify", paymentHandler.VerifyPayment)
	api.GET("/store/:store_id/wallet", paymentHandler.GetWallet)
}
