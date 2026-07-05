package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Store struct {
	ID                 int64  `json:"id"`
	SellerID           int64  `json:"seller_id"`
	TelegramChatID     int64  `json:"telegram_chat_id"`
	TelegramChatTitle  string `json:"telegram_chat_title"`
	Status             string `json:"status"`
	VerificationStatus string `json:"verificationStatus"`
	Name               string `json:"name"`
	Category           string `json:"category"`
	Description        string `json:"description"`
	LogoImage          string `json:"logo_image"`
	CoverImage         string `json:"cover_image"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	Location           string `json:"location"`
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

type SubmitStoreKYCRequest struct {
	TINNumber                  string `json:"tinNumber"`
	BusinessRegistrationNumber string `json:"businessRegistrationNumber"`
	TINCertificateURL          string `json:"tinCertificateUrl"`
	BusinessLicenseURL         string `json:"businessLicenseUrl"`
}

func (r SubmitStoreKYCRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TINNumber, validation.Required),
		validation.Field(&r.BusinessRegistrationNumber, validation.Required),
		validation.Field(&r.TINCertificateURL, validation.Required, is.URL),
		validation.Field(&r.BusinessLicenseURL, validation.Required, is.URL),
	)
}

type StoreKYCResponse struct {
	StoreID                    int64      `json:"storeId"`
	StoreName                  string     `json:"storeName,omitempty"`
	VerificationStatus         string     `json:"verificationStatus"`
	TINNumber                  string     `json:"tinNumber,omitempty"`
	BusinessRegistrationNumber string     `json:"businessRegistrationNumber,omitempty"`
	TINCertificateURL          string     `json:"tinCertificateUrl,omitempty"`
	BusinessLicenseURL         string     `json:"businessLicenseUrl,omitempty"`
	ReviewNote                 string     `json:"reviewNote,omitempty"`
	SubmittedAt                *time.Time `json:"submittedAt,omitempty"`
	ReviewedAt                 *time.Time `json:"reviewedAt,omitempty"`
}

type RejectStoreKYCRequest struct {
	ReviewNote string `json:"reviewNote"`
}
