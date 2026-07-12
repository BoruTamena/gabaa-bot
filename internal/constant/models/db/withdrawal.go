package db

import (
	"encoding/json"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
)

type Withdrawal struct {
	BaseModel
	StoreID         int64                         `gorm:"column:store_id;not null" json:"store_id"`
	Amount          float64                       `gorm:"column:amount;type:numeric;not null" json:"amount"`
	Currency        string                        `gorm:"column:currency;not null;default:'ETB'" json:"currency"`
	PhoneNumber     string                        `gorm:"column:phone_number;not null" json:"phone_number"`
	Medium          string                        `gorm:"column:medium;not null" json:"medium"`
	Reference       string                        `gorm:"column:reference;not null;uniqueIndex" json:"reference"`
	TransactionID   *string                       `gorm:"column:transaction_id" json:"transaction_id"`
	Status          constant.WithdrawalStatus       `gorm:"column:status;not null;default:'initiated'" json:"status"`
	GatewayStatus   constant.GatewayPaymentStatus   `gorm:"column:gateway_status" json:"gateway_status"`
	GatewayResponse json.RawMessage                 `gorm:"column:gateway_response;type:jsonb" json:"gateway_response"`
	Store           Store                           `gorm:"foreignKey:StoreID;references:ID" json:"store"`
}
