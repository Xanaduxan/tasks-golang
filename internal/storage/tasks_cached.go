package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	taskPrefix = "task:"
	taskTTL    = 5 * time.Minute
)

type TaskCached struct {
	repo  *TaskStorage
	redis *redis.Client
}

func NewTaskCached(repo *TaskStorage, redisClient *redis.Client) *TaskCached {
	return &TaskCached{
		repo:  repo,
		redis: redisClient,
	}
}

func (r *TaskCached) cacheKey(id uuid.UUID) string {
	return taskPrefix + id.String()
}

func (r *TaskCached) setCache(ctx context.Context, task Task) error {
	key := r.cacheKey(task.ID)

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task: %w", err)
	}

	if err := r.redis.Set(ctx, key, data, taskTTL).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

func (r *TaskCached) deleteCache(ctx context.Context, id uuid.UUID) error {
	key := r.cacheKey(id)

	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}

	return nil
}

func (r *TaskCached) Create(task Task) error {
	task.Status = StatusCreated

	if err := r.repo.Create(task); err != nil {
		return err
	}

	if err := r.setCache(context.Background(), task); err != nil {
		log.Printf("cache create error: %v", err)
	}

	return nil
}

func (r *TaskCached) GetByID(id uuid.UUID) (Task, error) {
	ctx := context.Background()
	key := r.cacheKey(id)

	data, err := r.redis.Get(ctx, key).Bytes()
	if err == nil {
		var task Task
		if err := json.Unmarshal(data, &task); err == nil {
			return task, nil
		}
		log.Printf("redis unmarshal error: %v", err)
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("redis get error: %v", err)
	}

	task, err := r.repo.GetByID(id)
	if err != nil {
		return Task{}, err
	}

	if err := r.setCache(ctx, task); err != nil {
		log.Printf("cache refill error: %v", err)
	}

	return task, nil
}

func (r *TaskCached) Update(task Task) error {
	if err := r.repo.Update(task); err != nil {
		return err
	}

	if err := r.setCache(context.Background(), task); err != nil {
		log.Printf("cache update error: %v", err)
	}

	return nil
}

func (r *TaskCached) DeleteByID(id uuid.UUID) error {
	if err := r.repo.DeleteByID(id); err != nil {
		return err
	}

	if err := r.deleteCache(context.Background(), id); err != nil {
		log.Printf("cache delete error: %v", err)
	}

	return nil
}

func (r *TaskCached) HasAccess(taskID, userID uuid.UUID) (bool, error) {
	return r.repo.HasAccess(taskID, userID)
}

func (r *TaskCached) UpdateStatus(id uuid.UUID, status TaskStatus) error {
	if err := r.repo.UpdateStatus(id, status); err != nil {
		return err
	}

	if err := r.deleteCache(context.Background(), id); err != nil {
		log.Printf("cache delete after status update error: %v", err)
	}

	return nil
}

func (r *TaskCached) GetAllNotDone() ([]Task, error) {
	return r.repo.GetAllNotDone()
}

func (r *TaskCached) Count() (int, error) {
	return r.repo.Count()
}
