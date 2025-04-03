package order

import (
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing"
	"github.com/BoruTamena/gabaa-bot/internal/handler"
	"gopkg.in/telebot.v4"
)

func InitOrderRoute(group *telebot.Group, handler handler.Order) {

	routes := []routing.Router{
		{
			Path:    telebot.OnCallback,
			Handler: handler.HandleOrder,
		},
	}

	routing.Register(group, routes)

}
