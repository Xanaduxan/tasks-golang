package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type Product struct {
	ID    uuid.UUID
	Name  string
	Price int64
}

type ProductStorage struct {
	DB *sql.DB
}

func NewProductStorage(db *sql.DB) *ProductStorage {
	return &ProductStorage{DB: db}
}

func (s *ProductStorage) Create(product Product) error {
	_, err := s.DB.Exec(`
		INSERT INTO products (id, name, price)
		VALUES ($1, $2, $3)
	`, product.ID, product.Name, product.Price)

	return err
}

func (s *ProductStorage) GetByID(id uuid.UUID) (Product, error) {
	var product Product

	err := s.DB.QueryRow(`
		SELECT id, name, price
		FROM products
		WHERE id = $1
	`, id).Scan(&product.ID, &product.Name, &product.Price)

	return product, err
}

func (s *ProductStorage) Update(product Product) (int64, error) {
	res, err := s.DB.Exec(`
		UPDATE products 
		SET name = $2, price = $3
		WHERE id = $1
	`, product.ID, product.Name, product.Price)

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *ProductStorage) DeleteByID(id uuid.UUID) (int64, error) {
	res, err := s.DB.Exec(`
		DELETE FROM products
		WHERE id = $1
	`, id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *ProductStorage) GetProducts() ([]Product, error) {
	rows, err := s.DB.Query(`
		SELECT id, name, price
		FROM products
		ORDER BY created_at
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, rows.Err()
}
