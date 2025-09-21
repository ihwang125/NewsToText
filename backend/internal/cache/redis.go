package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	SetJSON(key string, value interface{}, expiration time.Duration) error
	GetJSON(key string, dest interface{}) error
}

type redisCache struct {
	client *redis.Client
}

func Initialize(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewCache(client *redis.Client) Cache {
	return &redisCache{client: client}
}

func (c *redisCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisCache) Get(key string) (string, error) {
	ctx := context.Background()
	return c.client.Get(ctx, key).Result()
}

func (c *redisCache) Delete(key string) error {
	ctx := context.Background()
	return c.client.Del(ctx, key).Err()
}

func (c *redisCache) Exists(key string) (bool, error) {
	ctx := context.Background()
	count, err := c.client.Exists(ctx, key).Result()
	return count > 0, err
}

func (c *redisCache) SetJSON(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(key, data, expiration)
}

func (c *redisCache) GetJSON(key string, dest interface{}) error {
	data, err := c.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}