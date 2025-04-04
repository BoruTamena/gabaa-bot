package rediscache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	rdb *redis.Client
}

func NewRedis(redisClient *RedisClient) *redisCache {

	if redisClient == nil {
		return nil
	}
	// Ensure the Redis client is properly initialized
	if redisClient.Addr == "" {
		return nil
	}
	if redisClient.DB < 0 {
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisClient.Addr,
		Password: redisClient.Password,
		DB:       redisClient.DB,
	})

	// Return the Redis instance
	return &redisCache{
		rdb: client,
	}
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}) error {
	err := r.rdb.Set(ctx, key, value,
		time.Duration(48*time.Hour)).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
func (r *redisCache) Delete(ctx context.Context, key string) error {
	err := r.rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.rdb.Expire(ctx, key, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	val, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (r *redisCache) HSet(ctx context.Context,
	key string, values map[string]interface{}) error {
	err := r.rdb.HSet(ctx, key, values).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) HGetAll(ctx context.Context, key string) (
	*redis.StringStringMapCmd, error) {

	val := r.rdb.HGetAll(ctx, key)
	return val, nil
}

func (r *redisCache) HExists(ctx context.Context, key string, field string) (bool, error) {
	val, err := r.rdb.HExists(ctx, key, field).Result()
	if err != nil {
		return false, err
	}
	return val, nil
}
func (r *redisCache) HDel(ctx context.Context, key string, fields ...string) error {
	err := r.rdb.HDel(ctx, key, fields...).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *redisCache) HGet(ctx context.Context, key string, field string) (string, error) {
	val, err := r.rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
func (r *redisCache) HKeys(ctx context.Context, key string) ([]string, error) {
	val, err := r.rdb.HKeys(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}
