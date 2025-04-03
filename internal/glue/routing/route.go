package routing

import (
	tele "gopkg.in/telebot.v4"
)

type Router struct {
	Path    string
	Handler tele.HandlerFunc
	// Middlewares []tele.MiddlewareFunc
}
type CallbackRouter struct {
	Btn     tele.Btn
	Handler tele.HandlerFunc
	// Middlewares []tele.MiddlewareFunc
}

func Register(group *tele.Group, routes []Router) {

	for _, route := range routes {
		// group.Handle(route.Path, route.Handler, route.Middlewares...)
		group.Handle(route.Path, route.Handler)

	}
}

func RegisterCallback(group *tele.Group, routes []CallbackRouter) {

	for _, route := range routes {
		// group.Handle(&route.Btn, route.Handler, route.Middlewares...)
		group.Handle(&route.Btn, route.Handler)
	}

}

func BindCallback(group *tele.Bot, routes []Router) {

	for _, route := range routes {
		// Bind the handler to the bot's callback event
		group.Handle(&tele.Callback{Data: route.Path}, route.Handler)
	}

}
