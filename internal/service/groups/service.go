package groups

import (
	"database/sql"
	"errors"

	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/google/uuid"
)

type GroupInterface interface {
	Create(group storage.Group) error
	GetByID(id uuid.UUID) (storage.Group, error)
}
type GroupService struct {
	groups GroupInterface
}

func NewGroupService(groups GroupInterface) *GroupService {
	return &GroupService{
		groups: groups,
	}
}

func (s *GroupService) getGroup(groupID uuid.UUID) (storage.Group, error) {
	if groupID == uuid.Nil {
		return storage.Group{}, ErrInvalidInput
	}

	group, err := s.groups.GetByID(groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Group{}, ErrNotFound
		}
		return storage.Group{}, err
	}

	return group, nil
}

func (s *GroupService) CreateGroup(userID uuid.UUID, name string) (uuid.UUID, error) {
	if userID == uuid.Nil || name == "" {
		return uuid.Nil, ErrInvalidInput
	}

	group := storage.Group{
		ID:        uuid.New(),
		Name:      name,
		CreatedBy: userID,
	}

	if err := s.groups.Create(group); err != nil {
		return uuid.Nil, err
	}

	return group.ID, nil
}

func (s *GroupService) GetGroup(groupID uuid.UUID) (storage.Group, error) {
	return s.getGroup(groupID)
}
