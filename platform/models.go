package platform

import "time"

type PaymentRequestPayload struct {
	CancelUrl      string        `json:"cancelUrl"`
	Phone          string        `json:"phone"`
	Email          string        `json:"email"`
	Nonce          string        `json:"nonce"`
	ErrorUrl       string        `json:"errorUrl"`
	NotifyUrl      string        `json:"notifyUrl"`
	SuccessUrl     string        `json:"successUrl"`
	PaymentMethods []string      `json:"paymentMethods"`
	ExpireDate     time.Time     `json:"expireDate"`
	Items          []interface{} `json:"items"`
	Beneficiaries  []struct {
		AccountNumber string  `json:"accountNumber"`
		Bank          string  `json:"bank"`
		Amount        float64 `json:"amount"`
	} `json:"beneficiaries"`
	Lang string `json:"lang"`
}

type PaymentResponse struct {
	Error   bool     `json:"error"`
	Message string   `json:"msg"`
	Data    DataInfo `json:"data"`
}

type DataInfo struct {
	SessionId   string `json:"sessionId"`
	PaymentUrl  string `json:"paymentUrl"`
	CancelUrl   string `json:"cancelUrl"`
	TotalAmount int    `json:"totalAmount"`
}
