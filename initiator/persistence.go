package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/storage/persistence"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
)


type Persistence struct {
	UserStorage           storage.UserStorage
	StoreStorage          storage.StoreStorage
	ProductStorage        storage.ProductStorage
	OrderStorage          storage.OrderStorage
	WalletStorage         storage.WalletStorage
	PaymentStorage        storage.PaymentStorage
	PaymentWebhookStorage storage.PaymentWebhookStorage
	EscrowStorage         storage.EscrowStorage
	WithdrawalStorage     storage.WithdrawalStorage
	CartStorage           storage.CartStorage
	CategoryStorage       storage.CategoryStorage
	AddressStorage        storage.AddressStorage
	StoryStorage          storage.StoryStorage
	FavoriteStorage       storage.FavoriteStorage
	PreferenceStorage     storage.PreferenceStorage
	RecommendationStorage storage.RecommendationStorage
	AuthSessionStorage    storage.AuthSessionStorage
	StoreKYCStorage       storage.StoreKYCStorage
	AnalyticsStorage      storage.AnalyticsStorage
	DeliveryStorage       storage.DeliveryStorage
}

func InitPersistence(db persistencedb.PersistenceDb, redis platform.Redis, logger platform.Logger) Persistence {
	return Persistence{
		UserStorage:           persistence.NewPersistence(db.DB, logger),
		StoreStorage:          persistence.NewStorePersistence(db.DB, logger),
		ProductStorage:        persistence.NewProductPersistence(db.DB, logger),
		OrderStorage:          persistence.NewOrderPersistence(db.DB, logger),
		WalletStorage:         persistence.NewWalletPersistence(db.DB, logger),
		PaymentStorage:        persistence.NewPaymentPersistence(db.DB, logger),
		PaymentWebhookStorage: persistence.NewPaymentWebhookPersistence(db.DB, logger),
		EscrowStorage:         persistence.NewEscrowPersistence(db.DB, logger),
		WithdrawalStorage:     persistence.NewWithdrawalPersistence(db.DB, logger),
		CartStorage:           persistence.NewCartPersistence(db.DB, logger),
		CategoryStorage:       persistence.NewCategoryPersistence(db.DB, logger),
		AddressStorage:        persistence.NewAddressPersistence(db.DB, logger),
		StoryStorage:          persistence.NewStoryPersistence(db.DB, logger),
		FavoriteStorage:       persistence.NewFavoritePersistence(db.DB, logger),
		PreferenceStorage:     persistence.NewPreferencePersistence(db.DB, logger),
		RecommendationStorage: persistence.NewRecommendationPersistence(db.DB, logger),
		AuthSessionStorage:    persistence.NewAuthSessionPersistence(db.DB, logger),
		StoreKYCStorage:       persistence.NewStoreKYCPersistence(db.DB, logger),
		AnalyticsStorage:      persistence.NewAnalyticsPersistence(db.DB, logger),
		DeliveryStorage:       persistence.NewDeliveryStorage(db.DB, logger),
	}
}



