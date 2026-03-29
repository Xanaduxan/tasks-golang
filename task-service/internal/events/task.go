package events

import (
	"github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	"github.com/google/uuid"
)

type TaskStatusUpdated struct {
	TaskID  uuid.UUID
	GroupID *uuid.UUID
	UserID  uuid.UUID
	Status  storage.TaskStatus
}
