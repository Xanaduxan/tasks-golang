package events

import (
	"github.com/google/uuid"

	"github.com/Xanaduxan/tasks-golang/internal/storage"
)

type TaskStatusUpdated struct {
	TaskID  uuid.UUID
	GroupID *uuid.UUID
	UserID  uuid.UUID
	Status  storage.TaskStatus
}
