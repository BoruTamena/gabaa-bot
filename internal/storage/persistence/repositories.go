package persistence

import (
	"context"
	"fmt"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"gorm.io/gorm"
)

type persistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewPersistence(db *gorm.DB, logger platform.Logger) storage.UserStorage {
	return &persistence{db: db, logger: logger}
}

// UserStorage
func (p *persistence) CreateUser(ctx context.Context, user *db.User) error {
	err := p.db.WithContext(ctx).Create(user).Error
	if err != nil {
		p.logger.Error("Failed to create user", "error", err)
	}
	return err
}

func (p *persistence) GetUserByTelegramID(ctx context.Context, telegramID int64) (*db.User, error) {
	var user db.User
	err := p.db.WithContext(ctx).Where("telegram_user_id = ?", telegramID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		p.logger.Error("Failed to get user by telegram ID", "error", err, "telegramID", telegramID)
	}
	return &user, err
}

func (p *persistence) GetUserByID(ctx context.Context, id int64) (*db.User, error) {
	var user db.User
	err := p.db.WithContext(ctx).First(&user, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		p.logger.Error("Failed to get user by ID", "error", err, "userID", id)
	}
	return &user, err
}

func (p *persistence) UpdateUser(ctx context.Context, user *db.User) error {
	err := p.db.WithContext(ctx).Save(user).Error
	if err != nil {
		p.logger.Error("Failed to update user", "error", err, "userID", user.ID)
	}
	return err
}

// StoreStorage
type storePersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewStorePersistence(db *gorm.DB, logger platform.Logger) storage.StoreStorage {
	return &storePersistence{db: db, logger: logger}
}

func (p *storePersistence) CreateStore(ctx context.Context, store *db.Store) error {
	err := p.db.WithContext(ctx).Create(store).Error
	if err != nil {
		p.logger.Error("Failed to create store", "error", err)
	}
	return err
}

func (p *storePersistence) GetStoreByID(ctx context.Context, id int64) (*db.Store, error) {
	var store db.Store
	err := p.db.WithContext(ctx).First(&store, id).Error
	if err != nil {
		p.logger.Error("Failed to get store by ID", "error", err, "storeID", id)
	}
	return &store, err
}

func (p *storePersistence) GetStoreByChatID(ctx context.Context, chatID int64) (*db.Store, error) {
	var store db.Store
	err := p.db.WithContext(ctx).Where("telegram_chat_id = ?", chatID).First(&store).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		p.logger.Error("Failed to get store by chat ID", "error", err, "chatID", chatID)
	}
	return &store, err
}

func (p *storePersistence) GetStoresBySellerID(ctx context.Context, sellerID int64) ([]db.Store, error) {
	var stores []db.Store
	err := p.db.WithContext(ctx).Where("seller_id = ?", sellerID).Find(&stores).Error
	if err != nil {
		p.logger.Error("Failed to get stores by seller ID", "error", err, "sellerID", sellerID)
	}
	return stores, err
}

func (p *storePersistence) UpdateStore(ctx context.Context, store *db.Store) error {
	err := p.db.WithContext(ctx).Save(store).Error
	if err != nil {
		p.logger.Error("Failed to update store", "error", err, "storeID", store.ID)
	}
	return err
}

func (p *storePersistence) IncrementStoreViews(ctx context.Context, storeIDs []int64) error {
	if len(storeIDs) == 0 {
		return nil
	}

	// Use raw SQL for high performance batch upsert
	query := `
		INSERT INTO store_stats (store_id, views, updated_at)
		SELECT unnest(?::bigint[]), 1, NOW()
		ON CONFLICT (store_id) 
		DO UPDATE SET views = store_stats.views + 1, updated_at = EXCLUDED.updated_at
	`
	
	err := p.db.WithContext(ctx).Exec(query, storeIDs).Error
	if err != nil {
		p.logger.Error("Failed to increment store views", "error", err, "count", len(storeIDs))
	}
	return err
}

// ProductStorage
type productPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewProductPersistence(db *gorm.DB, logger platform.Logger) storage.ProductStorage {
	return &productPersistence{db: db, logger: logger}
}

