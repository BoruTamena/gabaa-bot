package routing

import (
	tele "gopkg.in/telebot.v4"
)

type Router struct {
	path        string
	Handler     tele.HandlerFunc
	Middlewares []tele.HandlerFunc
}

func (r *Router) Register(bot *tele.Bot) {
	group := bot.Group()
	// for _, middleware := range r.Middlewares {
	// 	group.Use(middleware)
	// }
	group.Handle(r.path, r.Handler)
}
