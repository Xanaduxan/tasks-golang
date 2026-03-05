package products

import (
	"database/sql"
	"errors"

	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service struct {
	products *storage.ProductStorage
}

func NewService(products *storage.ProductStorage) *Service {
	return &Service{products: products}
}

func (s *Service) GetProduct(id uuid.UUID) (storage.Product, error) {
	if id == uuid.Nil {
		return storage.Product{}, ErrInvalidInput
	}

	p, err := s.products.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Product{}, ErrNotFound
		}
		return storage.Product{}, err
	}

	return p, nil
}

func (s *Service) CreateProduct(name string, price decimal.Decimal) (uuid.UUID, error) {
	if name == "" || price.LessThanOrEqual(decimal.Zero) {
		return uuid.Nil, ErrInvalidInput
	}

	p := storage.Product{
		ID:    uuid.New(),
		Name:  name,
		Price: price,
	}

	if err := s.products.Create(p); err != nil {
		return uuid.Nil, err
	}

	return p.ID, nil
}

func (s *Service) DeleteProduct(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidInput
	}

	rows, err := s.products.DeleteByID(id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Service) UpdateProduct(id uuid.UUID, name string, price decimal.Decimal) error {
	if id == uuid.Nil || name == "" || price.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidInput
	}

	p := storage.Product{
		ID:    id,
		Name:  name,
		Price: price,
	}

	rows, err := s.products.Update(p)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