func (p *productPersistence) CreateProduct(ctx context.Context, product *db.Product) error {
	err := p.db.WithContext(ctx).Create(product).Error
	if err != nil {
		p.logger.Error("Failed to create product", "error", err)
	}
	return err
}

func (p *productPersistence) GetProductByID(ctx context.Context, id int64) (*db.Product, error) {
	var product db.Product
	err := p.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		p.logger.Error("Failed to get product by ID", "error", err, "productID", id)
	}
	return &product, err
}

func (p *productPersistence) GetProductsByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]db.Product, error) {
	var products []db.Product
	err := p.db.WithContext(ctx).Where("store_id = ?", storeID).Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		p.logger.Error("Failed to get products by store ID", "error", err, "storeID", storeID)
	}
	return products, err
}

func (p *productPersistence) GetProductsTotal(ctx context.Context, storeID int64) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&db.Product{}).Where("store_id = ?", storeID).Count(&count).Error
	if err != nil {
		p.logger.Error("Failed to get total products", "error", err, "storeID", storeID)
	}
	return count, err
}

func (p *productPersistence) ListAllProducts(ctx context.Context, filter dto.ProductFilterParams) ([]db.Product, int64, error) {
	var products []db.Product
	var count int64

	query := p.db.WithContext(ctx).Model(&db.Product{})

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if filter.Query != "" {
		searchTerm := "%" + filter.Query + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.MinStock != nil {
		query = query.Where("stock >= ?", *filter.MinStock)
	}

	if filter.MaxStock != nil {
		query = query.Where("stock <= ?", *filter.MaxStock)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(filter.GetLimit()).Offset(filter.GetOffset()).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (p *productPersistence) UpdateProduct(ctx context.Context, product *db.Product) error {
	err := p.db.WithContext(ctx).Save(product).Error
	if err != nil {
		p.logger.Error("Failed to update product", "error", err, "productID", product.ID)
	}
	return err
}

func (p *productPersistence) DeleteProduct(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Delete(&db.Product{}, id).Error
	if err != nil {
		p.logger.Error("Failed to delete product", "error", err, "productID", id)
	}
	return err
}

// OrderStorage
type orderPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewOrderPersistence(db *gorm.DB, logger platform.Logger) storage.OrderStorage {
	return &orderPersistence{db: db, logger: logger}
}

func (p *orderPersistence) CreateOrder(ctx context.Context, order *db.Order) error {
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		// Stock adjustment should be handled in module layer but DB transaction is here
		return nil
	})
	if err != nil {
		p.logger.Error("Failed to create order", "error", err)
	}
	return err
}

func (p *orderPersistence) GetOrderByID(ctx context.Context, id int64) (*db.Order, error) {
	var order db.Order
	err := p.db.WithContext(ctx).Preload("Items.Product").First(&order, id).Error
	if err != nil {
		p.logger.Error("Failed to get order by ID", "error", err, "orderID", id)
	}
	return &order, err
}

func (p *orderPersistence) GetOrdersByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]db.Order, error) {
	var orders []db.Order
	err := p.db.WithContext(ctx).Where("store_id = ?", storeID).Limit(limit).Offset(offset).Find(&orders).Error
	if err != nil {
		p.logger.Error("Failed to get orders by store ID", "error", err, "storeID", storeID)
	}
	return orders, err
}

func (p *orderPersistence) GetOrdersTotalByStoreID(ctx context.Context, storeID int64) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&db.Order{}).Where("store_id = ?", storeID).Count(&count).Error
	if err != nil {
		p.logger.Error("Failed to get total orders by store ID", "error", err, "storeID", storeID)
	}
	return count, err
}

func (p *orderPersistence) GetOrdersByCustomerID(ctx context.Context, customerID int64, limit, offset int) ([]db.Order, error) {
	var orders []db.Order
	err := p.db.WithContext(ctx).Where("user_id = ?", customerID).Limit(limit).Offset(offset).Find(&orders).Error
	if err != nil {
		p.logger.Error("Failed to get orders by customer ID", "error", err, "customerID", customerID)
	}
	return orders, err
}

func (p *orderPersistence) GetOrdersTotalByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&db.Order{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		p.logger.Error("Failed to get total orders by user ID", "error", err, "userID", userID)
	}
	return count, err
}

