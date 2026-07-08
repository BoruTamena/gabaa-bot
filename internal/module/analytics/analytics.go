package analytics

import (
	"context"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
)

type analyticsModule struct {
	analyticsStorage storage.AnalyticsStorage
}

func NewAnalyticsModule(aStorage storage.AnalyticsStorage) module.AnalyticsModule {
	return &analyticsModule{
		analyticsStorage: aStorage,
	}
}

// setDefaultDates assigns default from/to values if they are omitted (defaults to last 30 days)
func setDefaultDates(filter *dto.AnalyticsFilterParams) {
	now := time.Now()
	if filter.To == nil {
		filter.To = &now
	}
	if filter.From == nil {
		defaultFrom := filter.To.AddDate(0, 0, -30)
		filter.From = &defaultFrom
	}
}

func (m *analyticsModule) GetSalesAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.SalesAnalytics, error) {
	setDefaultDates(&filter)
	logger.Info("Fetching sales analytics", zap.Int64("store_id", storeID), zap.Time("from", *filter.From), zap.Time("to", *filter.To))
	
	resp, err := m.analyticsStorage.GetSalesAnalytics(ctx, storeID, filter)
	if err != nil {
		logger.Error("Failed to fetch sales analytics in module", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, err
	}
	return resp, nil
}

func (m *analyticsModule) GetOrderAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.OrderAnalytics, error) {
	setDefaultDates(&filter)
	logger.Info("Fetching order analytics", zap.Int64("store_id", storeID), zap.Time("from", *filter.From), zap.Time("to", *filter.To))

	resp, err := m.analyticsStorage.GetOrderAnalytics(ctx, storeID, filter)
	if err != nil {
		logger.Error("Failed to fetch order analytics in module", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, err
	}
	return resp, nil
}

func (m *analyticsModule) GetProductAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.ProductAnalytics, error) {
	setDefaultDates(&filter)
	logger.Info("Fetching product analytics", zap.Int64("store_id", storeID))

	resp, err := m.analyticsStorage.GetProductAnalytics(ctx, storeID, filter)
	if err != nil {
		logger.Error("Failed to fetch product analytics in module", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, err
	}
	return resp, nil
}

func (m *analyticsModule) GetStoryAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.StoryAnalytics, error) {
	setDefaultDates(&filter)
	logger.Info("Fetching story analytics", zap.Int64("store_id", storeID))

	resp, err := m.analyticsStorage.GetStoryAnalytics(ctx, storeID, filter)
	if err != nil {
		logger.Error("Failed to fetch story analytics in module", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, err
	}
	return resp, nil
}
