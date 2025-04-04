package rediscache

import (
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

func (r *redisCache) Set(key string, value interface{}) error {
	err := r.rdb.Set(r.rdb.Context(), key, value,
		time.Duration(48*time.Hour)).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *redisCache) Get(key string) (string, error) {
	val, err := r.rdb.Get(r.rdb.Context(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
func (r *redisCache) Delete(key string) error {
	err := r.rdb.Del(r.rdb.Context(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Expire(key string, expiration time.Duration) error {
	err := r.rdb.Expire(r.rdb.Context(), key, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
