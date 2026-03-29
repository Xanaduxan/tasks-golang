package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	Create(user storage.User) error
	GetByEmail(email string) (storage.User, error)
	GetByID(id uuid.UUID) (storage.User, error)
	DeleteByID(id uuid.UUID) error
}
type Service struct {
	users     UserInterface
	jwtSecret []byte
}

func NewService(users UserInterface, jwtSecret []byte) *Service {
	return &Service{users: users, jwtSecret: jwtSecret}
}

func (s *Service) Registration(email, password string) (string, error) {

	if email == "" || password == "" {
		return "", ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := storage.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hash),
	}

	err = s.users.Create(user)
	if err != nil {
		return "", err
	}

	return s.generateToken(user.ID)
}

func (s *Service) Login(email, password string) (string, error) {

	if email == "" || password == "" {
		return "", ErrInvalidInput
	}

	user, err := s.users.GetByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateToken(user.ID)
}

func (s *Service) generateToken(id uuid.UUID) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
