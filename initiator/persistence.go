package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/storage/persistence"
	"github.com/BoruTamena/gabaa-bot/internal/storage/cache"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
)


type Persistence struct {
	UserStorage    storage.UserStorage
	StoreStorage   storage.StoreStorage
	ProductStorage storage.ProductStorage
	OrderStorage   storage.OrderStorage
	WalletStorage  storage.WalletStorage
	CartStorage    storage.CartStorage
	CategoryStorage storage.CategoryStorage
}

func InitPersistence(db persistencedb.PersistenceDb, redis platform.Redis, logger platform.Logger) Persistence {
	return Persistence{
		UserStorage:    persistence.NewPersistence(db.DB, logger),
		StoreStorage:   persistence.NewStorePersistence(db.DB, logger),
		ProductStorage: persistence.NewProductPersistence(db.DB, logger),
		OrderStorage:   persistence.NewOrderPersistence(db.DB, logger),
		WalletStorage:  persistence.NewWalletPersistence(db.DB, logger),
		CartStorage:    cache.NewCartCache(redis),
		CategoryStorage: persistence.NewCategoryPersistence(db.DB, logger),
	}
}


