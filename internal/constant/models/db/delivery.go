package db

import (
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
)

type DeliveryAgent struct {
	BaseModel
	Username       string                        `gorm:"column:username;not null" json:"username"`
	FullName       string                        `gorm:"column:full_name;not null" json:"full_name"`
	Phone          string                        `gorm:"column:phone;not null" json:"phone"`
	UserID         *int64                        `gorm:"column:user_id" json:"user_id"`
	TelegramUserID *int64                        `gorm:"column:telegram_user_id" json:"telegram_user_id"`
	LoyaltyScore   int                           `gorm:"column:loyalty_score;default:0" json:"loyalty_score"`
	Status         constant.DeliveryAgentStatus  `gorm:"column:status;not null;default:'pending_invite'" json:"status"`
}

type StoreDeliveryLink struct {
	BaseModel
	StoreID           int64 `gorm:"column:store_id;not null" json:"store_id"`
	DeliveryAgentID   int64 `gorm:"column:delivery_agent_id;not null" json:"delivery_agent_id"`
	ShareEnabled      bool  `gorm:"column:share_enabled;default:false" json:"share_enabled"`
	ConnectedByUserID *int64 `gorm:"column:connected_by_user_id" json:"connected_by_user_id"`

	Store         Store         `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	DeliveryAgent DeliveryAgent `gorm:"foreignKey:DeliveryAgentID" json:"delivery_agent,omitempty"`
	Routes        []DeliveryRoute `gorm:"foreignKey:StoreDeliveryLinkID" json:"routes,omitempty"`
}

type DeliveryRoute struct {
	BaseModel
	StoreDeliveryLinkID int64  `gorm:"column:store_delivery_link_id;not null" json:"store_delivery_link_id"`
	Label               string `gorm:"column:label;not null" json:"label"`
	IsActive            bool   `gorm:"column:is_active;default:true" json:"is_active"`

	Locations []DeliveryRouteLocation `gorm:"foreignKey:DeliveryRouteID" json:"locations,omitempty"`
}

type DeliveryRouteLocation struct {
	BaseModel
	DeliveryRouteID  int64                        `gorm:"column:delivery_route_id;not null" json:"delivery_route_id"`
	LocationType     constant.DeliveryLocationType  `gorm:"column:location_type;not null" json:"location_type"`
	Label            string                       `gorm:"column:label" json:"label"`
	Country          string                       `gorm:"column:country;default:'Ethiopia'" json:"country"`
	Region           string                       `gorm:"column:region" json:"region"`
	City             string                       `gorm:"column:city" json:"city"`
	Street           string                       `gorm:"column:street" json:"street"`
	Landmark         string                       `gorm:"column:landmark" json:"landmark"`
	Notes            string                       `gorm:"column:notes" json:"notes"`
	UseStoreLocation bool                         `gorm:"column:use_store_location;default:false" json:"use_store_location"`
}

type DeliveryAgentShare struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerStoreID    int64     `gorm:"column:owner_store_id;not null" json:"owner_store_id"`
	DeliveryAgentID int64     `gorm:"column:delivery_agent_id;not null" json:"delivery_agent_id"`
	AdoptedStoreID  int64     `gorm:"column:adopted_store_id;not null" json:"adopted_store_id"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}

type DeliveryAreaPreset struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Region   string `gorm:"column:region" json:"region"`
	City     string `gorm:"column:city" json:"city"`
	Street   string `gorm:"column:street" json:"street"`
	IsActive bool   `gorm:"column:is_active;default:true" json:"is_active"`
}
