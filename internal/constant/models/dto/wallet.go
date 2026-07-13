package dto

import (
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
)

type Wallet struct {
	ID               int64   `json:"id"`
	StoreID          int64   `json:"store_id"`
	Currency         string  `json:"currency"`
	PendingBalance   float64 `json:"pending_balance"`
	AvailableBalance float64 `json:"available_balance"`
	LockedBalance    float64 `json:"locked_balance"`
	TotalEarned      float64 `json:"total_earned"`
	TotalWithdrawn   float64 `json:"total_withdrawn"`
}

type WithdrawalRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,oneof=ETB USD"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Medium      string  `json:"medium" binding:"required,oneof=MPESA TELEBIRR CBE ETHSWITCH"`
}

type Withdrawal struct {
	ID            int64                         `json:"id"`
	StoreID       int64                         `json:"store_id"`
	Amount        float64                       `json:"amount"`
	Currency      string                        `json:"currency"`
	PhoneNumber   string                        `json:"phone_number"`
	Medium        string                        `json:"medium"`
	Reference     string                        `json:"reference"`
	TransactionID *string                       `json:"transaction_id,omitempty"`
	Status        constant.WithdrawalStatus       `json:"status"`
	GatewayStatus constant.GatewayPaymentStatus `json:"gateway_status"`
	CreatedAt     time.Time                     `json:"created_at"`
}

func (w Withdrawal) IsTerminal() bool {
	return w.Status.IsTerminal()
}
