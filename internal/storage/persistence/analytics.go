package persistence

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"gorm.io/gorm"
)

type analyticsPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewAnalyticsPersistence(db *gorm.DB, logger platform.Logger) storage.AnalyticsStorage {
	return &analyticsPersistence{db: db, logger: logger}
}

func (p *analyticsPersistence) GetSalesAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.SalesAnalytics, error) {
	from := *filter.From
	to := *filter.To
	duration := to.Sub(from)
	prevFrom := from.Add(-duration)
	prevTo := from

	// 1. Current Period Total Revenue & Orders
	var currentStats struct {
		Revenue float64
		Orders  int64
	}
	err := p.db.WithContext(ctx).
		Table("orders").
		Select("COALESCE(SUM(total_price), 0) as revenue, COUNT(*) as orders").
		Where("store_id = ? AND status != ? AND created_at >= ? AND created_at <= ?", storeID, "cancelled", from, to).
		Scan(&currentStats).Error
	if err != nil {
		p.logger.Error("Failed to fetch current period sales analytics", "error", err, "storeID", storeID)
		return nil, err
	}

	// 2. Previous Period Total Revenue (for percentage change)
	var prevStats struct {
		Revenue float64
	}
	err = p.db.WithContext(ctx).
		Table("orders").
		Select("COALESCE(SUM(total_price), 0) as revenue").
		Where("store_id = ? AND status != ? AND created_at >= ? AND created_at <= ?", storeID, "cancelled", prevFrom, prevTo).
		Scan(&prevStats).Error
	if err != nil {
		p.logger.Error("Failed to fetch previous period sales analytics", "error", err, "storeID", storeID)
		return nil, err
	}

	// Calculate percentage change
	var revenueChangePct float64
	if prevStats.Revenue > 0 {
		revenueChangePct = ((currentStats.Revenue - prevStats.Revenue) / prevStats.Revenue) * 100
	} else if currentStats.Revenue > 0 {
		revenueChangePct = 100.0 // 100% increase if there was no revenue before
	}

	// Calculate average order value
	var avgOrderValue float64
	if currentStats.Orders > 0 {
		avgOrderValue = currentStats.Revenue / float64(currentStats.Orders)
	}

	// 3. Period Revenue (Daily)
	var periodRevenues []dto.PeriodRevenue
	err = p.db.WithContext(ctx).
		Table("orders").
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as period, COALESCE(SUM(total_price), 0) as revenue, COUNT(*) as orders").
		Where("store_id = ? AND status != ? AND created_at >= ? AND created_at <= ?", storeID, "cancelled", from, to).
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("period ASC").
		Scan(&periodRevenues).Error
	if err != nil {
		p.logger.Error("Failed to fetch period revenues", "error", err, "storeID", storeID)
		return nil, err
	}

	// If empty, return an empty slice instead of nil
	if periodRevenues == nil {
		periodRevenues = []dto.PeriodRevenue{}
	}

	// 4. Top Selling Products
	var topProducts []dto.TopProduct
	err = p.db.WithContext(ctx).
		Table("order_items oi").
		Select("oi.product_id, p.name as product_name, COALESCE(SUM(oi.quantity * oi.price), 0) as revenue, COALESCE(SUM(oi.quantity), 0) as units_sold").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Joins("JOIN products p ON p.id = oi.product_id").
		Where("o.store_id = ? AND o.status != ? AND o.created_at >= ? AND o.created_at <= ?", storeID, "cancelled", from, to).
		Group("oi.product_id, p.name").
		Order("units_sold DESC, revenue DESC").
		Limit(5).
		Scan(&topProducts).Error
	if err != nil {
		p.logger.Error("Failed to fetch top selling products", "error", err, "storeID", storeID)
		return nil, err
	}

	if topProducts == nil {
		topProducts = []dto.TopProduct{}
	}

	return &dto.SalesAnalytics{
		TotalRevenue:       currentStats.Revenue,
		RevenueChangePct:   math.Round(revenueChangePct*100) / 100, // round to 2 decimal places
		TotalOrders:        currentStats.Orders,
		AverageOrderValue:  math.Round(avgOrderValue*100) / 100,
		RevenueByPeriod:    periodRevenues,
		TopSellingProducts: topProducts,
	}, nil
}

