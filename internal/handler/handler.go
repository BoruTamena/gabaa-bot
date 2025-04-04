package handler

import "gopkg.in/telebot.v4"

type User interface {
	CreateUser(c telebot.Context) error
}
type Product interface {
	StartProductCreation(c telebot.Context) error
	CreateProduct(c telebot.Context) error
}

type Order interface {
	HandleOrder(c telebot.Context) error
}
