package storage

import (
	"database/sql"

	"github.com/Xanaduxan/tasks-golang/etl-worker/internal/events"
)

type TaskEventLogStorage struct {
	db *sql.DB
}

func NewTaskEventLogStorage(db *sql.DB) *TaskEventLogStorage {
	return &TaskEventLogStorage{db: db}
}

func (s *TaskEventLogStorage) Save(event events.TaskEvent) error {
	_, err := s.db.Exec(`
		INSERT INTO task_event_log (
			event_id,
			task_id,
			user_id,
			group_id,
			event_type,
			status,
			previous_status,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (event_id) DO NOTHING
	`,
		event.EventID,
		event.TaskID,
		event.UserID,
		event.GroupID,
		event.EventType,
		event.Status,
		event.PreviousStatus,
		event.Timestamp,
	)

	return err
}
