package redis

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	taskPrefix = "task:"
	taskTTL    = 5 * time.Minute
)

func (r *Redis) GetTask(ctx context.Context, id uuid.UUID) (*storage.Task, bool) {
	key := taskPrefix + id.String()

	data, err := r.redis.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false
		}

		log.Printf("redis GetTask error: %v", err)
		return nil, false
	}

	var task storage.Task
	if err := json.Unmarshal(data, &task); err != nil {
		log.Printf("redis unmarshal error: %v", err)
		return nil, false
	}

	return &task, true
}

func (r *Redis) SetTask(ctx context.Context, task storage.Task) {
	key := taskPrefix + task.ID.String()

	data, err := json.Marshal(task)
	if err != nil {
		log.Printf("redis marshal error: %v", err)
		return
	}

	err = r.redis.Set(ctx, key, data, taskTTL).Err()
	if err != nil {
		log.Printf("redis SetTask error: %v", err)
	}
}

func (r *Redis) DeleteTask(ctx context.Context, id uuid.UUID) {
	key := taskPrefix + id.String()

	if err := r.redis.Del(ctx, key).Err(); err != nil {
		log.Printf("redis DeleteTask error: %v", err)
	}
}
