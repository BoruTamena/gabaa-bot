package order

import (
	"context"
	"fmt"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/errors"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
)

type orderModule struct {
	productStorage storage.ProductStorage
	cache          platform.Redis
}

func InitOrderModule(pStorage storage.ProductStorage, cache platform.Redis) module.OrderModule {

	return &orderModule{
		productStorage: pStorage,
		cache:          cache,
	}

}

func (order *orderModule) AddToCart(cxt context.Context, user_id, productId string) error {

	// Check if the product exists in the storage
	product, err := order.productStorage.GetProductByID(cxt, productId)
	if err != nil {

		return err
	}

	// Check if the product is already in the cart
	exists, err := order.cache.HExists(cxt, user_id, product.ID)
	if err != nil {
		return err

	}
	if exists {
		err := errors.CartItemAlreadyExistsErr.Wrap(fmt.
			Errorf("product %v already exists in the cart",
				product.ID), "product already exists in the cart")

		return err
	}

	// Add the product to the cart
	err = order.cache.HSet(cxt, user_id, map[string]interface{}{
		product.ID: product.ID,
		"qty":      1,
	})

	if err != nil {
		return err
	}

	// Set the expiration time for the cart
	// This will remove the cart after 24 hours of inactivity
	// This is to prevent the cart from growing indefinitely
	// and to free up memory
	order.cache.Expire(cxt, user_id, time.Duration(24)*time.Hour)

	return nil

}
func (order *orderModule) CreateOrder() error {

	return nil

}
