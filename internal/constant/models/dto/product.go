package dto

import (
	"time"

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
	Status      string   `json:"status"`
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
	Status      string   `json:"status"`
}

type ProductFilterParams struct {
	PaginationParams
	StoreID  int64    `form:"store_id"`
	Category string   `form:"category"`
	Query    string   `form:"query"`
	Title    string   `form:"title"`
	Status   string   `form:"status"`
	MinStock *int     `form:"min_stock"`
	MaxStock *int     `form:"max_stock"`
	MinPrice *float64 `form:"min_price"`
	MaxPrice *float64 `form:"max_price"`
}

type CreateProductInquiryRequest struct {
	Quantity int    `json:"quantity"`
	Note     string `json:"note"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

func (r CreateProductInquiryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Quantity, validation.Required, validation.Min(1)),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Phone, validation.Required),
	)
}

type ReviewProductInquiryRequest struct {
	ReviewNote string `json:"review_note"`
}

type ProductInquiry struct {
	ID         int64      `json:"id"`
	ProductID  int64      `json:"product_id"`
	StoreID    int64      `json:"store_id"`
	CustomerID int64      `json:"customer_id"`
	Status     string     `json:"status"`
	Quantity   int        `json:"quantity"`
	Note       string     `json:"note"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	ReviewedBy *int64     `json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	ReviewNote string     `json:"review_note,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Product    *Product   `json:"product,omitempty"`
}

type ProductInquiryFilterParams struct {
	PaginationParams
	ProductID int64  `form:"product_id"`
	Status    string `form:"status"`
}

func (r UpdateProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Price, validation.Min(0.0)),
		validation.Field(&r.Stock, validation.Min(0)),
	)
}

// ── Product Story DTOs ──────────────────────────────────────────────────────

// ProductStory is the full story response returned to the client.
type ProductStory struct {
	ID        int64    `json:"id"`
	StoreID   int64    `json:"store_id"`
	ProductID int64    `json:"product_id"`
	Caption   string   `json:"caption"`
	MediaURLs []string `json:"media_urls"`
	MediaType string   `json:"media_type"`
	StartsAt  string   `json:"starts_at"` // RFC3339
	EndsAt    string   `json:"ends_at"`   // RFC3339
	IsActive  bool     `json:"is_active"`
	Views     int64    `json:"views"`
	CreatedAt string   `json:"created_at"`
	// Product detail is populated only on single-story fetch (GetStory)
	Product *Product `json:"product,omitempty"`
}

// CreateProductStoryRequest is the payload for creating a new story ad.
type CreateProductStoryRequest struct {
	ProductID int64    `json:"product_id"`
	Caption   string   `json:"caption"`
	MediaURLs []string `json:"media_urls"`
	MediaType string   `json:"media_type"` // "image" | "video"
	StartsAt  string   `json:"starts_at"`  // RFC3339
	EndsAt    string   `json:"ends_at"`    // RFC3339
}

func (r CreateProductStoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProductID, validation.Required, validation.Min(int64(1))),
		validation.Field(&r.MediaURLs, validation.Required, validation.Length(1, 10)),
		validation.Field(&r.MediaType, validation.Required, validation.In("image", "video")),
		validation.Field(&r.StartsAt, validation.Required),
		validation.Field(&r.EndsAt, validation.Required),
	)
}

// UpdateProductStoryRequest supports partial story updates.
type UpdateProductStoryRequest struct {
	Caption   string   `json:"caption"`
	MediaURLs []string `json:"media_urls"`
	MediaType string   `json:"media_type"`
	StartsAt  string   `json:"starts_at"`
	EndsAt    string   `json:"ends_at"`
	IsActive  *bool    `json:"is_active"`
}

func (r UpdateProductStoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MediaType, validation.In("image", "video")),
		validation.Field(&r.MediaURLs, validation.Length(0, 10)),
	)
}

// ProductStoryFilterParams is used to list stories (merchant-scoped or public).
type ProductStoryFilterParams struct {
	PaginationParams
	StoreID   int64  `form:"-"`          // injected server-side
	StoreName string `form:"store_name"` // optional ILIKE search on store name
	ProductID *int64 `form:"product_id"` // optional client filter
	IsActive  *bool  `form:"is_active"`
	Type      string `form:"type"`    // video or image
	Search    string `form:"search"`  // search by caption/title
	SortBy    string `form:"sort_by"` // newest, popular
}

// ── Favorite DTOs ──────────────────────────────────────────────────────────

type FavoriteResponse struct {
	ID        int64    `json:"id"`
	UserID    int64    `json:"user_id"`
	ProductID int64    `json:"product_id"`
	CreatedAt string   `json:"created_at"`
	Product   *Product `json:"product,omitempty"`
}
