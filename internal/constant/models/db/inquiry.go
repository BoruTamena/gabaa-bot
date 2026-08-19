package db

import "time"

type ProductInquiry struct {
	BaseModel
	ProductID  int64      `gorm:"column:product_id;not null;index" json:"product_id"`
	StoreID    int64      `gorm:"column:store_id;not null;index" json:"store_id"`
	CustomerID int64      `gorm:"column:customer_id;not null;index" json:"customer_id"`
	Status     string     `gorm:"column:status;not null;default:pending;index" json:"status"`
	Quantity   int        `gorm:"column:quantity;not null" json:"quantity"`
	Note       string     `gorm:"column:note" json:"note"`
	Name       string     `gorm:"column:name;not null" json:"name"`
	Phone      string     `gorm:"column:phone;not null" json:"phone"`
	ReviewedBy *int64     `gorm:"column:reviewed_by" json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time `gorm:"column:reviewed_at" json:"reviewed_at,omitempty"`
	ReviewNote string     `gorm:"column:review_note" json:"review_note,omitempty"`
	Product    Product    `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Store      Store      `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	Customer   User       `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Reviewer   *User      `gorm:"foreignKey:ReviewedBy;references:ID" json:"reviewer,omitempty"`
}
