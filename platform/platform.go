package platform

import (
	"context"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"

	"github.com/go-redis/redis/v8"

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
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
	HSet(ctx context.Context,
		key string, values map[string]interface{}) error

	HGetAll(ctx context.Context, key string) (
		*redis.StringStringMapCmd, error)

	HExists(ctx context.Context, key string, field string) (bool, error)

	HDel(ctx context.Context, key string, fields ...string) error

	HGet(ctx context.Context, key string, field string) (string, error)

	HKeys(ctx context.Context, key string) ([]string, error)

	// Add other methods as needed
	// Example: Get(key string) (string, error)
	// Example: Delete(key string) error
	// Example: Expire(key string, duration time.Duration) error
	// Example: Keys(pattern string) ([]string, error)
}
