package group_members

import (
	"database/sql"
	"errors"

	"github.com/Xanaduxan/tasks-golang/internal/service/auth"
	"github.com/Xanaduxan/tasks-golang/internal/service/groups"
	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/google/uuid"
)

type GroupMemberInterface interface {
	Create(member storage.GroupMember) error
	GetByID(groupID, userID uuid.UUID) (storage.GroupMember, error)
	GetByGroupID(groupID uuid.UUID) ([]storage.GroupMember, error)
	Update(member storage.GroupMember) error
	DeleteByID(groupID, userID uuid.UUID) error
	IsMember(groupID, userID uuid.UUID) (bool, error)
}
type GroupMemberService struct {
	members GroupMemberInterface
	groups  groups.GroupInterface
	users   auth.UserInterface
}

func NewGroupMemberService(
	members GroupMemberInterface,
	groups groups.GroupInterface,
	users auth.UserInterface,
) *GroupMemberService {
	return &GroupMemberService{
		members: members,
		groups:  groups,
		users:   users,
	}
}

func (s *GroupMemberService) getGroup(groupID uuid.UUID) (storage.Group, error) {
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

func (s *GroupMemberService) getUser(userID uuid.UUID) (storage.User, error) {
	if userID == uuid.Nil {
		return storage.User{}, ErrInvalidInput
	}

	user, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.User{}, ErrNotFound
		}
		return storage.User{}, err
	}

	return user, nil
}

func (s *GroupMemberService) getMember(groupID, userID uuid.UUID) (storage.GroupMember, error) {
	if groupID == uuid.Nil || userID == uuid.Nil {
		return storage.GroupMember{}, ErrInvalidInput
	}

	member, err := s.members.GetByID(groupID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.GroupMember{}, ErrNotFound
		}
		return storage.GroupMember{}, err
	}

	return member, nil
}

func (s *GroupMemberService) CreateMember(groupID, userID uuid.UUID, role string) error {
	if role == "" {
		role = "member"
	}

	_, err := s.getGroup(groupID)
	if err != nil {
		return err
	}

	_, err = s.getUser(userID)
	if err != nil {
		return err
	}

	member := storage.GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    role,
	}

	if err := s.members.Create(member); err != nil {
		return err
	}

	return nil
}

func (s *GroupMemberService) GetMembers(groupID uuid.UUID) ([]storage.GroupMember, error) {
	if groupID == uuid.Nil {
		return nil, ErrInvalidInput
	}

	_, err := s.getGroup(groupID)
	if err != nil {
		return nil, err
	}

	return s.members.GetByGroupID(groupID)
}

func (s *GroupMemberService) UpdateMember(groupID, userID uuid.UUID, role string) error {
	if role == "" {
		return ErrInvalidInput
	}

	member, err := s.getMember(groupID, userID)
	if err != nil {
		return err
	}

	member.Role = role

	if err := s.members.Update(member); err != nil {
		return err
	}

	return nil
}

func (s *GroupMemberService) DeleteMember(groupID, userID uuid.UUID) error {
	_, err := s.getMember(groupID, userID)
	if err != nil {
		return err
	}

	if err := s.members.DeleteByID(groupID, userID); err != nil {
		return err
	}

	return nil
}
func (s *GroupMemberService) IsMember(groupID, userID uuid.UUID) (bool, error) {
	return s.members.IsMember(groupID, userID)
}
