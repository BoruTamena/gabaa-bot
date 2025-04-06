package dto

type PaymentWebHook struct {
	Id                string          `json:"uuid" `
	Nonce             string          `json:"nonce "`
	Phone             string          `json:"phone"`
	PaymentMethod     string          `json:"paymentMethod"`
	TotalAmount       int             `json:"totalAmount"`
	TransactionStatus string          `json:"transactionStatus"`
	Transaction       TransactionInfo `json:"transaction"`
	NotificationUrl   string          `json:"notificationUrl"`
	SessionId         string          `json:"sessionId"`
}

type TransactionInfo struct {
	Id                string `json:"transactionId"`
	TransactionStatus string `json:"transactionStatus"`
}
