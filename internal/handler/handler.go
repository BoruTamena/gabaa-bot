package handler

import "gopkg.in/telebot.v4"

type User interface {
	CreateUser(c telebot.Context) error
}
type Product interface {
	StartProductCreation(c telebot.Context) error
	CreateProduct(c telebot.Context) error
	//GetProducts
	//GetProductsById
}
type Order interface {
	HandleOrder(c telebot.Context) error
	CreateOrder(c telebot.Context, productID string) error
	AddToCart(c telebot.Context, productId string) error
}
