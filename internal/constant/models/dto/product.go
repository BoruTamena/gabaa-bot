package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Product struct {
	ID          string  `json:"id ,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (p Product) Validate() error {

	return validation.ValidateStruct(
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Description, validation.Required, validation.Min(5)),
		validation.Field(&p.Price, validation.Required, is.Float),
	)

}
