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
		SellerId:    product.SellerId,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
	}
	res := ps.db.WithContext(ctx).Create(&p)

	if err := res.Error; err != nil {

		log.Println("can't create product ::", err)
		return err, uuid.New()
	}

	return nil, p.ID
}

func (ps *productStorage) GetProductByID(ctx context.Context, id string) (dto.Product, error) {
	var product db.Product
	res := ps.db.WithContext(ctx).Where("id = ?", id).First(&product)

	if err := res.Error; err != nil {
		log.Println("can't get product ::", err)
		return dto.Product{}, err
	}

	return dto.Product{
		ID:          product.ID.String(),
		SellerId:    product.SellerId,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}
