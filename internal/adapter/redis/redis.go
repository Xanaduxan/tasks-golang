package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	idempotencyKeyPrefix = "Idempotency:"
	ttl                  = time.Hour
)

type Redis struct {
	redis *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{redis: client}
}
