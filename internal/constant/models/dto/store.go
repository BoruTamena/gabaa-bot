package dto

type Store struct {
	ID       int64  `json:"id"`
	SellerID int64  `json:"seller_id"`
	ChatID   int64  `json:"chat_id"`
	ChatType string `json:"chat_type"`
	Name     string `json:"name"`
}
