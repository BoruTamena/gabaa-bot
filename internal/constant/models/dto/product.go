package dto

type Product struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}
