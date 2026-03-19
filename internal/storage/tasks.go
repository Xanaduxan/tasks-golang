package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusCreated    TaskStatus = "created"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID       uuid.UUID
	Name     string
	Deadline *time.Time
	UserID   uuid.UUID
	GroupID  *uuid.UUID
	Status   TaskStatus
}

type TaskStorage struct {
	DB *sql.DB
}

func NewTaskStorage(db *sql.DB) *TaskStorage {
	return &TaskStorage{DB: db}
}

func (s *TaskStorage) Create(task Task) error {
	task.Status = StatusCreated
	_, err := s.DB.Exec(`
		INSERT INTO tasks (id, name, deadline, user_id, group_id, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, task.ID, task.Name, task.Deadline, task.UserID, task.GroupID, task.Status)

	return err
}

func (s *TaskStorage) GetByID(id uuid.UUID) (Task, error) {
	var task Task

	err := s.DB.QueryRow(`
		SELECT id, name, deadline, user_id, group_id, status
		FROM tasks
		WHERE id = $1
	`, id).Scan(&task.ID, &task.Name, &task.Deadline, &task.UserID, &task.GroupID, &task.Status)

	return task, err
}

func (s *TaskStorage) Update(task Task) error {
	_, err := s.DB.Exec(`
		UPDATE tasks 
		SET name = $2, deadline = $3, group_id = $4, status = $5
		WHERE id = $1
	`, task.ID, task.Name, task.Deadline, task.GroupID, task.Status)

	return err
}

func (s *TaskStorage) DeleteByID(id uuid.UUID) error {

	_, err := s.DB.Exec(`
		DELETE FROM tasks
		WHERE id = $1
	`, id)

	return err
}

func (s *TaskStorage) HasAccess(taskID, userID uuid.UUID) (bool, error) {
	var exists bool

	err := s.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM tasks t
			LEFT JOIN group_members gm ON gm.group_id = t.group_id
			WHERE t.id = $1
			  AND (t.user_id = $2 OR gm.user_id = $2)
		)
	`, taskID, userID).Scan(&exists)

	return exists, err
}
func (s *TaskStorage) UpdateStatus(id uuid.UUID, status TaskStatus) error {
	_, err := s.DB.Exec(`
		UPDATE tasks
		SET status = $2
		WHERE id = $1
	`, id, status)

	return err
}
func (s *TaskStorage) GetAllNotDone() ([]Task, error) {
	rows, err := s.DB.Query(`
		SELECT id, name, deadline, user_id, group_id, status
		FROM tasks
		WHERE status <> $1
	`, StatusDone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task

		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Deadline,
			&t.UserID,
			&t.GroupID,
			&t.Status,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
func (s *TaskStorage) Count() (int, error) {
	var count int

	err := s.DB.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
