package tasks

import "github.com/Xanaduxan/tasks-golang/internal/events"

type Notifier interface {
	NotifyTaskStatusUpdated(event events.TaskStatusUpdated) error
}
