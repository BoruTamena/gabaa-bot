package lakipay

import "strings"

type DirectPaymentRequest struct {
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	PhoneNumber string            `json:"phone_number"`
	Medium      string            `json:"medium"`
	Description string            `json:"description,omitempty"`
	Reference   string            `json:"reference"`
	CallbackURL string            `json:"callback_url"`
	Redirects   *PaymentRedirects `json:"redirects,omitempty"`
}

type PaymentRedirects struct {
	Success string `json:"success,omitempty"`
	Failure string `json:"failure,omitempty"`
}

type DirectPaymentResponse struct {
	Success bool              `json:"success"`
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    DirectPaymentData `json:"data"`
}

func (r DirectPaymentResponse) IsSuccess() bool {
	if r.Success {
		return true
	}
	return strings.ToUpper(strings.TrimSpace(r.Status)) == "SUCCESS"
}

type DirectPaymentData struct {
	TransactionID string  `json:"transaction_id"`
	Reference     string  `json:"reference"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	Medium        string  `json:"medium"`
	CreatedAt     string  `json:"created_at"`
}

type WithdrawalRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	PhoneNumber string  `json:"phone_number"`
	Medium      string  `json:"medium"`
	Reference   string  `json:"reference"`
	CallbackURL string  `json:"callback_url"`
}

type WithdrawalResponse struct {
	Success bool           `json:"success"`
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    WithdrawalData `json:"data"`
}

func (r WithdrawalResponse) IsSuccess() bool {
	if r.Success {
		return true
	}
	return strings.ToUpper(strings.TrimSpace(r.Status)) == "SUCCESS"
}

type WithdrawalData struct {
	TransactionID string  `json:"transaction_id"`
	Reference     string  `json:"reference"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	Medium        string  `json:"medium"`
	CreatedAt     string  `json:"created_at"`
}

type WebhookPayload struct {
	Event         string  `json:"event"`
	TransactionID string  `json:"transaction_id"`
	Reference     string  `json:"reference"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	Medium        string  `json:"medium"`
	Timestamp     string  `json:"timestamp"`
	Signature     string  `json:"signature"`
}
