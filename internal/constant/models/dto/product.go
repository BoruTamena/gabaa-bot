package dto


type Product struct {
	ID          int64   `json:"id"`
	StoreID     int64   `json:"store_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Images      string  `json:"images"`
}

