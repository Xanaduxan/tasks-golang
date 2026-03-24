package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type DeliveryItem struct {
	ID         uuid.UUID
	DeliveryID uuid.UUID
	ProductID  uuid.UUID
	Quantity   int64
}

type DeliveryItemStorage struct {
	DB *sql.DB
}

func NewDeliveryItemStorage(db *sql.DB) *DeliveryItemStorage {
	return &DeliveryItemStorage{DB: db}
}

func (s *DeliveryItemStorage) Create(deliveryItem DeliveryItem) error {
	_, err := s.DB.Exec(`
		INSERT INTO delivery_items (id, delivery_id, product_id, quantity)
		VALUES ($1, $2, $3, $4)
	`, deliveryItem.ID, deliveryItem.DeliveryID, deliveryItem.ProductID, deliveryItem.Quantity)

	return err
}

func (s *DeliveryItemStorage) GetByID(id uuid.UUID) (DeliveryItem, error) {
	var deliveryItem DeliveryItem

	err := s.DB.QueryRow(`
		SELECT id, delivery_id, product_id, quantity
		FROM delivery_items
		WHERE id = $1
	`, id).Scan(&deliveryItem.ID, &deliveryItem.DeliveryID, &deliveryItem.ProductID, &deliveryItem.Quantity)

	return deliveryItem, err
}

func (s *DeliveryItemStorage) Update(deliveryItem DeliveryItem) (int64, error) {
	res, err := s.DB.Exec(`
		UPDATE delivery_items
		SET quantity=$2
		WHERE id = $1
	`, deliveryItem.ID, deliveryItem.Quantity)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *DeliveryItemStorage) DeleteByID(id uuid.UUID) (int64, error) {
	res, err := s.DB.Exec(`
		DELETE FROM delivery_items
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
func (s *DeliveryItemStorage) GetByDeliveryID(deliveryID uuid.UUID) ([]DeliveryItem, error) {
	rows, err := s.DB.Query(`
		SELECT id, delivery_id, product_id, quantity
		FROM delivery_items
		WHERE delivery_id = $1
		ORDER BY created_at DESC
	`, deliveryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveryItems []DeliveryItem

	for rows.Next() {
		var d DeliveryItem
		if err := rows.Scan(&d.ID, &d.DeliveryID, &d.ProductID, &d.Quantity); err != nil {
			return nil, err
		}
		deliveryItems = append(deliveryItems, d)
	}

	return deliveryItems, rows.Err()
}
