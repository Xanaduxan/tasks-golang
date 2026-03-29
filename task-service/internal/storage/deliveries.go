package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type DeliveryStatus string

const (
	StatusAwaiting   DeliveryStatus = "awaiting"
	StatusOnPath     DeliveryStatus = "on_path"
	StatusProcessing DeliveryStatus = "processing"
	StatusChecked    DeliveryStatus = "checked"
	StatusAccepted   DeliveryStatus = "accepted"
)

type Delivery struct {
	ID        uuid.UUID
	Status    DeliveryStatus
	UserID    uuid.UUID
	UpdatedAt time.Time
}

type DeliveryStorage struct {
	DB *sql.DB
}

func NewDeliveryStorage(db *sql.DB) *DeliveryStorage {
	return &DeliveryStorage{DB: db}
}

func (s *DeliveryStorage) Create(delivery Delivery) error {
	_, err := s.DB.Exec(`
		INSERT INTO deliveries (id, status, user_id, updated_at)
		VALUES ($1, $2, $3, $4)
	`, delivery.ID, delivery.Status, delivery.UserID, delivery.UpdatedAt)

	return err
}

func (s *DeliveryStorage) GetByID(id uuid.UUID) (Delivery, error) {
	var delivery Delivery

	err := s.DB.QueryRow(`
		SELECT id, status, user_id, updated_at
		FROM deliveries
		WHERE id = $1
	`, id).Scan(&delivery.ID, &delivery.Status, &delivery.UserID, &delivery.UpdatedAt)

	return delivery, err
}

func (s *DeliveryStorage) Update(delivery Delivery) (int64, error) {
	res, err := s.DB.Exec(`
		UPDATE deliveries
		SET status = $2, updated_at = $3
		WHERE id = $1
	`, delivery.ID, delivery.Status, delivery.UpdatedAt)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *DeliveryStorage) DeleteByID(id uuid.UUID) (int64, error) {
	res, err := s.DB.Exec(`
		DELETE FROM deliveries
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
func (s *DeliveryStorage) GetByUserID(userID uuid.UUID) ([]Delivery, error) {
	rows, err := s.DB.Query(`
		SELECT id, status, user_id, updated_at
		FROM deliveries
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []Delivery

	for rows.Next() {
		var d Delivery
		if err := rows.Scan(&d.ID, &d.Status, &d.UserID, &d.UpdatedAt); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}

	return deliveries, rows.Err()
}

func (s *DeliveryStorage) GetAllNotAccepted() ([]Delivery, error) {
	rows, err := s.DB.Query(`
		SELECT id, status, user_id, updated_at
		FROM deliveries
		WHERE status <> $1
	`, StatusAccepted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []Delivery
	for rows.Next() {
		var d Delivery
		if err := rows.Scan(&d.ID, &d.Status, &d.UserID, &d.UpdatedAt); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deliveries, nil
}
