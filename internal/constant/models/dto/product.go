package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Product struct {
	ID          int64    `json:"id"`
	SellerID    int64    `json:"seller_id"`
	StoreID     *int64   `json:"store_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
	IsPosted    bool     `json:"is_posted"`
	IsBoosted   bool     `json:"is_boosted"`
}

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
	IsPosted    bool     `json:"is_posted"`
}

func (r CreateProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Price, validation.Required, validation.Min(0.0)),
		validation.Field(&r.Stock, validation.Required, validation.Min(0)),
	)
}

type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
}

type ProductFilterParams struct {
	PaginationParams
	Category string `form:"category"`
	Query    string `form:"query"`
}

func (r UpdateProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Price, validation.Min(0.0)),
		validation.Field(&r.Stock, validation.Min(0)),
	)
}
