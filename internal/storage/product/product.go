package product

import (
	"context"
	"log"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/google/uuid"
)

type productStorage struct {
	db persistencedb.PersistenceDb
}

func InitProductStorage(db persistencedb.PersistenceDb) storage.ProductStorage {
	return &productStorage{
		db: db,
	}
}

func (ps *productStorage) CreateProduct(ctx context.Context, product dto.Product) (error, uuid.UUID) {

	p := db.Product{
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
	}
	res := ps.db.WithContext(ctx).Create(&p)

	if err := res.Error; err != nil {

		log.Println("cant create product ::", err)
		return err, uuid.New()
	}

	return nil, p.ID
}
