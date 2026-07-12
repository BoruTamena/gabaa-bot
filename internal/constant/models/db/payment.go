package db

import (
	"encoding/json"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
)

type Payment struct {
	BaseModel
	OrderID         int64           `gorm:"column:order_id;not null" json:"order_id"`
	Status          constant.PaymentStatus        `gorm:"column:status;not null;default:'initiated'" json:"status"`
	Method          string                        `gorm:"column:method;not null" json:"method"`
	Reference       string                        `gorm:"column:reference" json:"reference"`
	TransactionID   *string                       `gorm:"column:transaction_id" json:"transaction_id"`
	Amount          float64                       `gorm:"column:amount;type:numeric" json:"amount"`
	Currency        string                        `gorm:"column:currency;default:'ETB'" json:"currency"`
	PhoneNumber     string                        `gorm:"column:phone_number" json:"phone_number"`
	Medium          string                        `gorm:"column:medium" json:"medium"`
	GatewayStatus   constant.GatewayPaymentStatus `gorm:"column:gateway_status" json:"gateway_status"`
	GatewayResponse json.RawMessage `gorm:"column:gateway_response;type:jsonb" json:"gateway_response"`
	Order           Order           `gorm:"foreignKey:OrderID;references:ID" json:"order"`
}

type PaymentWebhook struct {
	ID            int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	PaymentID     *int64          `gorm:"column:payment_id" json:"payment_id"`
	WithdrawalID  *int64          `gorm:"column:withdrawal_id" json:"withdrawal_id"`
	TransactionID string          `gorm:"column:transaction_id" json:"transaction_id"`
	Event         string          `gorm:"column:event" json:"event"`
	Status        string          `gorm:"column:status" json:"status"`
	Payload       json.RawMessage `gorm:"column:payload;type:jsonb;not null" json:"payload"`
	Signature     string          `gorm:"column:signature" json:"signature"`
	Verified      bool            `gorm:"column:verified;default:false" json:"verified"`
	Processed     bool            `gorm:"column:processed;default:false" json:"processed"`
	ReceivedAt    time.Time       `gorm:"column:received_at" json:"received_at"`
}

type Escrow struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      int64      `gorm:"column:order_id;not null;uniqueIndex" json:"order_id"`
	StoreID      int64      `gorm:"column:store_id;not null" json:"store_id"`
	Amount       float64    `gorm:"column:amount;type:numeric;not null" json:"amount"`
	Currency     string     `gorm:"column:currency;default:'ETB'" json:"currency"`
	Status       constant.EscrowStatus `gorm:"column:status;not null;default:'held'" json:"status"`
	ReleaseAt    *time.Time `gorm:"column:release_at" json:"release_at"`
	ReleasedAt   *time.Time `gorm:"column:released_at" json:"released_at"`
	RefundAmount float64    `gorm:"column:refund_amount;type:numeric;default:0" json:"refund_amount"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
}
