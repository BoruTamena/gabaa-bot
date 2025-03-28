package platform

import (
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"gopkg.in/telebot.v4"
)

// define your platform interfaces here

type Telegram interface {
	Start()
	Group() telebot.Group
	AddOrderButtonToProduct(c telebot.Context, data dto.Product) error
}
