package tasks

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	tasks *storage.TaskStorage
	users *storage.UserStorage
}

func NewService(tasks *storage.TaskStorage, users *storage.UserStorage) *Service {
	return &Service{tasks: tasks, users: users}
}

func (s *Service) getUser(id uuid.UUID) (storage.User, error) {
	if id == uuid.Nil {
		return storage.User{}, ErrInvalidInput
	}
	u, err := s.users.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.User{}, ErrNotFound
		}
		return storage.User{}, err
	}
	return u, nil
}

func (s *Service) getTask(taskID uuid.UUID) (storage.Task, error) {
	if taskID == uuid.Nil {
		return storage.Task{}, ErrInvalidInput
	}
	t, err := s.tasks.GetById(taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Task{}, ErrNotFound
		}
		return storage.Task{}, err
	}
	return t, nil
}

func (s *Service) getOwnedTask(userID, taskID uuid.UUID) (storage.Task, error) {
	u, err := s.getUser(userID)
	if err != nil {
		return storage.Task{}, err
	}
	t, err := s.getTask(taskID)
	if err != nil {
		return storage.Task{}, err
	}
	if t.UserID != u.ID {
		return storage.Task{}, ErrForbidden
	}
	return t, nil
}

func (s *Service) GetTask(id, taskID uuid.UUID) (storage.Task, error) {
	return s.getOwnedTask(id, taskID)
}

func (s *Service) CreateTask(id uuid.UUID, name string, deadline *time.Time) (uuid.UUID, error) {
	if name == "" {
		return uuid.Nil, ErrInvalidInput
	}
	u, err := s.getUser(id)
	if err != nil {
		return uuid.Nil, err
	}

	t := storage.Task{
		ID:       uuid.New(),
		Name:     name,
		Deadline: deadline,
		UserID:   u.ID,
	}
	if err := s.tasks.Create(t); err != nil {
		return uuid.Nil, err
	}
	return t.ID, nil
}

func (s *Service) DeleteTask(id, taskID uuid.UUID) error {
	_, err := s.getOwnedTask(id, taskID)
	if err != nil {
		return err
	}
	return s.tasks.DeleteByID(taskID)
}

func (s *Service) UpdateTask(id, taskID uuid.UUID, name string, deadline *time.Time) error {
	if name == "" {
		return ErrInvalidInput
	}
	t, err := s.getOwnedTask(id, taskID)
	if err != nil {
		return err
	}
	t.Name = name
	t.Deadline = deadline
	return s.tasks.Update(t)
}
