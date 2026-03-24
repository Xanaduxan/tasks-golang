package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis(addr string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func Close(client *redis.Client) error {
	if client == nil {
		return nil
	}

	return client.Close()
}
