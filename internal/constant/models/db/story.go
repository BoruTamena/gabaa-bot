package db

import "time"

// ProductStory represents a story ad created by a store owner for a specific product.
// Stories are only publicly visible within the [StartsAt, EndsAt] date range.
type ProductStory struct {
	BaseModel
	StoreID   int64     `gorm:"column:store_id;not null;index"                json:"store_id"`
	ProductID int64     `gorm:"column:product_id;not null;index"               json:"product_id"`
	Caption   string    `gorm:"column:caption"                                 json:"caption"`
	MediaURLs string    `gorm:"column:media_urls;type:jsonb;not null"          json:"media_urls"` // marshalled []string
	MediaType string    `gorm:"column:media_type;not null;default:image"       json:"media_type"` // "image" | "video"
	StartsAt  time.Time `gorm:"column:starts_at;not null"                      json:"starts_at"`
	EndsAt    time.Time `gorm:"column:ends_at;not null"                        json:"ends_at"`
	IsActive  bool      `gorm:"column:is_active;default:true"                  json:"is_active"`
	Views     int64     `gorm:"column:views;default:0"                         json:"views"`

	// Associations (loaded on demand)
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
