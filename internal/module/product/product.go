package product

import (
	"context"
	"encoding/json"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
)

type productModule struct {
	productStorage storage.ProductStorage
}

func NewProductModule(pStorage storage.ProductStorage) module.ProductModule {
	return &productModule{
		productStorage: pStorage,
	}
}

func (m *productModule) CreateProduct(ctx context.Context, storeID int64, req dto.CreateProductRequest) (*dto.Product, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	imagesBytes, _ := json.Marshal(req.Images)
	dbProduct := &db.Product{
		StoreID:     storeID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Images:      string(imagesBytes),
	}
	if err := m.productStorage.CreateProduct(ctx, dbProduct); err != nil {
		logger.Error("failed to create product", zap.Error(err), zap.Int64("store_id", storeID))
		return nil, err
	}

	logger.Info("product created successfully", zap.Int64("product_id", dbProduct.ID), zap.Int64("store_id", storeID))

	return m.mapToDTO(dbProduct), nil
}

func (m *productModule) GetProduct(ctx context.Context, id int64) (*dto.Product, error) {
	product, err := m.productStorage.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.mapToDTO(product), nil
}

func (m *productModule) ListProducts(ctx context.Context, storeID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error) {
	limit := params.GetLimit()
	offset := params.GetOffset()

	products, err := m.productStorage.GetProductsByStoreID(ctx, storeID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := m.productStorage.GetProductsTotal(ctx, storeID)
	if err != nil {
		return nil, err
	}

	dtoProducts := make([]dto.Product, len(products))
	for i, p := range products {
		dtoProducts[i] = *m.mapToDTO(&p)
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoProducts,
	}, nil
}

func (m *productModule) ListAllProducts(ctx context.Context, filter dto.ProductFilterParams) (*dto.PaginatedResponse, error) {
	products, total, err := m.productStorage.ListAllProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	dtoProducts := make([]dto.Product, len(products))
	for i, p := range products {
		dtoProducts[i] = *m.mapToDTO(&p)
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoProducts,
	}, nil
}

func (m *productModule) UpdateProduct(ctx context.Context, id int64, req dto.UpdateProductRequest) (*dto.Product, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	product, err := m.productStorage.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Stock != 0 {
		product.Stock = req.Stock
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if len(req.Images) > 0 {
		imagesBytes, _ := json.Marshal(req.Images)
		product.Images = string(imagesBytes)
	}

	if err := m.productStorage.UpdateProduct(ctx, product); err != nil {
		logger.Error("failed to update product", zap.Error(err), zap.Int64("product_id", id))
		return nil, err
	}

	logger.Info("product updated successfully", zap.Int64("product_id", id))

	return m.mapToDTO(product), nil
}

func (m *productModule) DeleteProduct(ctx context.Context, id int64) error {
	err := m.productStorage.DeleteProduct(ctx, id)
	if err != nil {
		logger.Error("failed to delete product", zap.Error(err), zap.Int64("product_id", id))
	} else {
		logger.Info("product deleted successfully", zap.Int64("product_id", id))
	}
	return err
}

func (m *productModule) mapToDTO(p *db.Product) *dto.Product {
	var images []string
	if p.Images != "" {
		_ = json.Unmarshal([]byte(p.Images), &images)
	}
	if images == nil {
		images = []string{}
	}

	return &dto.Product{
		ID:          p.ID,
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		Images:      images,
	}
}
