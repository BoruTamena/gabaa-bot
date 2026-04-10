package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Category struct {
	ID      int64  `json:"id"`
	StoreID int64  `json:"store_id"`
	Name    string `json:"name"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (r CreateCategoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 100)),
	)
}
