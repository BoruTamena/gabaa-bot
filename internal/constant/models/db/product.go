package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {

	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil

}

type Product struct {
	BaseModel
	Title       string  `gorm:"title" json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
