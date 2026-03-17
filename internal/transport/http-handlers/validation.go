package http_handlers

import (
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/service/tasks"
)

func validateTask(name string, deadline *time.Time) error {
	if name == "" {
		return tasks.ErrInvalidInput
	}

	if deadline != nil {
		if deadline.IsZero() {
			return tasks.ErrInvalidInput
		}

		now := time.Now()

		if deadline.Before(now) {
			return tasks.ErrInvalidInput
		}

	}

	return nil
}
