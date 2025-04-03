package product

import (
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing"
	"github.com/BoruTamena/gabaa-bot/internal/handler"
	"gopkg.in/telebot.v4"
)

func InitProductRoute(group *telebot.Group, handler handler.Product) {

	routes := []routing.Router{
		{
			Path:    "/start",
			Handler: handler.StartProductCreation,
		},
		{
			Path:    telebot.OnText,
			Handler: handler.CreateProduct,
		},
	}

	routing.Register(group, routes)

}
