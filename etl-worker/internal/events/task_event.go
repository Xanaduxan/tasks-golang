package events

import (
	"time"

	"github.com/google/uuid"
)

type TaskEventType string

const (
	TaskCreatedEvent   TaskEventType = "task_created"
	TaskUpdatedEvent   TaskEventType = "task_updated"
	TaskCompletedEvent TaskEventType = "task_completed"
	TaskStatusEvent    TaskEventType = "task_status_updated"
)

const TaskEventsTopic = "task_events"

type TaskEvent struct {
	EventID        uuid.UUID     `json:"event_id"`
	TaskID         uuid.UUID     `json:"task_id"`
	UserID         uuid.UUID     `json:"user_id"`
	GroupID        *uuid.UUID    `json:"group_id,omitempty"`
	EventType      TaskEventType `json:"event_type"`
	Status         string        `json:"status"`
	PreviousStatus *string       `json:"previous_status,omitempty"`
	Timestamp      time.Time     `json:"timestamp"`
}
