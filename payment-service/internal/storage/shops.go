package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type Shop struct {
	ID   uuid.UUID
	Name string
}

type ShopStorage struct {
	DB *sql.DB
}

func NewShopStorage(db *sql.DB) *ShopStorage {
	return &ShopStorage{DB: db}
}

func (s *ShopStorage) Create(shop Shop) error {
	_, err := s.DB.Exec(`
		INSERT INTO shops (id, name)
		VALUES ($1, $2)
	`, shop.ID, shop.Name)

	return err
}

func (s *ShopStorage) GetByID(id uuid.UUID) (Shop, error) {
	var shop Shop

	err := s.DB.QueryRow(`
		SELECT id, name
		FROM shops
		WHERE id = $1
	`, id).Scan(&shop.ID, &shop.Name)

	return shop, err
}

func (s *ShopStorage) Update(shop Shop) (int64, error) {
	res, err := s.DB.Exec(`
		UPDATE shops 
		SET name = $2
		WHERE id = $1
	`, shop.ID, shop.Name)

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *ShopStorage) DeleteByID(id uuid.UUID) (int64, error) {
	res, err := s.DB.Exec(`
		DELETE FROM shops
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

func (s *ShopStorage) GetShops() ([]Shop, error) {
	rows, err := s.DB.Query(`
		SELECT id, name
		FROM shops
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shops []Shop

	for rows.Next() {
		var p Shop
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		shops = append(shops, p)
	}

	return shops, rows.Err()
}
