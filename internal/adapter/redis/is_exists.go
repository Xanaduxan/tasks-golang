package redis

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
)

func (r *Redis) IsExists(ctx context.Context, idempotencyKey string) bool {
	key := idempotencyKeyPrefix + idempotencyKey

	err := r.redis.Get(ctx, key).Err()
	if err == nil {
		return true
	}

	if !errors.Is(err, redis.Nil) {
		log.Printf("redis Get failed: %v", err)
	}

	err = r.redis.Set(ctx, key, []byte{}, ttl).Err()
	if err != nil {
		log.Printf("redis Set failed: %v", err)
	}

	return false
}
