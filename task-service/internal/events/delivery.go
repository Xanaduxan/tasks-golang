package events

import (
	"github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	"github.com/google/uuid"
)

type DeliveryStatusUpdated struct {
	DeliveryID uuid.UUID
	UserID     uuid.UUID
	Status     storage.DeliveryStatus
}
