package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type GroupMember struct {
	GroupID   uuid.UUID
	UserID    uuid.UUID
	Role      string
	CreatedAt time.Time
}

type GroupMemberStorage struct {
	DB *sql.DB
}

func NewGroupMemberStorage(db *sql.DB) *GroupMemberStorage {
	return &GroupMemberStorage{DB: db}
}

func (s *GroupMemberStorage) Create(member GroupMember) error {
	_, err := s.DB.Exec(`
		INSERT INTO group_members (group_id, user_id, role)
		VALUES ($1, $2, $3)
	`, member.GroupID, member.UserID, member.Role)

	return err
}

func (s *GroupMemberStorage) GetByID(groupID, userID uuid.UUID) (GroupMember, error) {
	var member GroupMember

	err := s.DB.QueryRow(`
		SELECT group_id, user_id, role, created_at
		FROM group_members
		WHERE group_id = $1 AND user_id = $2
	`, groupID, userID).Scan(
		&member.GroupID,
		&member.UserID,
		&member.Role,
		&member.CreatedAt,
	)

	return member, err
}

func (s *GroupMemberStorage) GetByGroupID(groupID uuid.UUID) ([]GroupMember, error) {
	rows, err := s.DB.Query(`
		SELECT group_id, user_id, role, created_at
		FROM group_members
		WHERE group_id = $1
		ORDER BY created_at
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []GroupMember
	for rows.Next() {
		var member GroupMember
		if err := rows.Scan(
			&member.GroupID,
			&member.UserID,
			&member.Role,
			&member.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, rows.Err()
}

func (s *GroupMemberStorage) Update(member GroupMember) error {
	_, err := s.DB.Exec(`
		UPDATE group_members
		SET role = $3
		WHERE group_id = $1 AND user_id = $2
	`, member.GroupID, member.UserID, member.Role)

	return err
}

func (s *GroupMemberStorage) DeleteByID(groupID, userID uuid.UUID) error {
	_, err := s.DB.Exec(`
		DELETE FROM group_members
		WHERE group_id = $1 AND user_id = $2
	`, groupID, userID)

	return err
}
func (s *GroupMemberStorage) IsMember(groupID, userID uuid.UUID) (bool, error) {
	var exists bool

	err := s.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM group_members
			WHERE group_id = $1 AND user_id = $2
		)
	`, groupID, userID).Scan(&exists)

	return exists, err
}
