package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
}

type UserStorage struct {
	DB *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{DB: db}
}

func (s *UserStorage) Create(user User) error {
	_, err := s.DB.Exec(`
		INSERT INTO users (id, email, password)
		VALUES ($1, $2, $3)
	`, user.ID, user.Email, user.Password)

	return err
}

func (s *UserStorage) GetByEmail(email string) (User, error) {
	var user User

	err := s.DB.QueryRow(`
		SELECT id, email, password
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.Password)

	return user, err
}

func (s *UserStorage) GetById(id uuid.UUID) (User, error) {
	var user User

	err := s.DB.QueryRow(`
		SELECT id, email, password
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Email, &user.Password)

	return user, err
}

func (s *UserStorage) DeleteByID(id uuid.UUID) error {

	_, err := s.DB.Exec(`
		DELETE FROM users
		WHERE id = $1
	`, id)

	return err
}
