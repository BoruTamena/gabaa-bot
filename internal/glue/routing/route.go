package routing

import (
	tele "gopkg.in/telebot.v4"
)

type Router struct {
	path        string
	Handler     tele.HandlerFunc
	Middlewares []tele.HandlerFunc
}

func Register(bot *tele.Bot, routes []Router) {
	group := bot.Group()
	// for _, middleware := range r.Middlewares {
	// 	group.Use(middleware)
	// }
	for _, route := range routes {
		group.Handle(route.path, route.Handler)
	}

}
