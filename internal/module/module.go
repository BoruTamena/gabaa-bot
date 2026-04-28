package module

import (
	"context"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
)

type AuthModule interface {
	TelegramAuth(ctx context.Context, initData string) (*dto.AuthResponse, error)
}

type UserModule interface {
	GetOrCreateUser(ctx context.Context, telegramID int64, username string) (*dto.User, error)
}

type StoreModule interface {
	CreateStore(ctx context.Context, userID int64, req dto.CreateStoreRequest) (*dto.Store, error)
	GetAdminDashboard(ctx context.Context, userID int64, chatID int64) (string, *dto.Store, error)
	GetStore(ctx context.Context, id int64) (*dto.Store, error)
	UpdateStore(ctx context.Context, id int64, req dto.UpdateStoreRequest) (*dto.Store, error)
}

type ProductModule interface {
	CreateProduct(ctx context.Context, sellerID int64, storeID int64, req dto.CreateProductRequest) (*dto.Product, error)
	GetProduct(ctx context.Context, id int64) (*dto.Product, error)
	ListProducts(ctx context.Context, storeID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	ListAllProducts(ctx context.Context, filter dto.ProductFilterParams) (*dto.PaginatedResponse, error)
	UpdateProduct(ctx context.Context, id int64, req dto.UpdateProductRequest) (*dto.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	PostProduct(ctx context.Context, productID int64, storeID int64) (*dto.Product, error)
}

type OrderModule interface {
	AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error
	GetCart(ctx context.Context, userID int64) (map[int64]int, error)
	GetUserCart(ctx context.Context, userID int64) (*dto.CartResponse, error)
	Checkout(ctx context.Context, userID int64, storeID int64) (*dto.Order, error)
	ListOrders(ctx context.Context, storeID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	GetUserOrders(ctx context.Context, userID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
}

type WalletModule interface {
	GetBalance(ctx context.Context, storeID int64) (float64, error)
	CreditWallet(ctx context.Context, storeID int64, amount float64) error
}

type CategoryModule interface {
	CreateCategory(ctx context.Context, storeID int64, req dto.CreateCategoryRequest) (*dto.Category, error)
	ListAllCategories(ctx context.Context, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	ListStoreCategories(ctx context.Context, storeID int64) ([]dto.Category, error)
}
