package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
)

type cartCache struct {
	redis platform.Redis
}

func NewCartCache(redis platform.Redis) storage.CartStorage {
	return &cartCache{redis: redis}
}

func (c *cartCache) GetCart(ctx context.Context, userID int64) (map[string]int, error) {
	key := fmt.Sprintf("cart:%d", userID)
	resCmd, err := c.redis.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	res, err := resCmd.Result()
	if err != nil {
		return nil, err
	}

	cart := make(map[string]int)
	for k, v := range res {
		qty, _ := strconv.Atoi(v)
		cart[k] = qty
	}
	return cart, nil
}

func (c *cartCache) AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error {
	key := fmt.Sprintf("cart:%d", userID)
	field := fmt.Sprintf("p:%d", productID)
	
	if quantity <= 0 {
		return c.redis.HDel(ctx, key, field)
	}
	
	return c.redis.HSet(ctx, key, map[string]interface{}{field: quantity})
}

func (c *cartCache) ClearCart(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("cart:%d", userID)
	return c.redis.Delete(ctx, key)
}

