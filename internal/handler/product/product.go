package product

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/handler"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"gopkg.in/telebot.v4"
)

var (
	userState   = make(map[int64]string)
	productData = make(map[int64]dto.Product)
	mu          sync.Mutex // Ensures thread safety
)

type productHandler struct {
	productModule module.ProductModule
}

func InitProductHandler(pmd module.ProductModule) handler.Product {

	return &productHandler{
		productModule: pmd,
	}

}

func (p *productHandler) StartProductCreation(c telebot.Context) error {

	userID := c.Sender().ID
	mu.Lock()
	defer mu.Unlock()
	userState[userID] = "waiting_for_title"
	productData[userID] = dto.Product{}
	return c.Send("📝 Please enter the product title:")

}

func (p *productHandler) CreateProduct(c telebot.Context) error {
	userID := c.Sender().ID

	mu.Lock()
	defer mu.Unlock()

	switch userState[userID] {

	case "waiting_for_title":
		product := productData[userID]
		product.Title = c.Text()
		productData[userID] = product
		userState[userID] = "waiting_for_description"
		return c.Send("📄 Please enter the product description:")

	case "waiting_for_description":
		product := productData[userID]
		product.Description = c.Text()
		productData[userID] = product
		userState[userID] = "waiting_for_price"
		return c.Send("💰 Please enter the product price (numeric value):")

	case "waiting_for_price":

		price, err := strconv.ParseFloat(c.Text(), 64)
		if err != nil {
			return c.Send("❌ Invalid price format. Please enter a valid number:" + c.Text())
		}

		product := productData[userID]
		product.Price = price
		product.SellerId = userID
		productData[userID] = product

		fmt.Println("product", productData[userID])
		// Save the product in the database
		if err := p.productModule.CreateProduct(c, productData[userID]); err != nil {
			return c.Send("❌ Failed to create product: " + err.Error())
		}
		// Reset user state after success
		delete(userState, userID)
		delete(productData, userID)

		return c.Send("✅ Product successfully created!")

	}

	return nil
}
