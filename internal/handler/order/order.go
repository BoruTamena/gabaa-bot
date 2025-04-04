package order

import (
	"fmt"
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/handler"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"gopkg.in/telebot.v4"
)

type orderHandler struct {
	orderModule module.OrderModule
}

func InitOrderHandler(orderModule module.OrderModule) handler.Order {

	return &orderHandler{
		orderModule: orderModule,
	}
}

func (o *orderHandler) HandleOrder(c telebot.Context) error {
	data := c.Callback().Data

	parts := strings.Split(data, "/")
	if len(parts) != 2 {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Invalid command!", ShowAlert: true})
	}

	action := strings.TrimSpace(parts[0])
	productID := strings.TrimSpace(parts[1])

	fmt.Println("Action:", action, "Product ID:", productID)

	switch action {
	case "order":
		return c.Send(fmt.Sprintf("✅ Order placed for Product ID: %s", productID))
	case "cart":
		return c.Send(fmt.Sprintf("🛒 Added to cart: Product ID: %s", productID))

	}

	return c.Respond(&telebot.CallbackResponse{Text: "Unknown action!", ShowAlert: true})

}

func (o *orderHandler) AddToCart(c telebot.Context, productId string) error {

	if err := o.orderModule.AddToCart(c, productId); err != nil {
	}

	return nil

}
