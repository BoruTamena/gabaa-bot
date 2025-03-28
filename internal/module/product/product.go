package product

import (
	"context"
	"log"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type productModule struct {
	productStorage storage.ProductStorage
}

func InitProductModule(pStorage storage.ProductStorage) module.ProductModule {

	return &productModule{
		productStorage: pStorage,
	}

}

func (p *productModule) CreateProduct(ctx context.Context, product dto.Product) error {

	if err := product.Validate(); err != nil {
		log.Print("validation error::", err)
		return err
	}

	if err := p.productStorage.CreateProduct(ctx, product); err != nil {
		return err
	}

	//TODO use telebot api to add create order inline button to post

	return nil

}
