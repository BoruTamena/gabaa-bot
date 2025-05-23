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
	fmt.Println("bot started ")

	tg.bot.Start()

	fmt.Println("listining... ")
}

func (tg *telegram) GetBot() *telebot.Bot {
	return tg.bot
}

func (tg *telegram) Group() telebot.Group {
	return *tg.bot.Group()
}

// add order now inline button
func (tg *telegram) AddButtonToProduct(c telebot.Context, data dto.Product) error {
	inline := &telebot.ReplyMarkup{}

	addToCart := inline.Data(" 🛒 Add to cart", "cart/"+data.ID)
	orderNow := inline.Data("🛍️ Order Now", "order/"+data.ID)
	inline.Row(orderNow)
	inline.Inline(inline.Row(addToCart, orderNow))

	message := fmt.Sprintf("*Product name :* %s \n *Description: *%s \n *Price: * %v \n  --- \n powered by Gabaa Place",
		data.Title, data.Description, data.Price)

	return c.Send(message, inline, telebot.ModeMarkdown)

}
