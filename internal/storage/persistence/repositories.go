package persistence

import (
	"context"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"gorm.io/gorm"
)

type persistence struct {
	db *gorm.DB
}

func NewPersistence(db *gorm.DB) storage.UserStorage {
	return &persistence{db: db}
}

// UserStorage
func (p *persistence) CreateUser(ctx context.Context, user *db.User) error {
	return p.db.WithContext(ctx).Create(user).Error
}

func (p *persistence) GetUserByTelegramID(ctx context.Context, telegramID int64) (*db.User, error) {
	var user db.User
	err := p.db.WithContext(ctx).Where("telegram_user_id = ?", telegramID).First(&user).Error
	return &user, err
}

func (p *persistence) UpdateUser(ctx context.Context, user *db.User) error {
	return p.db.WithContext(ctx).Save(user).Error
}

// StoreStorage
type storePersistence struct {
	db *gorm.DB
}

func NewStorePersistence(db *gorm.DB) storage.StoreStorage {
	return &storePersistence{db: db}
}

func (p *storePersistence) CreateStore(ctx context.Context, store *db.Store) error {
	return p.db.WithContext(ctx).Create(store).Error
}

func (p *storePersistence) GetStoreByID(ctx context.Context, id int64) (*db.Store, error) {
	var store db.Store
	err := p.db.WithContext(ctx).First(&store, id).Error
	return &store, err
}

func (p *storePersistence) GetStoreByChatID(ctx context.Context, chatID int64) (*db.Store, error) {
	var store db.Store
	err := p.db.WithContext(ctx).Where("telegram_chat_id = ?", chatID).First(&store).Error
	return &store, err
}

func (p *storePersistence) GetStoresBySellerID(ctx context.Context, sellerID int64) ([]db.Store, error) {
	var stores []db.Store
	err := p.db.WithContext(ctx).Where("seller_id = ?", sellerID).Find(&stores).Error
	return stores, err
}

// ProductStorage
type productPersistence struct {
	db *gorm.DB
}

func NewProductPersistence(db *gorm.DB) storage.ProductStorage {
	return &productPersistence{db: db}
}

func (p *productPersistence) CreateProduct(ctx context.Context, product *db.Product) error {
	return p.db.WithContext(ctx).Create(product).Error
}

func (p *productPersistence) GetProductByID(ctx context.Context, id int64) (*db.Product, error) {
	var product db.Product
	err := p.db.WithContext(ctx).First(&product, id).Error
	return &product, err
}

func (p *productPersistence) GetProductsByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]db.Product, error) {
	var products []db.Product
	err := p.db.WithContext(ctx).Where("store_id = ?", storeID).Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

func (p *productPersistence) GetProductsTotal(ctx context.Context, storeID int64) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&db.Product{}).Where("store_id = ?", storeID).Count(&count).Error
	return count, err
}


func (p *productPersistence) UpdateProduct(ctx context.Context, product *db.Product) error {
	return p.db.WithContext(ctx).Save(product).Error
}

func (p *productPersistence) DeleteProduct(ctx context.Context, id int64) error {
	return p.db.WithContext(ctx).Delete(&db.Product{}, id).Error
}

// OrderStorage
type orderPersistence struct {
	db *gorm.DB
}

func NewOrderPersistence(db *gorm.DB) storage.OrderStorage {
	return &orderPersistence{db: db}
}

func (p *orderPersistence) CreateOrder(ctx context.Context, order *db.Order) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		// Stock adjustment should be handled in module layer but DB transaction is here
		return nil
	})
}

func (p *orderPersistence) GetOrderByID(ctx context.Context, id int64) (*db.Order, error) {
	var order db.Order
	err := p.db.WithContext(ctx).Preload("Items.Product").First(&order, id).Error
	return &order, err
}

func (p *orderPersistence) GetOrdersByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]db.Order, error) {
	var orders []db.Order
	err := p.db.WithContext(ctx).Where("store_id = ?", storeID).Limit(limit).Offset(offset).Find(&orders).Error
	return orders, err
}

func (p *orderPersistence) GetOrdersTotalByStoreID(ctx context.Context, storeID int64) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&db.Order{}).Where("store_id = ?", storeID).Count(&count).Error
	return count, err
}

func (p *orderPersistence) GetOrdersByCustomerID(ctx context.Context, customerID int64, limit, offset int) ([]db.Order, error) {
	var orders []db.Order
	err := p.db.WithContext(ctx).Where("user_id = ?", customerID).Limit(limit).Offset(offset).Find(&orders).Error
	return orders, err
}


func (p *orderPersistence) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	return p.db.WithContext(ctx).Model(&db.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// WalletStorage
type walletPersistence struct {
	db *gorm.DB
}

func NewWalletPersistence(db *gorm.DB) storage.WalletStorage {
	return &walletPersistence{db: db}
}

func (p *walletPersistence) GetWalletByStoreID(ctx context.Context, storeID int64) (*db.Wallet, error) {
	var wallet db.Wallet
	err := p.db.WithContext(ctx).Where("store_id = ?", storeID).First(&wallet).Error
	return &wallet, err
}

func (p *walletPersistence) UpdateWalletBalance(ctx context.Context, storeID int64, amount float64) error {
	return p.db.WithContext(ctx).Model(&db.Wallet{}).Where("store_id = ?", storeID).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}
