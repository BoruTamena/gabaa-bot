package module

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
)

type AuthModule interface {
	TelegramAuth(ctx context.Context, initData string) (*dto.AuthResponse, error)
	StartBotLoginSession(ctx context.Context) (*dto.TelegramLoginSessionResponse, error)
	CompleteBotLoginSession(ctx context.Context, sessionID string, tgUser *dto.TelegramUser) error
	PollBotLoginSession(ctx context.Context, sessionID string) (*dto.TelegramLoginPollResponse, error)
}

type UserModule interface {
	GetOrCreateUser(ctx context.Context, telegramID int64, username string) (*dto.User, error)
}

type AddressModule interface {
	CreateAddress(ctx context.Context, userID int64, req dto.CreateAddressRequest) (*dto.Address, error)
	GetAddress(ctx context.Context, id int64) (*dto.Address, error)
	GetAddressesByUser(ctx context.Context, userID int64) ([]dto.Address, error)
	UpdateAddress(ctx context.Context, userID int64, id int64, req dto.UpdateAddressRequest) (*dto.Address, error)
	DeleteAddress(ctx context.Context, userID int64, id int64) error
	SetDefaultAddress(ctx context.Context, userID int64, id int64) error
}

type StoreModule interface {
	CreateStore(ctx context.Context, userID int64, req dto.CreateStoreRequest) (*dto.Store, error)
	GetAdminDashboard(ctx context.Context, userID int64, chatID int64) (string, *dto.Store, error)
	GetStore(ctx context.Context, id int64) (*dto.Store, error)
	GetStoreStatus(ctx context.Context, id int64) (string, error)
	UpdateStore(ctx context.Context, id int64, req dto.UpdateStoreRequest) (*dto.Store, error)
	SubmitStoreKYC(ctx context.Context, storeID, sellerID int64, req dto.SubmitStoreKYCRequest) (*dto.StoreKYCResponse, error)
	GetStoreKYC(ctx context.Context, storeID, sellerID int64) (*dto.StoreKYCResponse, error)
	ListStoreVerifications(ctx context.Context, status string) ([]dto.StoreKYCResponse, error)
	ApproveStoreKYC(ctx context.Context, storeID int64) (*dto.StoreKYCResponse, error)
	RejectStoreKYC(ctx context.Context, storeID int64, req dto.RejectStoreKYCRequest) (*dto.StoreKYCResponse, error)
	IsStoreVerified(ctx context.Context, storeID int64) (bool, error)
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

type CartModule interface {
	AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error
	GetCart(ctx context.Context, userID int64) (map[int64]int, error)
	GetUserCart(ctx context.Context, userID int64) (*dto.CartResponse, error)
	UpdateCartItem(ctx context.Context, userID int64, productID int64, action string) error
	RemoveFromCart(ctx context.Context, userID int64, productID int64) error
	ClearCart(ctx context.Context, userID int64) error
}

type OrderModule interface {
	Checkout(ctx context.Context, userID int64, storeID int64, addressID int64, medium, phone string) (*dto.CheckoutResponse, error)
	GetOrder(ctx context.Context, orderID int64) (*dto.Order, error)
	ListOrders(ctx context.Context, storeID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	GetUserOrders(ctx context.Context, userID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
	CancelOrder(ctx context.Context, userID int64, orderID int64) error
	// Merchant-scoped:
	GetMyStoreOrders(ctx context.Context, filter dto.OrderFilterParams) (*dto.PaginatedResponse, error)
	GetMyStoreOrder(ctx context.Context, storeID int64, orderID int64) (*dto.Order, error)
	UpdateMyStoreOrderStatus(ctx context.Context, storeID int64, orderID int64, status string) error
	OnPaymentSuccess(ctx context.Context, orderID int64) error
	OnPaymentFailed(ctx context.Context, orderID int64) error
	SetPaymentModule(pm PaymentModule)
}

type WalletModule interface {
	GetWalletSummary(ctx context.Context, storeID int64) (*dto.Wallet, error)
	RequestWithdrawal(ctx context.Context, storeID int64, req dto.WithdrawalRequest) (*dto.Withdrawal, error)
	ListWithdrawals(ctx context.Context, storeID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	GetMyStoreWithdrawal(ctx context.Context, storeID, withdrawalID int64) (*dto.Withdrawal, error)
}

type PaymentModule interface {
	InitiateForOrder(ctx context.Context, order *db.Order, medium, phone string) (*dto.Payment, error)
	HandleWebhook(ctx context.Context, rawBody []byte) dto.WebhookResult
	ListStoreTransactions(ctx context.Context, filter dto.PaymentFilterParams) (*dto.PaginatedResponse, error)
}

type CategoryModule interface {
	CreateCategory(ctx context.Context, storeID int64, req dto.CreateCategoryRequest) (*dto.Category, error)
	ListAllCategories(ctx context.Context, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	ListStoreCategories(ctx context.Context, storeID int64) ([]dto.Category, error)
}

type BotModule interface {
	// Any specific bot-related methods can go here
}

type UploadModule interface {
	UploadImages(ctx context.Context, files []interface{}, fileNames []string) ([]string, error)
	UploadDocuments(ctx context.Context, files []interface{}, fileNames []string) ([]string, error)
}

type StoryModule interface {
	CreateStory(ctx context.Context, storeID int64, req dto.CreateProductStoryRequest) (*dto.ProductStory, error)
	GetStory(ctx context.Context, id int64) (*dto.ProductStory, error)
	ListMyStories(ctx context.Context, filter dto.ProductStoryFilterParams) (*dto.PaginatedResponse, error)
	UpdateStory(ctx context.Context, storeID int64, storyID int64, req dto.UpdateProductStoryRequest) (*dto.ProductStory, error)
	DeleteStory(ctx context.Context, storeID int64, storyID int64) error
	ListActiveStories(ctx context.Context, params dto.PaginationParams) (*dto.PaginatedResponse, error)
	StartExpiryJob(ctx context.Context)
}

type FavoriteModule interface {
	AddFavorite(ctx context.Context, userID, productID int64) error
	RemoveFavorite(ctx context.Context, userID, productID int64) error
	ListUserFavorites(ctx context.Context, userID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error)
}

type RecommendationModule interface {
	GetPreferences(ctx context.Context, userID int64) (*dto.UserPreferences, error)
	SetPreferences(ctx context.Context, userID int64, req dto.UpdateUserPreferencesRequest) (*dto.UserPreferences, error)
	SetBotStarted(ctx context.Context, telegramUserID int64, username string) error
	SetRecommendationsEnabled(ctx context.Context, telegramUserID int64, enabled bool) error
	ToggleCategory(ctx context.Context, telegramUserID int64, category string) (added bool, err error)
	NotifyMatchingUsers(ctx context.Context, product *db.Product, sellerUserID int64)
}

type AnalyticsModule interface {
	GetSalesAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.SalesAnalytics, error)
	GetOrderAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.OrderAnalytics, error)
	GetProductAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.ProductAnalytics, error)
	GetStoryAnalytics(ctx context.Context, storeID int64, filter dto.AnalyticsFilterParams) (*dto.StoryAnalytics, error)
}
