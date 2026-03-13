package events

import (
	"github.com/google/uuid"

	"github.com/Xanaduxan/tasks-golang/internal/storage"
)

type DeliveryStatusUpdated struct {
	DeliveryID uuid.UUID
	UserID     uuid.UUID
	Status     storage.DeliveryStatus
}
