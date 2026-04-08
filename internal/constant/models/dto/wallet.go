package dto

type Wallet struct {
	ID      int64   `json:"id"`
	StoreID int64   `json:"store_id"`
	Balance float64 `json:"balance"`
}
