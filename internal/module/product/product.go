package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

type productModule struct {
	productStorage storage.ProductStorage
	storeStorage   storage.StoreStorage
	tele           platform.Telegram
	appURL         string
}

func NewProductModule(pStorage storage.ProductStorage, sStorage storage.StoreStorage, tele platform.Telegram, appURL string) module.ProductModule {
	return &productModule{
		productStorage: pStorage,
		storeStorage:   sStorage,
		tele:           tele,
		appURL:         appURL,
	}
}

func (m *productModule) incrementStoreViewsAsync(products []db.Product) {
	if len(products) == 0 {
		return
	}

	// Extract unique StoreIDs
	storeIDMap := make(map[int64]bool)
	for _, p := range products {
		if p.StoreID != nil {
			storeIDMap[*p.StoreID] = true
		}
	}

	if len(storeIDMap) == 0 {
		return
	}

	storeIDs := make([]int64, 0, len(storeIDMap))
	for id := range storeIDMap {
		storeIDs = append(storeIDs, id)
	}

	// Fire and forget in a goroutine
	go func() {
		// Use a fresh context for background task
		ctx := context.Background()
		if err := m.storeStorage.IncrementStoreViews(ctx, storeIDs); err != nil {
			logger.Error("failed to increment store views in background", zap.Error(err))
		}
	}()
}

func (m *productModule) CreateProduct(ctx context.Context, sellerID int64, storeID int64, req dto.CreateProductRequest) (*dto.Product, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	imagesBytes, _ := json.Marshal(req.Images)
	dbProduct := &db.Product{
		SellerID:    sellerID,
		StoreID:     &storeID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Images:      string(imagesBytes),
		Status:      constant.ProductStatusDraft, // Default to draft
		IsPosted:    req.IsPosted,
	}
	if err := m.productStorage.CreateProduct(ctx, dbProduct); err != nil {
		logger.Error("failed to create product", zap.Error(err), zap.Int64("store_id", storeID), zap.Int64("seller_id", sellerID))
		return nil, err
	}

	logger.Info("product created successfully", zap.Int64("product_id", dbProduct.ID), zap.Int64("store_id", storeID))

	return m.mapToDTO(dbProduct), nil
}

func (m *productModule) PostProduct(ctx context.Context, productID int64, storeID int64) (*dto.Product, error) {
	product, err := m.productStorage.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	product.IsPosted = true
	if err := m.productStorage.UpdateProduct(ctx, product); err != nil {
		logger.Error("failed to post product", zap.Error(err), zap.Int64("product_id", productID))
		return nil, err
	}

	logger.Info("product posted successfully", zap.Int64("product_id", productID))
	return m.mapToDTO(product), nil
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

	// Increment views asynchronously
	m.incrementStoreViewsAsync(products)

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

	// Increment views asynchronously
	m.incrementStoreViewsAsync(products)

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

	if req.Status != "" {
		// If transitioning to 'published' for the first time
		if product.Status != constant.ProductStatusPublished && req.Status == constant.ProductStatusPublished {
			go m.pushProductToTelegram(product)
		}
		product.Status = req.Status
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
		SellerID:    p.SellerID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		Images:      images,
		Status:      p.Status,
		IsPosted:    p.IsPosted,
		IsBoosted:   p.IsBoosted,
	}
}

func (m *productModule) pushProductToTelegram(p *db.Product) {
	if p.StoreID == nil {
		return
	}

	// 1. Get store to find chat ID
	ctx := context.Background()
	store, err := m.storeStorage.GetStoreByID(ctx, *p.StoreID)
	if err != nil || store.TelegramChatID == 0 {
		logger.Error("cannot push product: store not found or not linked", zap.Error(err))
		return
	}

	// 2. Construct message
	caption := fmt.Sprintf("🛍 *New Product Available!*\n\n*Name:* %s\n*Price:* %.2f ETB\n\n%s",
		p.Name, p.Price, p.Description)

	// 3. Construct "Order Now" button
	productURL := fmt.Sprintf("%s/product/%d", m.appURL, p.ID)
	selector := &telebot.ReplyMarkup{}
	btn := selector.URL("🛒 Order Now", productURL)
	selector.Inline(selector.Row(btn))

	// 4. Send with image if available
	var images []string
	_ = json.Unmarshal([]byte(p.Images), &images)

	bot := m.tele.GetBot()
	chat := &telebot.Chat{ID: store.TelegramChatID}

	if len(images) > 0 {
		photo := &telebot.Photo{
			File:    telebot.FromURL(images[0]),
			Caption: caption,
		}
		_, err = bot.Send(chat, photo, telebot.ModeMarkdown, selector)
	} else {
		_, err = bot.Send(chat, caption, telebot.ModeMarkdown, selector)
	}

	if err != nil {
		logger.Error("failed to push product to telegram", zap.Error(err), zap.Int64("product_id", p.ID))
	} else {
		logger.Info("product pushed to telegram successfully", zap.Int64("product_id", p.ID))
	}
}