func (p *orderPersistence) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	err := p.db.WithContext(ctx).Model(&db.Order{}).Where("id = ?", orderID).Update("status", status).Error
	if err != nil {
		p.logger.Error("Failed to update order status", "error", err, "orderID", orderID, "status", status)
	}
	return err
}

// WalletStorage
type walletPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewWalletPersistence(db *gorm.DB, logger platform.Logger) storage.WalletStorage {
	return &walletPersistence{db: db, logger: logger}
}

func (p *walletPersistence) GetWalletByStoreID(ctx context.Context, storeID int64) (*db.Wallet, error) {
	var wallet db.Wallet
	err := p.db.WithContext(ctx).Where("store_id = ?", storeID).First(&wallet).Error
	if err != nil {
		p.logger.Error("Failed to get wallet by store ID", "error", err, "storeID", storeID)
	}
	return &wallet, err
}

func (p *walletPersistence) UpdateWalletBalance(ctx context.Context, storeID int64, amount float64) error {
	err := p.db.WithContext(ctx).Model(&db.Wallet{}).Where("store_id = ?", storeID).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		p.logger.Error("Failed to update wallet balance", "error", err, "storeID", storeID, "amount", amount)
	}
	return err
}

// CategoryStorage
type categoryPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewCategoryPersistence(db *gorm.DB, logger platform.Logger) storage.CategoryStorage {
	return &categoryPersistence{db: db, logger: logger}
}

func (p *categoryPersistence) CreateCategory(ctx context.Context, category *db.Category) error {
	err := p.db.WithContext(ctx).Create(category).Error
	if err != nil {
		p.logger.Error("Failed to create category", "error", err, "name", category.Name)
	}
	return err
}

func (p *categoryPersistence) GetAllCategories(ctx context.Context, limit, offset int) ([]db.Category, int64, error) {
	var categories []db.Category
	var count int64
	err := p.db.WithContext(ctx).Model(&db.Category{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = p.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&categories).Error
	return categories, count, err
}

func (p *categoryPersistence) GetCategoriesByStoreID(ctx context.Context, storeID int64) ([]db.Category, error) {
	var categories []db.Category
	// Get store specific + global (storeID = 0)
	err := p.db.WithContext(ctx).Where("store_id = ? OR store_id = 0", storeID).Find(&categories).Error
	return categories, err
}

func (p *categoryPersistence) GetCategoryByName(ctx context.Context, name string, storeID int64) (*db.Category, error) {
	var category db.Category
	err := p.db.WithContext(ctx).Where("name = ? AND (store_id = ? OR store_id = 0)", name, storeID).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CartStorage
type cartPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewCartPersistence(db *gorm.DB, logger platform.Logger) storage.CartStorage {
	return &cartPersistence{db: db, logger: logger}
}

func (p *cartPersistence) GetCart(ctx context.Context, userID int64) (map[string]int, error) {
	var items []db.CartItem
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).Find(&items).Error
	if err != nil {
		p.logger.Error("Failed to get cart items from DB", "error", err, "userID", userID)
		return nil, err
	}

	cart := make(map[string]int)
	for _, item := range items {
		cart[fmt.Sprintf("p:%d", item.ProductID)] = item.Quantity
	}
	return cart, nil
}

func (p *cartPersistence) AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error {
	if quantity <= 0 {
		return p.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).Delete(&db.CartItem{}).Error
	}

	item := db.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	// Use OnConflict to handle update if already exists
	return p.db.WithContext(ctx).Save(&item).Error
}

func (p *cartPersistence) UpdateCartItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	if quantity <= 0 {
		return p.RemoveFromCart(ctx, userID, productID)
	}

	// Update existing record, or create if not exists
	item := db.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	return p.db.WithContext(ctx).Save(&item).Error
}

func (p *cartPersistence) RemoveFromCart(ctx context.Context, userID int64, productID int64) error {
	return p.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).Delete(&db.CartItem{}).Error
}

func (p *cartPersistence) ClearCart(ctx context.Context, userID int64) error {
	return p.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&db.CartItem{}).Error
}
