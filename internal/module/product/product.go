package product

import (
	"context"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type productModule struct {
	productStorage storage.ProductStorage
}

func NewProductModule(pStorage storage.ProductStorage) module.ProductModule {
	return &productModule{
		productStorage: pStorage,
	}
}

func (m *productModule) CreateProduct(ctx context.Context, product *dto.Product) error {
	dbProduct := &db.Product{
		StoreID:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Images:      product.Images,
	}
	if err := m.productStorage.CreateProduct(ctx, dbProduct); err != nil {
		return err
	}
	product.ID = dbProduct.ID
	return nil
}

func (m *productModule) GetProduct(ctx context.Context, id int64) (*dto.Product, error) {
	product, err := m.productStorage.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.Product{
		ID:          product.ID,
		StoreID:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Images:      product.Images,
	}, nil
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
		dtoProducts[i] = dto.Product{
			ID:          p.ID,
			StoreID:     p.StoreID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Images:      p.Images,
		}
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoProducts,
	}, nil
}

func (m *productModule) UpdateProduct(ctx context.Context, product *dto.Product) error {
	dbProduct := &db.Product{
		StoreID:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Images:      product.Images,
	}
	dbProduct.ID = product.ID
	return m.productStorage.UpdateProduct(ctx, dbProduct)
}

func (m *productModule) DeleteProduct(ctx context.Context, id int64) error {
	return m.productStorage.DeleteProduct(ctx, id)
}


