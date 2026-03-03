package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID       uuid.UUID
	Name     string
	Deadline *time.Time
	UserID   uuid.UUID
}

type TaskStorage struct {
	DB *sql.DB
}

func NewTaskStorage(db *sql.DB) *TaskStorage {
	return &TaskStorage{DB: db}
}

func (s *TaskStorage) Create(task Task) error {
	_, err := s.DB.Exec(`
		INSERT INTO tasks (id, name, deadline, user_id)
		VALUES ($1, $2, $3, $4)
	`, task.ID, task.Name, task.Deadline, task.UserID)

	return err
}

func (s *TaskStorage) GetById(id uuid.UUID) (Task, error) {
	var task Task

	err := s.DB.QueryRow(`
		SELECT id, name, deadline, user_id
		FROM tasks
		WHERE id = $1
	`, id).Scan(&task.ID, &task.Name, &task.Deadline, &task.UserID)

	return task, err
}

func (s *TaskStorage) Update(task Task) error {
	_, err := s.DB.Exec(`
		UPDATE tasks 
		SET name = $2, deadline = $3
		WHERE id = $1
	`, task.ID, task.Name, task.Deadline)

	return err
}

func (s *TaskStorage) DeleteByID(id uuid.UUID) error {

	_, err := s.DB.Exec(`
		DELETE FROM tasks
		WHERE id = $1
	`, id)

	return err
}
