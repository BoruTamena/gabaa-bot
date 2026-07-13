package dto

type DeliveryRouteLocationInput struct {
	Label            string `json:"label"`
	Country          string `json:"country"`
	Region           string `json:"region"`
	City             string `json:"city"`
	Street           string `json:"street"`
	Landmark         string `json:"landmark"`
	Notes            string `json:"notes"`
	UseStoreLocation bool   `json:"use_store_location"`
}

type DeliveryRouteInput struct {
	Label             string                       `json:"label"`
	PickupLocations   []DeliveryRouteLocationInput `json:"pickup_locations"`
	DeliveryLocations []DeliveryRouteLocationInput `json:"delivery_locations"`
}

type ConnectDeliveryAgentRequest struct {
	Username     string               `json:"username"`
	FullName     string               `json:"full_name"`
	Phone        string               `json:"phone"`
	ShareEnabled bool                 `json:"share_enabled"`
	Routes       []DeliveryRouteInput `json:"routes"`
}

type UpdateDeliveryAgentRequest struct {
	FullName     *string              `json:"full_name"`
	Phone        *string              `json:"phone"`
	ShareEnabled *bool                `json:"share_enabled"`
	Routes       []DeliveryRouteInput `json:"routes"`
}

type DeliveryRouteLocationResponse struct {
	ID               int64  `json:"id"`
	LocationType     string `json:"location_type"`
	Label            string `json:"label"`
	Country          string `json:"country"`
	Region           string `json:"region"`
	City             string `json:"city"`
	Street           string `json:"street"`
	Landmark         string `json:"landmark"`
	Notes            string `json:"notes"`
	UseStoreLocation bool   `json:"use_store_location"`
}

type DeliveryRouteResponse struct {
	ID                int64                           `json:"id"`
	Label             string                          `json:"label"`
	IsActive          bool                            `json:"is_active"`
	PickupLocations   []DeliveryRouteLocationResponse `json:"pickup_locations"`
	DeliveryLocations []DeliveryRouteLocationResponse `json:"delivery_locations"`
}

type DeliveryAgentResponse struct {
	ID           int64                   `json:"id"`
	Username     string                  `json:"username"`
	FullName     string                  `json:"full_name"`
	Phone        string                  `json:"phone"`
	Status       string                  `json:"status"`
	LoyaltyScore int                     `json:"loyalty_score"`
	ShareEnabled bool                    `json:"share_enabled"`
	Routes       []DeliveryRouteResponse `json:"routes"`
}

type DeliveryAreaPreset struct {
	ID     int64  `json:"id"`
	Region string `json:"region"`
	City   string `json:"city"`
	Street string `json:"street"`
}

type DeliverySuggestion struct {
	AgentID      int64  `json:"agent_id"`
	RouteID      int64  `json:"route_id"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	LoyaltyScore int    `json:"loyalty_score"`
	Score        int    `json:"score"`
	MatchSummary string `json:"match_summary"`
	Suggested    bool   `json:"suggested"`
}

type ShipOrderRequest struct {
	Status           string `json:"status"`
	DeliveryAgentID  *int64 `json:"delivery_agent_id"`
	DeliveryRouteID  *int64 `json:"delivery_route_id"`
}

type DeliveryProfileResponse struct {
	ID           int64                   `json:"id"`
	Username     string                  `json:"username"`
	FullName     string                  `json:"full_name"`
	Phone        string                  `json:"phone"`
	LoyaltyScore int                     `json:"loyalty_score"`
	Status       string                  `json:"status"`
	Routes       []DeliveryRouteResponse `json:"routes"`
}

type DeliveryOrderResponse struct {
	Order
	PickupLocation string `json:"pickup_location"`
	StoreName      string `json:"store_name"`
}

type DeliveryOrderStatusRequest struct {
	Status string `json:"status"`
}

type AddDeliveryRouteRequest struct {
	Label             string                       `json:"label"`
	PickupLocations   []DeliveryRouteLocationInput `json:"pickup_locations"`
	DeliveryLocations []DeliveryRouteLocationInput `json:"delivery_locations"`
}

type UpdateDeliveryRouteRequest struct {
	Label             *string                      `json:"label"`
	IsActive          *bool                        `json:"is_active"`
	PickupLocations   []DeliveryRouteLocationInput `json:"pickup_locations"`
	DeliveryLocations []DeliveryRouteLocationInput `json:"delivery_locations"`
}
