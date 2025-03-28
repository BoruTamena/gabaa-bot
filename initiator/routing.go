package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing/product"
	"gopkg.in/telebot.v4"
)

func InitRoute(group *telebot.Group, handler Handler) {

	// initalizing product bot handler
	product.InitProductRoute(group, handler.productHandler)

}
