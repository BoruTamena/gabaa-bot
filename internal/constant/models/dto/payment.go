package dto

import (
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
)

type Payment struct {
	ID            int64                         `json:"id"`
	OrderID       int64                         `json:"order_id"`
	Status        constant.PaymentStatus        `json:"status"`
	Method        string                        `json:"method"`
	Reference     string                        `json:"reference"`
	TransactionID *string                       `json:"transaction_id,omitempty"`
	Amount        float64                       `json:"amount"`
	Currency      string                        `json:"currency"`
	PhoneNumber   string                        `json:"phone_number"`
	Medium        string                        `json:"medium"`
	GatewayStatus constant.GatewayPaymentStatus `json:"gateway_status"`
	CreatedAt     time.Time                     `json:"created_at"`
}

type CheckoutResponse struct {
	Order   Order   `json:"order"`
	Payment Payment `json:"payment"`
}

type WebhookResult struct {
	StatusCode int
	Message    string
}
