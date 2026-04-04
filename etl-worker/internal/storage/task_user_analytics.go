package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TaskUserAnalyticsStorage struct {
	db *sql.DB
}

func NewTaskUserAnalyticsStorage(db *sql.DB) *TaskUserAnalyticsStorage {
	return &TaskUserAnalyticsStorage{db: db}
}

func (s *TaskUserAnalyticsStorage) IncrementCompleted(userID uuid.UUID, eventTime time.Time) error {
	_, err := s.db.Exec(`
		INSERT INTO task_user_analytics (
			user_id,
			tasks_completed,
			last_event_at
		)
		VALUES ($1, 1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET
			tasks_completed = task_user_analytics.tasks_completed + 1,
			last_event_at = EXCLUDED.last_event_at
	`, userID, eventTime)

	return err
}
