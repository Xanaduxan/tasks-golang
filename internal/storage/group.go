package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID
	Name      string
	CreatedBy uuid.UUID
	CreatedAt time.Time
}

type GroupStorage struct {
	DB *sql.DB
}

func NewGroupStorage(db *sql.DB) *GroupStorage {
	return &GroupStorage{DB: db}
}

func (s *GroupStorage) Create(group Group) error {
	_, err := s.DB.Exec(`
		INSERT INTO groups (id, name, created_by)
		VALUES ($1, $2, $3)
	`, group.ID, group.Name, group.CreatedBy)

	return err
}

func (s *GroupStorage) GetByID(id uuid.UUID) (Group, error) {
	var group Group

	err := s.DB.QueryRow(`
		SELECT id, name, created_by, created_at
		FROM groups
		WHERE id = $1
	`, id).Scan(
		&group.ID,
		&group.Name,
		&group.CreatedBy,
		&group.CreatedAt,
	)

	return group, err
}
