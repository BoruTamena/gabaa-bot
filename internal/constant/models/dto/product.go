package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Product struct {
	ID          string  `json:"id ,omitempty"`
	SellerId    int64   `json:"seller_id,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (p Product) Validate() error {

	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Description, validation.Required),
		validation.Field(&p.Price, validation.Required),
	)

}
