package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
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

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	setting := telebot.Settings{

		Token:  viper.GetString("tg.token"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Client: client,
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

	idStr := strconv.FormatInt(data.ID, 10)
	addToCart := inline.Data(" 🛒 Add to cart", "cart/"+idStr)
	orderNow := inline.Data("🛍️ Order Now", "order/"+idStr)
	inline.Row(orderNow)
	inline.Inline(inline.Row(addToCart, orderNow))

	message := fmt.Sprintf("*Product name :* %s \n *Description: *%s \n *Price: * %v \n  --- \n powered by Gabaa Place",
		data.Name, data.Description, data.Price)

	return c.Send(message, inline, telebot.ModeMarkdown)
}

func (tg *telegram) ValidateInitData(initData string) (bool, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return false, err
	}

	hash := values.Get("hash")
	if hash == "" {
		return false, fmt.Errorf("hash missing")
	}

	var keys []string
	for k := range values {
		if k != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var dataCheckString strings.Builder
	for i, k := range keys {
		if i > 0 {
			dataCheckString.WriteString("\n")
		}
		dataCheckString.WriteString(k)
		dataCheckString.WriteString("=")
		dataCheckString.WriteString(values.Get(k))
	}

	secretKey := hmacSHA256([]byte(tg.bot.Token), []byte("WebAppData"))
	calculatedHash := hex.EncodeToString(hmacSHA256(secretKey, []byte(dataCheckString.String())))

	return calculatedHash == hash, nil
}

func (tg *telegram) ParseInitData(initData string) (*dto.TelegramUser, int64, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, 0, err
	}

	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, 0, fmt.Errorf("user data missing in initData")
	}

	var user dto.TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, 0, err
	}

	chatIDStr := values.Get("chat_id")
	var chatID int64
	if chatIDStr != "" {
		chatID, _ = strconv.ParseInt(chatIDStr, 10, 64)
	}

	return &user, chatID, nil
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func (tg *telegram) IsChatAdmin(chatID int64, userID int64) (bool, error) {
	admins, err := tg.bot.AdminsOf(&telebot.Chat{ID: chatID})
	if err != nil {
		return false, err
	}

	for _, admin := range admins {
		if admin.User.ID == userID {
			return true, nil
		}
	}
	return false, nil
}
