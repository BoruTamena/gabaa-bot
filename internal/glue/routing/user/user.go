package user

import (
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing"
	"github.com/BoruTamena/gabaa-bot/internal/handler"
	"gopkg.in/telebot.v4"
)

func InitUserRoute(group *telebot.Group, handler handler.User) {

	routes := []routing.Router{
		{
			Path:    "/register",
			Handler: handler.CreateUser,
		},
	}

	routing.Register(group, routes)

}
