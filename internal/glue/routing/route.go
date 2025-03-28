package routing

import (
	tele "gopkg.in/telebot.v4"
)

type Router struct {
	Path    string
	Handler tele.HandlerFunc
	// Middlewares []tele.MiddlewareFunc
}

func Register(group *tele.Group, routes []Router) {

	for _, route := range routes {
		// group.Handle(route.Path, route.Handler, route.Middlewares...)
		group.Handle(route.Path, route.Handler)

	}

}
