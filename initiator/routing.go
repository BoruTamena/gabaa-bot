package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing/order"
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing/product"
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing/user"
	"gopkg.in/telebot.v4"
)

func InitRoute(group *telebot.Group, handler Handler) {
	// intalizing user registeration bot handler
	user.InitUserRoute(group, handler.userHandler)
	// initalizing product bot handler
	product.InitProductRoute(group, handler.productHandler)
	// intalizing orders bot handler
	order.InitOrderRoute(group, handler.orderHandler)
}
