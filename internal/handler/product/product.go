package product

import (
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/handler"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"gopkg.in/telebot.v4"
)

type productHandler struct {
	productModule module.ProductModule
}

func InitProductHandler(pmd module.ProductModule) handler.Product {

	return &productHandler{
		productModule: pmd,
	}

}

func (p *productHandler) CreateProduct(c telebot.Context) error {

	userState := make(map[int64]interface{})
	productData := make(map[int64]dto.Product)
	user_id := c.Sender().ID

	switch userState[user_id] {

	case "":

		userState[user_id] = "waiting_for_title"
		productData[user_id] = dto.Product{}

		return c.Send("please enter product title:")

	case "waiting_for_title":
		product := productData[user_id]
		product.Title = c.Text()
		userState[user_id] = "waiting_for_description"
		productData[user_id] = product
		return c.Send("please enter product description:")

	case "waiting_for_description":
		product := productData[user_id]
		product.Description = c.Text()
		userState[user_id] = "waiting_for_price"
		return c.Send("please enter product price")

	case "waiting_for_price":
		price, err := strconv.ParseFloat(c.Text(), 64)

		if err != nil {
			return err
		}

		product := productData[user_id]
		product.Price = price
		userState[user_id] = ""

	}

	err := p.productModule.CreateProduct(c, productData[user_id])

	if err != nil {

		c.Send(err.Error())
		return err
	}

	// reseting the state
	delete(userState, user_id)
	delete(productData, user_id)
	return nil
}
