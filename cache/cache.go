package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type client struct {
	client *redis.Client
}

func NewClient(ctx context.Context, address string) (Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: address,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client{redisClient}, nil
}

func (c client) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	return value, nil
}

func (c client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	status := c.client.Set(ctx, key, value, expiration)
	err := status.Err()
	if err != nil {
		return err
	}

	return nil
}
