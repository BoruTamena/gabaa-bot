package handler

import "gopkg.in/telebot.v4"

// define your handlers interface here

type Product interface {
	CreateProduct(c telebot.Context)
}
