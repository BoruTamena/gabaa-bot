package telegram

import (
	"time"

	"github.com/BoruTamena/gabaa-bot/platform"
	"gopkg.in/telebot.v4"
)

type telegram struct {
	bot *telebot.Bot
}

func InitTelBot() platform.Telegram {

	setting := telebot.Settings{

		// TODO  add token
		// token:
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(setting)

	if err != nil {
		panic(err)
	}

	return &telegram{
		bot: bot,
	}
}
