package storage

import (
	"context"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *db.User) error
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*db.User, error)
	UpdateUser(ctx context.Context, user *db.User) error
}

type StoreStorage interface {
	CreateStore(ctx context.Context, store *db.Store) error
	GetStoreByID(ctx context.Context, id int64) (*db.Store, error)
	GetStoreByChatID(ctx context.Context, chatID int64) (*db.Store, error)
	GetStoresBySellerID(ctx context.Context, sellerID int64) ([]db.Store, error)
}

type ProductStorage interface {
	CreateProduct(ctx context.Context, product *db.Product) error
	GetProductByID(ctx context.Context, id int64) (*db.Product, error)
	GetProductsByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]db.Product, error)
	GetProductsTotal(ctx context.Context, storeID int64) (int64, error)
	UpdateProduct(ctx context.Context, product *db.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

type OrderStorage interface {
	CreateOrder(ctx context.Context, order *db.Order) error
	GetOrderByID(ctx context.Context, id int64) (*db.Order, error)
	GetOrdersByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]db.Order, error)
	GetOrdersTotalByStoreID(ctx context.Context, storeID int64) (int64, error)
	GetOrdersByCustomerID(ctx context.Context, customerID int64, limit, offset int) ([]db.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
}


type WalletStorage interface {
	GetWalletByStoreID(ctx context.Context, storeID int64) (*db.Wallet, error)
	UpdateWalletBalance(ctx context.Context, storeID int64, amount float64) error
}

type CartStorage interface {
	GetCart(ctx context.Context, userID int64) (map[string]int, error)
	AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error
	ClearCart(ctx context.Context, userID int64) error
}

