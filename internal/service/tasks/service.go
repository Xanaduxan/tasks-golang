package tasks

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/events"
	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/Xanaduxan/tasks-golang/metrics"
	"github.com/google/uuid"
)

type tasksInterface interface {
	Create(task storage.Task) error
	GetByID(id uuid.UUID) (storage.Task, error)
	Update(task storage.Task) error
	DeleteByID(id uuid.UUID) error
	HasAccess(taskID, userID uuid.UUID) (bool, error)
	UpdateStatus(id uuid.UUID, status storage.TaskStatus) error
	GetAllNotDone() ([]storage.Task, error)
	Count() (int, error)
}
type Service struct {
	tasks        tasksInterface
	users        *storage.UserStorage
	groups       *storage.GroupStorage
	groupMembers *storage.GroupMemberStorage
	notifier     Notifier
}

func NewService(
	tasks tasksInterface,
	users *storage.UserStorage,
	groups *storage.GroupStorage,
	groupMembers *storage.GroupMemberStorage,
	notifier Notifier,

) *Service {
	return &Service{
		tasks:        tasks,
		users:        users,
		groups:       groups,
		groupMembers: groupMembers,
		notifier:     notifier,
	}
}

func (s *Service) getUser(id uuid.UUID) (storage.User, error) {
	if id == uuid.Nil {
		return storage.User{}, ErrInvalidInput
	}

	u, err := s.users.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.User{}, ErrNotFound
		}
		return storage.User{}, err
	}

	return u, nil
}

func (s *Service) getGroup(id *uuid.UUID) (*storage.Group, error) {
	if id == nil {
		return nil, nil
	}
	if *id == uuid.Nil {
		return nil, ErrInvalidInput
	}

	g, err := s.groups.GetByID(*id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &g, nil
}

func (s *Service) ensureUserInGroup(userID uuid.UUID, groupID *uuid.UUID) error {
	if groupID == nil {
		return nil
	}

	isMember, err := s.groupMembers.IsMember(*groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrForbidden
	}

	return nil
}

func (s *Service) getTask(taskID uuid.UUID) (storage.Task, error) {
	if taskID == uuid.Nil {
		return storage.Task{}, ErrInvalidInput
	}

	t, err := s.tasks.GetByID(taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Task{}, ErrNotFound
		}
		return storage.Task{}, err
	}

	return t, nil
}

func (s *Service) getOwnedTask(userID, taskID uuid.UUID) (storage.Task, error) {
	_, err := s.getUser(userID)
	if err != nil {
		return storage.Task{}, err
	}

	t, err := s.getTask(taskID)
	if err != nil {
		return storage.Task{}, err
	}

	if t.UserID != userID {
		return storage.Task{}, ErrForbidden
	}

	return t, nil
}

func (s *Service) getAccessibleTask(userID, taskID uuid.UUID) (storage.Task, error) {
	_, err := s.getUser(userID)
	if err != nil {
		return storage.Task{}, err
	}

	t, err := s.getTask(taskID)
	if err != nil {
		return storage.Task{}, err
	}

	hasAccess, err := s.tasks.HasAccess(taskID, userID)
	if err != nil {
		return storage.Task{}, err
	}
	if !hasAccess {
		return storage.Task{}, ErrForbidden
	}

	return t, nil
}

func (s *Service) validateStatus(status storage.TaskStatus) error {
	switch status {
	case storage.StatusCreated, storage.StatusInProgress, storage.StatusDone:
		return nil
	default:
		return ErrInvalidInput
	}
}

func (s *Service) prepareTaskUpdate(userID uuid.UUID, groupID *uuid.UUID) error {
	_, err := s.getGroup(groupID)
	if err != nil {
		return err
	}

	if err := s.ensureUserInGroup(userID, groupID); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTask(id, taskID uuid.UUID) (storage.Task, error) {
	return s.getAccessibleTask(id, taskID)
}

func (s *Service) GetTaskForWorker(taskID uuid.UUID) (storage.Task, error) {
	return s.getTask(taskID)
}

func (s *Service) CreateTask(id uuid.UUID, name string, deadline *time.Time, groupID *uuid.UUID) (uuid.UUID, error) {
	if name == "" {
		return uuid.Nil, ErrInvalidInput
	}

	u, err := s.getUser(id)
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.prepareTaskUpdate(id, groupID); err != nil {
		return uuid.Nil, err
	}

	t := storage.Task{
		ID:       uuid.New(),
		Name:     name,
		Deadline: deadline,
		UserID:   u.ID,
		GroupID:  groupID,
		Status:   storage.StatusCreated,
	}

	if err := s.tasks.Create(t); err != nil {
		return uuid.Nil, err
	}
	time.Sleep(1 * time.Second)
	return t.ID, nil
}

func (s *Service) DeleteTask(id, taskID uuid.UUID) error {
	_, err := s.getOwnedTask(id, taskID)
	if err != nil {
		return err
	}

	if err := s.tasks.DeleteByID(taskID); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateTask(id, taskID uuid.UUID, name string, deadline *time.Time, groupID *uuid.UUID) error {
	if name == "" {
		return ErrInvalidInput
	}

	if err := s.prepareTaskUpdate(id, groupID); err != nil {
		return err
	}

	t, err := s.getAccessibleTask(id, taskID)
	if err != nil {
		return err
	}

	t.Name = name
	t.Deadline = deadline
	t.GroupID = groupID

	if err := s.tasks.Update(t); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateTaskStatus(taskID uuid.UUID, status storage.TaskStatus) error {
	if taskID == uuid.Nil {
		return ErrInvalidInput
	}

	if err := s.validateStatus(status); err != nil {
		return err
	}

	task, err := s.getTask(taskID)
	if err != nil {
		return err
	}

	if err := s.tasks.UpdateStatus(taskID, status); err != nil {
		return err
	}

	if s.notifier != nil {
		err := s.notifier.NotifyTaskStatusUpdated(events.TaskStatusUpdated{
			TaskID:  task.ID,
			GroupID: task.GroupID,
			UserID:  task.UserID,
			Status:  status,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetAllNotDone() ([]storage.Task, error) {
	return s.tasks.GetAllNotDone()
}

func (s *Service) InitMetrics() {
	count, err := s.tasks.Count()
	if err != nil {
		log.Printf("failed to count tasks: %v", err)
		return
	}

	metrics.TasksCurrent.Set(float64(count))
}
