package platform

import (
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"gopkg.in/telebot.v4"
)

// define your platform interfaces here

type Telegram interface {
	Start()
	GetBot() *telebot.Bot
	Group() telebot.Group
	AddButtonToProduct(c telebot.Context, data dto.Product) error
}

type Redis interface {
	// define your Redis methods here
	// Example: Set(key string, value interface{}) error
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Delete(key string) error
	// Add other methods as needed
	// Example: Get(key string) (string, error)
	// Example: Delete(key string) error
	// Example: Increment(key string) (int, error)
	// Example: Decrement(key string) (int, error)
	// Example: Exists(key string) (bool, error)
	// Example: Expire(key string, duration time.Duration) error
	// Example: FlushAll() error
	// Example: Keys(pattern string) ([]string, error)
}
