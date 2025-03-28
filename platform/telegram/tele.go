package telegram

import (
	"fmt"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/spf13/viper"
	"gopkg.in/telebot.v4"
)

type telegram struct {
	bot *telebot.Bot
}

func InitTelBot() platform.Telegram {

	setting := telebot.Settings{

		Token:  viper.GetString("tg.token"),
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

func (tg *telegram) Start() {
	tg.bot.Start()
}

func (tg *telegram) Group() telebot.Group {

	return *tg.bot.Group()

}

// add order now inline button
func (tg *telegram) AddOrderButtonToProduct(c telebot.Context, data dto.Product) error {
	inline := &telebot.ReplyMarkup{}

	btn := inline.Data(" 🛒 Order Now", data.ID)
	inline.Inline(inline.Row(btn))

	message := fmt.Sprintf("*Product name :* %s \n *Description: *%s \n *Price: * %d \n  --- \n powered by Gabaa Place",
		data.Title, data.Description, data.Price)

	return c.Send(message, inline, telebot.ModeMarkdown)

}
