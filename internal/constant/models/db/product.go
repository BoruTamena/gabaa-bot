package db

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string  `gorm:"title" json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
