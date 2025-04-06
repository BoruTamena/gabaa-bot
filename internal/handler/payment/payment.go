package payment

import "github.com/gin-gonic/gin"

type paymentHandler struct {
}

func InitPaymentHandler() {

}

func (p *paymentHandler) PaymentCancelHandler(*gin.Context) {

	// handler cancel payment fallback process
}

func (p *paymentHandler) PaymentErrorHandler(c *gin.Context) {
	// handler if payment fall
}

func (p *paymentHandler) PayementSuccessHandler(c *gin.Context) {

	//handle successfully payment completion

}

func (p *paymentHandler) PaymentNotifyHandler(c *gin.Context) {
	// successfull webhook notification

}
