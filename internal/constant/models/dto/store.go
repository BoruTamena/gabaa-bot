package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Store struct {
	ID             int64  `json:"id"`
	SellerID       int64  `json:"seller_id"`
	TelegramChatID int64  `json:"telegram_chat_id"`
	TelegramChatTitle string `json:"telegram_chat_title"`
	Status         string `json:"status"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	LogoImage      string `json:"logo_image"`
	CoverImage     string `json:"cover_image"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Location       string `json:"location"`
}

type CreateStoreRequest struct {
	TelegramChatID int64  `json:"telegram_chat_id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	LogoImage      string `json:"logo_image"`
	CoverImage     string `json:"cover_image"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Location       string `json:"location"`
}

func (r CreateStoreRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Category, validation.Required),
		validation.Field(&r.Phone, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
	)
}

type UpdateStoreRequest struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	LogoImage   string `json:"logo_image"`
	CoverImage  string `json:"cover_image"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Location    string `json:"location"`
}

func (r UpdateStoreRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, is.Email),
	)
}
