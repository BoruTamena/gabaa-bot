package product

import (
	"context"
	"log"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"gopkg.in/telebot.v4"
)

type productModule struct {
	productStorage storage.ProductStorage
	tele           platform.Telegram
}

func InitProductModule(pStorage storage.ProductStorage, tele platform.Telegram) module.ProductModule {

	return &productModule{
		productStorage: pStorage,
		tele:           tele,
	}

}

func (p *productModule) CreateProduct(c telebot.Context, product dto.Product) error {

	if err := product.Validate(); err != nil {
		log.Print("validation error::", err)
		return err
	}

	err, id := p.productStorage.CreateProduct(context.Background(), product)
	if err != nil {
		return err
	}

	product.ID = id.String()

	if err := p.tele.AddOrderButtonToProduct(c, product); err != nil {

		log.Println("can't add button to product list ", err)

		return err
	}

	return nil

}
