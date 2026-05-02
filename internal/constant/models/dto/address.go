package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Address struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Label         string    `json:"label"`
	RecipientName string    `json:"recipient_name"`
	Phone         string    `json:"phone"`
	Street        string    `json:"street"`
	City          string    `json:"city"`
	Region        string    `json:"region"`
	Country       string    `json:"country"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateAddressRequest struct {
	Label         string `json:"label"`           // home, work, other (optional — default: home)
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	Street        string `json:"street"`
	City          string `json:"city"`
	Region        string `json:"region"`
	Country       string `json:"country"`
	IsDefault     bool   `json:"is_default"`
}

func (r CreateAddressRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RecipientName, validation.Required),
		validation.Field(&r.Phone, validation.Required),
		validation.Field(&r.Street, validation.Required),
		validation.Field(&r.City, validation.Required),
	)
}

type UpdateAddressRequest struct {
	Label         string `json:"label"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	Street        string `json:"street"`
	City          string `json:"city"`
	Region        string `json:"region"`
	Country       string `json:"country"`
	IsDefault     bool   `json:"is_default"`
}
