package dto

import "time"

// AnalyticsFilterParams holds optional date range filters
type AnalyticsFilterParams struct {
	From *time.Time `form:"from" time_format:"2006-01-02T15:04:05Z07:00"`
	To   *time.Time `form:"to" time_format:"2006-01-02T15:04:05Z07:00"`
}

// PeriodRevenue represents sales aggregated by date
type PeriodRevenue struct {
	Period  string  `json:"period"` // e.g. "2026-07-08"
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

// TopProduct represents a simplified product structure for analytics lists
type TopProduct struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Revenue     float64 `json:"revenue"`
	UnitsSold   int64   `json:"units_sold"`
}

// SalesAnalytics response DTO
type SalesAnalytics struct {
	TotalRevenue       float64         `json:"total_revenue"`
	RevenueChangePct   float64         `json:"revenue_change_pct"`
	TotalOrders        int64           `json:"total_orders"`
	AverageOrderValue  float64         `json:"average_order_value"`
	RevenueByPeriod    []PeriodRevenue `json:"revenue_by_period"`
	TopSellingProducts []TopProduct    `json:"top_selling_products"`
}

// OrdersByStatus represents order status distribution
type OrdersByStatus struct {
	Status     string  `json:"status"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// OrderAnalytics response DTO
type OrderAnalytics struct {
	TotalOrders       int64            `json:"total_orders"`
	OrdersByStatus    []OrdersByStatus `json:"orders_by_status"`
	AverageOrderValue float64          `json:"average_order_value"`
	RecentOrders      int64            `json:"recent_orders"` // last 7 days
	CancellationRate  float64          `json:"cancellation_rate_pct"`
}

// ProductByStatus represents product status distribution
type ProductByStatus struct {
	Status string `json:"status"` // draft, published, archived
	Count  int64  `json:"count"`
}

// TopViewedProduct represents a product with its view counts
type TopViewedProduct struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	Views       int64  `json:"views"`
	Category    string `json:"category"`
}

// ProductAnalytics response DTO
type ProductAnalytics struct {
	TotalProducts    int64             `json:"total_products"`
	ProductsByStatus []ProductByStatus `json:"products_by_status"`
	LowStockCount    int64             `json:"low_stock_count"`    // stock <= 5
	OutOfStockCount  int64             `json:"out_of_stock_count"` // stock == 0
	TopViewedProducts []TopViewedProduct `json:"top_viewed_products"`
}

// TopStory represents a story with its view counts
type TopStory struct {
	StoryID   int64  `json:"story_id"`
	ProductID int64  `json:"product_id"`
	Caption   string `json:"caption"`
	Views     int64  `json:"views"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
}

// StoryAnalytics response DTO
type StoryAnalytics struct {
	TotalStories   int64      `json:"total_stories"`
	ActiveStories  int64      `json:"active_stories"`
	ExpiredStories int64      `json:"expired_stories"`
	TotalViews     int64      `json:"total_views"`
	TopStories     []TopStory `json:"top_stories"`
}
