package service

import (
	"log"
	"time"

	"github.com/Xanaduxan/tasks-golang/etl-worker/internal/events"
	"github.com/google/uuid"
)

type EventStorage interface {
	Save(event events.TaskEvent) error
}

type UserAnalyticsStorage interface {
	IncrementCompleted(userID uuid.UUID, eventTime time.Time) error
}

type Processor struct {
	storage   EventStorage
	analytics UserAnalyticsStorage
}

func NewProcessor(storage EventStorage, analytics UserAnalyticsStorage) *Processor {
	return &Processor{
		storage:   storage,
		analytics: analytics,
	}
}

func (p *Processor) Process(event events.TaskEvent) error {
	if err := p.storage.Save(event); err != nil {
		return err
	}

	if event.EventType == events.TaskCompletedEvent {
		if err := p.analytics.IncrementCompleted(event.UserID, event.Timestamp); err != nil {
			return err
		}
	}

	log.Printf(
		"processed event: type=%s task_id=%s user_id=%s status=%s",
		event.EventType,
		event.TaskID.String(),
		event.UserID.String(),
		event.Status,
	)

	return nil
}
