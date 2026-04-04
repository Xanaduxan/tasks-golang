package tasks

import (
	"context"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/events"
)

type EventPublisher interface {
	PublishTaskEvent(ctx context.Context, event events.TaskEvent) error
}
