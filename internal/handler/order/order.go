package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/handler"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/google/uuid"
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
		return o.CreateOrder(c, productID)
	case "cart":
		return o.AddToCart(c, productID)

	}

	return c.Respond(&telebot.CallbackResponse{Text: "Unknown action!", ShowAlert: true})

}

// AddToCart handles the addition of a product to the cart
// It takes the context and productID as parameters
// and returns an error if any occurs during the process
// The function should check if the product exists in the storage
// and if the product can be added to the cart successfully
// It should also handle any errors that may occur
// during the addition process
// It should return nil if the product is added successfully
func (o *orderHandler) AddToCart(c telebot.Context, productId string) error {

	user_id := c.Sender().ID

	err := o.orderModule.AddToCart(context.Background(), fmt.Sprint(user_id), productId)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.Respond(&telebot.CallbackResponse{
				Text:      "Product already exists in the cart",
				ShowAlert: true,
			})
		}

		return c.Respond(&telebot.CallbackResponse{
			Text:      "Sorry we could not add the product to the cart",
			ShowAlert: true,
		})
	}
	return c.Respond(&telebot.CallbackResponse{
		Text:      "Product added to cart successfully !",
		ShowAlert: true,
	})

}

// CreateOrder handles the order creation process
// It should be called when the user confirms the order
// It takes the context and productID as parameters
// and returns an error if any occurs during the process
// The function should check if the product exists in the storage
// and if the order can be created successfully
// It should also handle any errors that may occur
// during the order creation process
// It should return nil if the order is created successfully
// It should also handle any errors that may occur
// during the order creation process
func (o *orderHandler) CreateOrder(c telebot.Context, productID string) error {
	// Implement order creation logic here

	productId, _ := uuid.Parse(productID)

	err := o.orderModule.CreateOrder(context.Background(), dto.Order{
		BuyerID: c.Sender().ID,
		Status:  "pending",
		OrderItems: []dto.OrderItem{
			{
				ProductId: productId,
				Price:     32,
				Quantity:  1,
			},
		},
	})

	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text:      "Sorry we could not create the order",
			ShowAlert: true,
		})
	}
	return c.Respond(&telebot.CallbackResponse{
		Text:      "Order created successfully!",
		ShowAlert: true,
	})
}
