package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/lib/pq"
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

	err := p.db.WithContext(ctx).Exec(query, pq.Array(storeIDs)).Error
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

	if filter.StoreID != 0 {
		query = query.Where("store_id = ?", filter.StoreID)
	}

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
	err := p.db.WithContext(ctx).Preload("Items.Product").Preload("ShippingAddress").First(&order, id).Error
	if err != nil {
		p.logger.Error("Failed to get order by ID", "error", err, "orderID", id)
	}
	return &order, err
}

func (p *orderPersistence) GetOrdersByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]db.Order, error) {
	var orders []db.Order
	err := p.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("ShippingAddress").
		Where("store_id = ?", storeID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error
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
	err := p.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("ShippingAddress").
		Where("user_id = ?", customerID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error
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

func (p *orderPersistence) UpdateOrderDispatch(ctx context.Context, orderID int64, status string, agentID, routeID int64) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":            status,
		"delivery_agent_id": agentID,
		"dispatched_at":     now,
	}
	if routeID > 0 {
		updates["delivery_route_id"] = routeID
	}
	err := p.db.WithContext(ctx).Model(&db.Order{}).Where("id = ?", orderID).Updates(updates).Error
	if err != nil {
		p.logger.Error("Failed to update order dispatch", "error", err, "orderID", orderID)
	}
	return err
}

func (p *orderPersistence) GetOrdersByDeliveryAgentID(ctx context.Context, agentID int64, status string, limit, offset int) ([]db.Order, error) {
	var orders []db.Order
	query := p.db.WithContext(ctx).
		Preload("User").
		Preload("Store").
		Preload("Items").
		Preload("Items.Product").
		Preload("ShippingAddress").
		Where("delivery_agent_id = ?", agentID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&orders).Error
	return orders, err
}

func (p *orderPersistence) GetOrdersTotalByDeliveryAgentID(ctx context.Context, agentID int64, status string) (int64, error) {
	var count int64
	query := p.db.WithContext(ctx).Model(&db.Order{}).Where("delivery_agent_id = ?", agentID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&count).Error
	return count, err
}

func (p *orderPersistence) GetOrdersByFilter(ctx context.Context, filter dto.OrderFilterParams) ([]db.Order, int64, error) {
	var orders []db.Order
	var count int64

	query := p.db.WithContext(ctx).Model(&db.Order{}).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("ShippingAddress").
		Where("store_id = ?", filter.StoreID)

	if filter.OrderID != nil {
		query = query.Where("id = ?", *filter.OrderID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(filter.GetLimit()).Offset(filter.GetOffset()).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, count, nil
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

func (p *walletPersistence) GetOrCreateWallet(ctx context.Context, storeID int64) (*db.Wallet, error) {
	wallet, err := p.GetWalletByStoreID(ctx, storeID)
	if err == nil {
		return wallet, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	wallet = &db.Wallet{StoreID: storeID, Currency: "ETB"}
	if err := p.db.WithContext(ctx).Create(wallet).Error; err != nil {
		p.logger.Error("Failed to create wallet", "error", err, "storeID", storeID)
		return nil, err
	}
	return wallet, nil
}

func (p *walletPersistence) AddPendingBalance(ctx context.Context, storeID int64, amount float64) error {
	_, err := p.GetOrCreateWallet(ctx, storeID)
	if err != nil {
		return err
	}
	err = p.db.WithContext(ctx).Model(&db.Wallet{}).Where("store_id = ?", storeID).
		UpdateColumn("pending_balance", gorm.Expr("pending_balance + ?", amount)).Error
	if err != nil {
		p.logger.Error("Failed to add pending balance", "error", err, "storeID", storeID, "amount", amount)
	}
	return err
}

func (p *walletPersistence) ReleaseEscrowFunds(ctx context.Context, storeID int64, amount float64) error {
	err := p.db.WithContext(ctx).Model(&db.Wallet{}).Where("store_id = ?", storeID).Updates(map[string]interface{}{
		"pending_balance":   gorm.Expr("pending_balance - ?", amount),
		"available_balance": gorm.Expr("available_balance + ?", amount),
		"total_earned":      gorm.Expr("total_earned + ?", amount),
	}).Error
	if err != nil {
		p.logger.Error("Failed to release escrow funds", "error", err, "storeID", storeID, "amount", amount)
	}
	return err
}

func (p *walletPersistence) LockForWithdrawal(ctx context.Context, storeID int64, amount float64) error {
	result := p.db.WithContext(ctx).Model(&db.Wallet{}).Where("store_id = ? AND available_balance >= ?", storeID, amount).Updates(map[string]interface{}{
		"available_balance": gorm.Expr("available_balance - ?", amount),
		"locked_balance":    gorm.Expr("locked_balance + ?", amount),
	})
	if result.Error != nil {
		p.logger.Error("Failed to lock funds for withdrawal", "error", result.Error, "storeID", storeID, "amount", amount)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient available balance")
	}
	return nil
}

func (p *walletPersistence) UnlockWithdrawal(ctx context.Context, storeID int64, amount float64) error {
	err := p.db.WithContext(ctx).Model(&db.Wallet{}).Where("store_id = ?", storeID).Updates(map[string]interface{}{
		"locked_balance":    gorm.Expr("locked_balance - ?", amount),
		"available_balance": gorm.Expr("available_balance + ?", amount),
	}).Error
	if err != nil {
		p.logger.Error("Failed to unlock withdrawal funds", "error", err, "storeID", storeID, "amount", amount)
	}
	return err
}

func (p *walletPersistence) CompleteWithdrawal(ctx context.Context, storeID int64, amount float64) error {
	err := p.db.WithContext(ctx).Model(&db.Wallet{}).Where("store_id = ?", storeID).Updates(map[string]interface{}{
		"locked_balance":   gorm.Expr("locked_balance - ?", amount),
		"total_withdrawn":  gorm.Expr("total_withdrawn + ?", amount),
	}).Error
	if err != nil {
		p.logger.Error("Failed to complete withdrawal", "error", err, "storeID", storeID, "amount", amount)
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

func (p *categoryPersistence) GetCategoryByID(ctx context.Context, id int64) (*db.Category, error) {
	var category db.Category
	err := p.db.WithContext(ctx).First(&category, id).Error
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

// AddressStorage
type addressPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewAddressPersistence(db *gorm.DB, logger platform.Logger) storage.AddressStorage {
	return &addressPersistence{db: db, logger: logger}
}

func (p *addressPersistence) CreateAddress(ctx context.Context, address *db.Address) error {
	return p.db.WithContext(ctx).Create(address).Error
}

func (p *addressPersistence) GetAddressByID(ctx context.Context, id int64) (*db.Address, error) {
	var address db.Address
	err := p.db.WithContext(ctx).First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (p *addressPersistence) GetAddressesByUserID(ctx context.Context, userID int64) ([]db.Address, error) {
	var addresses []db.Address
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error
	return addresses, err
}

func (p *addressPersistence) UpdateAddress(ctx context.Context, address *db.Address) error {
	return p.db.WithContext(ctx).Save(address).Error
}

func (p *addressPersistence) DeleteAddress(ctx context.Context, id int64) error {
	return p.db.WithContext(ctx).Delete(&db.Address{}, id).Error
}

func (p *addressPersistence) ClearDefaultAddress(ctx context.Context, userID int64) error {
	return p.db.WithContext(ctx).Model(&db.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error
}

// StoryStorage
type storyPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewStoryPersistence(db *gorm.DB, logger platform.Logger) storage.StoryStorage {
	return &storyPersistence{db: db, logger: logger}
}

func (p *storyPersistence) CreateStory(ctx context.Context, story *db.ProductStory) error {
	err := p.db.WithContext(ctx).Create(story).Error
	if err != nil {
		p.logger.Error("Failed to create story", "error", err, "productID", story.ProductID)
	}
	return err
}

func (p *storyPersistence) GetStoryByID(ctx context.Context, id int64) (*db.ProductStory, error) {
	var story db.ProductStory
	err := p.db.WithContext(ctx).
		Preload("Product").
		First(&story, id).Error
	if err != nil {
		p.logger.Error("Failed to get story by ID", "error", err, "storyID", id)
	}
	return &story, err
}

func (p *storyPersistence) ListStoriesByStore(ctx context.Context, filter dto.ProductStoryFilterParams) ([]db.ProductStory, int64, error) {
	var stories []db.ProductStory
	var count int64

	query := p.db.WithContext(ctx).Model(&db.ProductStory{}).
		Where("store_id = ?", filter.StoreID)

	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Product").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Order("created_at DESC").
		Find(&stories).Error
	if err != nil {
		p.logger.Error("Failed to list stories by store", "error", err, "storeID", filter.StoreID)
	}

	return stories, count, err
}

func (p *storyPersistence) ListActiveStories(ctx context.Context, params dto.PaginationParams) ([]db.ProductStory, int64, error) {
	var stories []db.ProductStory
	var count int64

	query := p.db.WithContext(ctx).Model(&db.ProductStory{}).
		Where("is_active = true AND starts_at <= NOW() AND ends_at >= NOW()")

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Product").
		Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order("created_at DESC").
		Find(&stories).Error
	if err != nil {
		p.logger.Error("Failed to list active stories", "error", err)
	}

	return stories, count, err
}

func (p *storyPersistence) UpdateStory(ctx context.Context, story *db.ProductStory) error {
	err := p.db.WithContext(ctx).Save(story).Error
	if err != nil {
		p.logger.Error("Failed to update story", "error", err, "storyID", story.ID)
	}
	return err
}

func (p *storyPersistence) DeleteStory(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Delete(&db.ProductStory{}, id).Error
	if err != nil {
		p.logger.Error("Failed to delete story", "error", err, "storyID", id)
	}
	return err
}

func (p *storyPersistence) IncrementStoryViews(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).
		Model(&db.ProductStory{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
	if err != nil {
		p.logger.Error("Failed to increment story views", "error", err, "storyID", id)
	}
	return err
}

func (p *storyPersistence) ExpireEndedStories(ctx context.Context) (int64, error) {
	result := p.db.WithContext(ctx).Model(&db.ProductStory{}).
		Where("is_active = ? AND ends_at < ?", true, time.Now()).
		Update("is_active", false)
	if result.Error != nil {
		p.logger.Error("Failed to expire ended stories", "error", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// FavoriteStorage
type favoritePersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewFavoritePersistence(db *gorm.DB, logger platform.Logger) storage.FavoriteStorage {
	return &favoritePersistence{db: db, logger: logger}
}

func (p *favoritePersistence) AddFavorite(ctx context.Context, favorite *db.Favorite) error {
	err := p.db.WithContext(ctx).Create(favorite).Error
	if err != nil {
		p.logger.Error("Failed to add favorite", "error", err, "userID", favorite.UserID, "productID", favorite.ProductID)
	}
	return err
}

func (p *favoritePersistence) RemoveFavorite(ctx context.Context, userID, productID int64) error {
	err := p.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).Delete(&db.Favorite{}).Error
	if err != nil {
		p.logger.Error("Failed to remove favorite", "error", err, "userID", userID, "productID", productID)
	}
	return err
}

func (p *favoritePersistence) ListUserFavorites(ctx context.Context, userID int64, params dto.PaginationParams) ([]db.Favorite, int64, error) {
	var favorites []db.Favorite
	var count int64

	query := p.db.WithContext(ctx).Model(&db.Favorite{}).Where("user_id = ?", userID)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Product").
		Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order("created_at DESC").
		Find(&favorites).Error
	if err != nil {
		p.logger.Error("Failed to list user favorites", "error", err, "userID", userID)
	}

	return favorites, count, err
}

func (p *favoritePersistence) IsFavorite(ctx context.Context, userID, productID int64) (bool, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&db.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error
	if err != nil {
		p.logger.Error("Failed to check if favorite", "error", err, "userID", userID, "productID", productID)
		return false, err
	}
	return count > 0, nil
}

// PreferenceStorage
type preferencePersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewPreferencePersistence(db *gorm.DB, logger platform.Logger) storage.PreferenceStorage {
	return &preferencePersistence{db: db, logger: logger}
}

func (p *preferencePersistence) GetUserPreferences(ctx context.Context, userID int64) ([]string, error) {
	pref, err := p.getPreferenceRow(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []string{}, nil
		}
		p.logger.Error("Failed to get user preferences", "error", err, "userID", userID)
		return nil, err
	}

	return categoriesFromJSON(pref.Categories), nil
}

func (p *preferencePersistence) SetUserPreferences(ctx context.Context, userID int64, categories []string) error {
	normalized := normalizeCategories(categories)
	categoriesJSON, err := categoriesToJSON(normalized)
	if err != nil {
		return err
	}

	pref, err := p.getPreferenceRow(ctx, userID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			p.logger.Error("Failed to load user preferences for update", "error", err, "userID", userID)
			return err
		}

		pref = &db.UserCategoryPreference{
			UserID:     userID,
			Categories: categoriesJSON,
		}
		if err := p.db.WithContext(ctx).Create(pref).Error; err != nil {
			p.logger.Error("Failed to create user preferences", "error", err, "userID", userID)
			return err
		}
		return nil
	}

	pref.Categories = categoriesJSON
	if err := p.db.WithContext(ctx).Save(pref).Error; err != nil {
		p.logger.Error("Failed to update user preferences", "error", err, "userID", userID)
		return err
	}
	return nil
}

func (p *preferencePersistence) ToggleUserCategory(ctx context.Context, userID int64, category string) (bool, error) {
	normalized := stringsTrimSpace(category)
	if normalized == "" {
		return false, fmt.Errorf("category cannot be empty")
	}

	current, err := p.GetUserPreferences(ctx, userID)
	if err != nil {
		return false, err
	}

	updated, added := toggleCategoryInList(current, normalized)
	if err := p.SetUserPreferences(ctx, userID, updated); err != nil {
		return false, err
	}
	return added, nil
}

func (p *preferencePersistence) GetUsersByCategories(ctx context.Context, categories []string) ([]db.User, error) {
	normalized := normalizeCategories(categories)
	if len(normalized) == 0 {
		return nil, nil
	}

	lowered := make([]string, len(normalized))
	for i, category := range normalized {
		lowered[i] = stringsToLower(category)
	}

	var users []db.User
	err := p.db.WithContext(ctx).
		Table("users").
		Select("DISTINCT users.*").
		Joins("JOIN user_category_preferences ON user_category_preferences.user_id = users.id AND user_category_preferences.deleted_at IS NULL").
		Where("users.recommendations_enabled = ?", true).
		Where("users.bot_started = ?", true).
		Where("users.telegram_user_id IS NOT NULL").
		Where(`EXISTS (
			SELECT 1
			FROM jsonb_array_elements_text(user_category_preferences.categories) AS cat
			WHERE lower(cat) IN ?
		)`, lowered).
		Find(&users).Error
	if err != nil {
		p.logger.Error("Failed to get users by categories", "error", err)
	}
	return users, err
}

func (p *preferencePersistence) getPreferenceRow(ctx context.Context, userID int64) (*db.UserCategoryPreference, error) {
	var pref db.UserCategoryPreference
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error
	if err != nil {
		return nil, err
	}
	return &pref, nil
}

func normalizeCategories(categories []string) []string {
	seen := make(map[string]bool)
	normalized := make([]string, 0, len(categories))
	for _, category := range categories {
		value := stringsTrimSpace(category)
		if value == "" {
			continue
		}
		key := stringsToLower(value)
		if seen[key] {
			continue
		}
		seen[key] = true
		normalized = append(normalized, value)
	}
	return normalized
}

func categoriesToJSON(categories []string) (string, error) {
	if categories == nil {
		categories = []string{}
	}
	bytes, err := json.Marshal(categories)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func categoriesFromJSON(raw string) []string {
	if raw == "" {
		return []string{}
	}

	var categories []string
	if err := json.Unmarshal([]byte(raw), &categories); err != nil {
		return []string{}
	}
	if categories == nil {
		return []string{}
	}
	return categories
}

func toggleCategoryInList(categories []string, category string) ([]string, bool) {
	target := stringsToLower(category)
	for i, existing := range categories {
		if stringsToLower(existing) == target {
			updated := append([]string{}, categories[:i]...)
			updated = append(updated, categories[i+1:]...)
			return updated, false
		}
	}
	return append(categories, category), true
}

// RecommendationStorage
type recommendationPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewRecommendationPersistence(db *gorm.DB, logger platform.Logger) storage.RecommendationStorage {
	return &recommendationPersistence{db: db, logger: logger}
}

func (p *recommendationPersistence) WasNotified(ctx context.Context, userID, productID int64) (bool, error) {
	var count int64
	err := p.db.WithContext(ctx).
		Model(&db.ProductRecommendation{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	if err != nil {
		p.logger.Error("Failed to check recommendation notification", "error", err, "userID", userID, "productID", productID)
		return false, err
	}
	return count > 0, nil
}

func (p *recommendationPersistence) RecordNotification(ctx context.Context, userID, productID int64) error {
	rec := &db.ProductRecommendation{
		UserID:    userID,
		ProductID: productID,
		SentAt:    timeNow(),
	}
	err := p.db.WithContext(ctx).Create(rec).Error
	if err != nil {
		p.logger.Error("Failed to record recommendation notification", "error", err, "userID", userID, "productID", productID)
	}
	return err
}

func stringsTrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func stringsToLower(s string) string {
	return strings.ToLower(s)
}

func timeNow() time.Time {
	return time.Now()
}
