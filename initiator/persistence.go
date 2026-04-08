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
}

func InitPersistence(db persistencedb.PersistenceDb, redis platform.Redis) Persistence {
	return Persistence{
		UserStorage:    persistence.NewPersistence(db.DB),
		StoreStorage:   persistence.NewStorePersistence(db.DB),
		ProductStorage: persistence.NewProductPersistence(db.DB),
		OrderStorage:   persistence.NewOrderPersistence(db.DB),
		WalletStorage:  persistence.NewWalletPersistence(db.DB),
		CartStorage:    cache.NewCartCache(redis),
	}
}


