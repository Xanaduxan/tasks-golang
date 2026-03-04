package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ProductID uuid.UUID
	Quantity  int64
	UpdatedAt time.Time
}

type StockStorage struct {
	DB *sql.DB
}

func NewStockStorage(db *sql.DB) *StockStorage {
	return &StockStorage{DB: db}
}

func (s *StockStorage) Create(stock Stock) error {
	_, err := s.DB.Exec(`
		INSERT INTO stocks (product_id, quantity, updated_at)
		VALUES ($1, $2, $3)
	`, stock.ProductID, stock.Quantity, stock.UpdatedAt)

	return err
}

func (s *StockStorage) GetByID(productId uuid.UUID) (Stock, error) {
	var stock Stock

	err := s.DB.QueryRow(`
		SELECT product_id, quantity, updated_at
		FROM stocks
		WHERE product_id = $1
	`, productId).Scan(&stock.ProductID, &stock.Quantity, &stock.UpdatedAt)

	return stock, err
}

func (s *StockStorage) Update(stock Stock) (int64, error) {
	res, err := s.DB.Exec(`
		UPDATE stocks
		SET quantity = $2, updated_at = $3
		WHERE product_id = $1
	`, stock.ProductID, stock.Quantity, stock.UpdatedAt)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *StockStorage) DeleteByID(productId uuid.UUID) (int64, error) {
	res, err := s.DB.Exec(`
		DELETE FROM stocks
		WHERE product_id = $1
	`, productId)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}
func (s *StockStorage) GetStocks() ([]Stock, error) {
	rows, err := s.DB.Query(`
		SELECT product_id, quantity, updated_at
		FROM stocks
		ORDER BY created_at DESC 
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []Stock

	for rows.Next() {
		var d Stock
		if err := rows.Scan(&d.ProductID, &d.Quantity, &d.UpdatedAt); err != nil {
			return nil, err
		}
		stocks = append(stocks, d)
	}

	return stocks, rows.Err()
}
