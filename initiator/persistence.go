package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/storage/product"

	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type Persistance struct {
	productStorage storage.ProductStorage
}

func InitPersistence(db persistencedb.PersistenceDb) Persistance {
	return Persistance{
		productStorage: product.InitProductStorage(db),
	}
}