func (p *analyticsPersistence) GetOrderAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.OrderAnalytics, error) {
	from := *filter.From
	to := *filter.To

	// 1. Total Orders & Average Order Value (including cancelled/pending for total metrics)
	var orderStats struct {
		TotalOrders int64
		TotalRevenue float64
	}
	err := p.db.WithContext(ctx).
		Table("orders").
		Select("COUNT(*) as total_orders, COALESCE(SUM(total_price), 0) as total_revenue").
		Where("store_id = ? AND created_at >= ? AND created_at <= ?", storeID, from, to).
		Scan(&orderStats).Error
	if err != nil {
		p.logger.Error("Failed to fetch order stats", "error", err, "storeID", storeID)
		return nil, err
	}

	// Average order value
	var avgOrderValue float64
	if orderStats.TotalOrders > 0 {
		avgOrderValue = orderStats.TotalRevenue / float64(orderStats.TotalOrders)
	}

	// 2. Orders by Status
	type statusCount struct {
		Status string
		Count  int64
	}
	var statuses []statusCount
	err = p.db.WithContext(ctx).
		Table("orders").
		Select("status, COUNT(*) as count").
		Where("store_id = ? AND created_at >= ? AND created_at <= ?", storeID, from, to).
		Group("status").
		Scan(&statuses).Error
	if err != nil {
		p.logger.Error("Failed to fetch orders by status", "error", err, "storeID", storeID)
		return nil, err
	}

	ordersByStatus := make([]dto.OrdersByStatus, 0, len(statuses))
	var cancelledCount int64
	for _, s := range statuses {
		pct := 0.0
		if orderStats.TotalOrders > 0 {
			pct = (float64(s.Count) / float64(orderStats.TotalOrders)) * 100
		}
		ordersByStatus = append(ordersByStatus, dto.OrdersByStatus{
			Status:     s.Status,
			Count:      s.Count,
			Percentage: math.Round(pct*100) / 100,
		})
		if s.Status == "cancelled" {
			cancelledCount = s.Count
		}
	}

	// 3. Recent Orders (last 7 days)
	var recentOrders int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	err = p.db.WithContext(ctx).
		Table("orders").
		Where("store_id = ? AND created_at >= ?", storeID, sevenDaysAgo).
		Count(&recentOrders).Error
	if err != nil {
		p.logger.Error("Failed to fetch recent orders count", "error", err, "storeID", storeID)
		return nil, err
	}

	// 4. Cancellation Rate
	var cancellationRate float64
	if orderStats.TotalOrders > 0 {
		cancellationRate = (float64(cancelledCount) / float64(orderStats.TotalOrders)) * 100
	}

	return &dto.OrderAnalytics{
		TotalOrders:       orderStats.TotalOrders,
		OrdersByStatus:    ordersByStatus,
		AverageOrderValue: math.Round(avgOrderValue*100) / 100,
		RecentOrders:      recentOrders,
		CancellationRate:  math.Round(cancellationRate*100) / 100,
	}, nil
}

func (p *analyticsPersistence) GetProductAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.ProductAnalytics, error) {
	// 1. Total Products
	var totalProducts int64
	err := p.db.WithContext(ctx).
		Table("products").
		Where("store_id = ?", storeID).
		Count(&totalProducts).Error
	if err != nil {
		p.logger.Error("Failed to fetch total products count", "error", err, "storeID", storeID)
		return nil, err
	}

	// 2. Products by Status
	type statusCount struct {
		Status string
		Count  int64
	}
	var statuses []statusCount
	err = p.db.WithContext(ctx).
		Table("products").
		Select("status, COUNT(*) as count").
		Where("store_id = ?", storeID).
		Group("status").
		Scan(&statuses).Error
	if err != nil {
		p.logger.Error("Failed to fetch products by status", "error", err, "storeID", storeID)
		return nil, err
	}

	productsByStatus := make([]dto.ProductByStatus, 0, len(statuses))
	for _, s := range statuses {
		productsByStatus = append(productsByStatus, dto.ProductByStatus{
			Status: s.Status,
			Count:  s.Count,
		})
	}

	// 3. Low Stock & Out of Stock Count (only for active/published products)
	var lowStockCount int64
	err = p.db.WithContext(ctx).
		Table("products").
		Where("store_id = ? AND status = ? AND stock <= ? AND stock > ?", storeID, "published", 5, 0).
		Count(&lowStockCount).Error
	if err != nil {
		p.logger.Error("Failed to fetch low stock count", "error", err, "storeID", storeID)
		return nil, err
	}

	var outOfStockCount int64
	err = p.db.WithContext(ctx).
		Table("products").
		Where("store_id = ? AND status = ? AND stock = ?", storeID, "published", 0).
		Count(&outOfStockCount).Error
	if err != nil {
		p.logger.Error("Failed to fetch out of stock count", "error", err, "storeID", storeID)
		return nil, err
	}

	// 4. Top Viewed Products (Proxied via order items purchase counts)
	var topViewed []dto.TopViewedProduct
	err = p.db.WithContext(ctx).
		Table("products p").
		Select("p.id as product_id, p.name as product_name, COALESCE(SUM(oi.quantity), 0) as views, p.category as category").
		Joins("LEFT JOIN order_items oi ON oi.product_id = p.id").
		Joins("LEFT JOIN orders o ON o.id = oi.order_id AND o.status != 'cancelled'").
		Where("p.store_id = ?", storeID).
		Group("p.id, p.name, p.category").
		Order("views DESC").
		Limit(5).
		Scan(&topViewed).Error
	if err != nil {
		p.logger.Error("Failed to fetch top viewed products", "error", err, "storeID", storeID)
		return nil, err
	}

	if topViewed == nil {
		topViewed = []dto.TopViewedProduct{}
	}

	return &dto.ProductAnalytics{
		TotalProducts:     totalProducts,
		ProductsByStatus:  productsByStatus,
		LowStockCount:     lowStockCount,
		OutOfStockCount:   outOfStockCount,
		TopViewedProducts: topViewed,
	}, nil
}

