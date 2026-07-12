package constant

import "strings"

// PaymentStatus is the application-level payment lifecycle state.
type PaymentStatus string

const (
	PaymentStatusInitiated PaymentStatus = "initiated"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
)

func (s PaymentStatus) String() string {
	return string(s)
}

func (s PaymentStatus) IsTerminal() bool {
	return s == PaymentStatusSuccess || s == PaymentStatusFailed
}

// GatewayPaymentStatus is the LakiPay gateway-reported payment state.
type GatewayPaymentStatus string

const (
	GatewayPaymentStatusPending   GatewayPaymentStatus = "PENDING"
	GatewayPaymentStatusSuccess   GatewayPaymentStatus = "SUCCESS"
	GatewayPaymentStatusFailed    GatewayPaymentStatus = "FAILED"
	GatewayPaymentStatusCancelled GatewayPaymentStatus = "CANCELLED"
)

const (
	WebhookEventDeposit    = "DEPOSIT"
	WebhookEventWithdrawal = "WITHDRAWAL"
)

func (s GatewayPaymentStatus) String() string {
	return string(s)
}

func ParseGatewayPaymentStatus(value string) GatewayPaymentStatus {
	return GatewayPaymentStatus(strings.ToUpper(strings.TrimSpace(value)))
}

// EscrowStatus tracks held funds until order fulfillment.
type EscrowStatus string

const (
	EscrowStatusHeld     EscrowStatus = "held"
	EscrowStatusReleased EscrowStatus = "released"
	EscrowStatusRefunded EscrowStatus = "refunded"
)

func (s EscrowStatus) String() string {
	return string(s)
}

// WithdrawalStatus is the application-level withdrawal lifecycle state.
type WithdrawalStatus string

const (
	WithdrawalStatusInitiated WithdrawalStatus = "initiated"
	WithdrawalStatusPending   WithdrawalStatus = "pending"
	WithdrawalStatusSuccess   WithdrawalStatus = "success"
	WithdrawalStatusFailed    WithdrawalStatus = "failed"
	WithdrawalStatusCancelled WithdrawalStatus = "cancelled"
)

func (s WithdrawalStatus) String() string {
	return string(s)
}

func (s WithdrawalStatus) IsTerminal() bool {
	return s == WithdrawalStatusSuccess || s == WithdrawalStatusFailed || s == WithdrawalStatusCancelled
}
