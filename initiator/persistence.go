package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/storage/order"
	"github.com/BoruTamena/gabaa-bot/internal/storage/product"
	"github.com/BoruTamena/gabaa-bot/internal/storage/user"

	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type Persistance struct {
	productStorage storage.ProductStorage
	orderStorage   storage.OrderStorage
	userStorage    storage.UserStorage
}

func InitPersistence(db persistencedb.PersistenceDb) Persistance {
	return Persistance{
		productStorage: product.InitProductStorage(db),
		orderStorage:   order.InitOrderStorage(db),
		userStorage:    user.InitUserStorage(db),
	}
}
