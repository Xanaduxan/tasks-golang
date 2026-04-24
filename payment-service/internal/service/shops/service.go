package shops

import (
	"database/sql"
	"errors"

	"github.com/Xanaduxan/tasks-golang/payment-service/internal/storage"
	"github.com/google/uuid"
)

type ShopInterface interface {
	Create(stock storage.Shop) error
	GetByID(id uuid.UUID) (storage.Shop, error)
	Update(stock storage.Shop) (int64, error)
	DeleteByID(id uuid.UUID) (int64, error)
	GetShops() ([]storage.Shop, error)
}
type Service struct {
	shops ShopInterface
}

func NewService(shops *storage.ShopStorage) *Service {
	return &Service{shops: shops}
}

func (s *Service) GetShops() ([]storage.Shop, error) {
	return s.shops.GetShops()
}
func (s *Service) GetShop(id uuid.UUID) (storage.Shop, error) {
	if id == uuid.Nil {
		return storage.Shop{}, ErrInvalidInput
	}

	p, err := s.shops.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Shop{}, ErrNotFound
		}
		return storage.Shop{}, err
	}

	return p, nil
}

func (s *Service) CreateShop(name string) (uuid.UUID, error) {
	if name == "" {
		return uuid.Nil, ErrInvalidInput
	}

	p := storage.Shop{
		ID:   uuid.New(),
		Name: name,
	}

	if err := s.shops.Create(p); err != nil {
		return uuid.Nil, err
	}

	return p.ID, nil
}

func (s *Service) DeleteShop(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidInput
	}

	rows, err := s.shops.DeleteByID(id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Service) UpdateShop(id uuid.UUID, name string) error {
	if id == uuid.Nil || name == "" {
		return ErrInvalidInput
	}

	p := storage.Shop{
		ID:   id,
		Name: name,
	}

	rows, err := s.shops.Update(p)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
