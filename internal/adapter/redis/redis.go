package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	idempotencyKeyPrefix = "Idempotency:"
	ttl                  = time.Minute * 5
)

type Redis struct {
	redis *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{redis: client}
}