func (p *analyticsPersistence) GetStoryAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.StoryAnalytics, error) {
	now := time.Now()

	// 1. Total Stories
	var totalStories int64
	err := p.db.WithContext(ctx).
		Table("product_stories").
		Where("store_id = ?", storeID).
		Count(&totalStories).Error
	if err != nil {
		p.logger.Error("Failed to fetch total stories count", "error", err, "storeID", storeID)
		return nil, err
	}

	// 2. Active Stories
	var activeStories int64
	err = p.db.WithContext(ctx).
		Table("product_stories").
		Where("store_id = ? AND is_active = ? AND starts_at <= ? AND ends_at >= ?", storeID, true, now, now).
		Count(&activeStories).Error
	if err != nil {
		p.logger.Error("Failed to fetch active stories count", "error", err, "storeID", storeID)
		return nil, err
	}

	// 3. Expired Stories
	var expiredStories int64
	err = p.db.WithContext(ctx).
		Table("product_stories").
		Where("store_id = ? AND (is_active = ? OR ends_at < ?)", storeID, false, now).
		Count(&expiredStories).Error
	if err != nil {
		p.logger.Error("Failed to fetch expired stories count", "error", err, "storeID", storeID)
		return nil, err
	}

	// 4. Total Views
	var totalViews sql.NullInt64
	err = p.db.WithContext(ctx).
		Table("product_stories").
		Select("SUM(views)").
		Where("store_id = ?", storeID).
		Row().
		Scan(&totalViews)
	if err != nil {
		p.logger.Error("Failed to fetch total story views", "error", err, "storeID", storeID)
		return nil, err
	}

	// 5. Top Stories (by views)
	type storyDB struct {
		ID        int64
		ProductID int64
		Caption   string
		Views     int64
		StartsAt  time.Time
		EndsAt    time.Time
	}
	var dbStories []storyDB
	err = p.db.WithContext(ctx).
		Table("product_stories").
		Select("id, product_id, caption, views, starts_at, ends_at").
		Where("store_id = ?", storeID).
		Order("views DESC").
		Limit(5).
		Scan(&dbStories).Error
	if err != nil {
		p.logger.Error("Failed to fetch top stories", "error", err, "storeID", storeID)
		return nil, err
	}

	topStories := make([]dto.TopStory, 0, len(dbStories))
	for _, s := range dbStories {
		topStories = append(topStories, dto.TopStory{
			StoryID:   s.ID,
			ProductID: s.ProductID,
			Caption:   s.Caption,
			Views:     s.Views,
			StartsAt:  s.StartsAt.Format(time.RFC3339),
			EndsAt:    s.EndsAt.Format(time.RFC3339),
		})
	}

	return &dto.StoryAnalytics{
		TotalStories:   totalStories,
		ActiveStories:  activeStories,
		ExpiredStories: expiredStories,
		TotalViews:     totalViews.Int64,
		TopStories:     topStories,
	}, nil
}
