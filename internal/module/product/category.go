package product

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
)

type categoryModule struct {
	categoryStorage storage.CategoryStorage
}

func NewCategoryModule(cStorage storage.CategoryStorage) module.CategoryModule {
	return &categoryModule{
		categoryStorage: cStorage,
	}
}

func (m *categoryModule) CreateCategory(ctx context.Context, storeID int64, req dto.CreateCategoryRequest) (*dto.Category, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	dbCategory := &db.Category{
		StoreID: storeID,
		Name:    req.Name,
	}

	if err := m.categoryStorage.CreateCategory(ctx, dbCategory); err != nil {
		logger.Error("failed to create category", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, err
	}

	logger.Info("category created successfully", zap.Int64("category_id", dbCategory.ID), zap.Int64("store_id", storeID))

	return m.mapToDTO(dbCategory), nil
}

func (m *categoryModule) ListAllCategories(ctx context.Context, params dto.PaginationParams) (*dto.PaginatedResponse, error) {
	categories, total, err := m.categoryStorage.GetAllCategories(ctx, params.GetLimit(), params.GetOffset())
	if err != nil {
		return nil, err
	}

	dtoCategories := make([]dto.Category, len(categories))
	for i, c := range categories {
		dtoCategories[i] = *m.mapToDTO(&c)
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoCategories,
	}, nil
}

func (m *categoryModule) ListStoreCategories(ctx context.Context, storeID int64) ([]dto.Category, error) {
	categories, err := m.categoryStorage.GetCategoriesByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	dtoCategories := make([]dto.Category, len(categories))
	for i, c := range categories {
		dtoCategories[i] = *m.mapToDTO(&c)
	}

	return dtoCategories, nil
}

func (m *categoryModule) mapToDTO(c *db.Category) *dto.Category {
	return &dto.Category{
		ID:      c.ID,
		StoreID: c.StoreID,
		Name:    c.Name,
	}
}
