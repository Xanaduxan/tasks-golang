package websocket

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type DeliveryStatusUpdatedData struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
	Status     string    `json:"status"`
}
type TaskStatusUpdatedData struct {
	TaskID uuid.UUID `json:"task_id"`
	Status string    `json:"status"`
}
