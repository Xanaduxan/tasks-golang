package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentStatus string

const (
	PaymentStatusNew                PaymentStatus = "NEW"
	PaymentStatusWaitingValidation2 PaymentStatus = "WAITING_FOR_VALIDATION_2"
	PaymentStatusFailed             PaymentStatus = "FAILED"
	PaymentStatusReadyForClosure    PaymentStatus = "READY_FOR_CLOSURE"
	PaymentStatusClosed             PaymentStatus = "CLOSED"
)

type Payment struct {
	ID                      uuid.UUID
	Status                  PaymentStatus
	ShopID                  uuid.UUID
	Amount                  decimal.Decimal
	WaitingForValidation2At sql.NullTime
	Attempts                int64
	UpdatedAt               time.Time
}

type PaymentStorage struct {
	DB *sql.DB
}

func NewPaymentStorage(db *sql.DB) *PaymentStorage {
	return &PaymentStorage{DB: db}
}

func (s *PaymentStorage) Create(payment Payment) error {
	_, err := s.DB.Exec(`
		INSERT INTO payments (id, status, shop_id, amount, attempts, waiting_for_validation_2_at)
		VALUES ($1, $2, $3, $4, $5, now())
	`, payment.ID, payment.Status, payment.ShopID, payment.Amount, payment.Attempts)

	return err
}

func (s *PaymentStorage) GetByID(id uuid.UUID) (Payment, error) {
	var payment Payment

	err := s.DB.QueryRow(`
		SELECT id, status, shop_id, amount, attempts, updated_at, waiting_for_validation_2_at
		FROM payments
		WHERE id = $1
	`, id).Scan(&payment.ID, &payment.Status, &payment.ShopID, &payment.Amount, &payment.Attempts, &payment.UpdatedAt, &payment.WaitingForValidation2At)

	return payment, err
}

func (s *PaymentStorage) Update(payment Payment) (int64, error) {
	current, err := s.GetByID(payment.ID)
	if err != nil {
		return 0, err
	}

	var waitingForValidation2At any

	switch {
	case current.Status != PaymentStatusWaitingValidation2 &&
		payment.Status == PaymentStatusWaitingValidation2:

		waitingForValidation2At = time.Now()

	case current.Status == PaymentStatusWaitingValidation2 &&
		payment.Status == PaymentStatusWaitingValidation2:

		if current.WaitingForValidation2At.Valid {
			waitingForValidation2At = current.WaitingForValidation2At.Time
		} else {
			waitingForValidation2At = time.Now()
		}

	default:

		waitingForValidation2At = nil
	}

	res, err := s.DB.Exec(`
		UPDATE payments
		SET status = $2,
		    attempts = $3,
		    amount=$4,
		    updated_at = now(),
		    waiting_for_validation_2_at = $5
		WHERE id = $1
	`, payment.ID, payment.Status, payment.Attempts, payment.Amount, waitingForValidation2At)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *PaymentStorage) Close(id uuid.UUID) (int64, error) {
	payment, err := s.GetByID(id)
	if err != nil {
		return 0, err
	}

	payment.Status = PaymentStatusClosed

	return s.Update(payment)
}

func (s *PaymentStorage) DeleteByID(id uuid.UUID) (int64, error) {
	res, err := s.DB.Exec(`
		DELETE FROM payments
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

func (s *PaymentStorage) GetShopPayments(shopId uuid.UUID) ([]Payment, error) {
	rows, err := s.DB.Query(`
		SELECT id, status, shop_id, amount, attempts, updated_at
		FROM payments
		WHERE shop_id = $1
	`, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.Status, &p.ShopID, &p.Amount, &p.Attempts, &p.UpdatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (s *PaymentStorage) GetPaymentsWaitingForValidation2() ([]Payment, error) {
	status := PaymentStatusWaitingValidation2
	rows, err := s.DB.Query(`
		SELECT id, status, shop_id, attempts, amount, updated_at, waiting_for_validation_2_at
		FROM payments
		WHERE status = $1
	`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.Status, &p.ShopID, &p.Attempts, &p.Amount, &p.UpdatedAt, &p.WaitingForValidation2At); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (s *PaymentStorage) GetPaymentsWaitingValidation2Expired() ([]Payment, error) {

	status := PaymentStatusWaitingValidation2

	rows, err := s.DB.Query(`
		SELECT id, status, shop_id, attempts, amount, updated_at, waiting_for_validation_2_at
		FROM payments
		WHERE status = $1
  AND waiting_for_validation_2_at <= now() - interval '1 hour'
	`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.Status, &p.ShopID, &p.Attempts, &p.Amount, &p.UpdatedAt, &p.WaitingForValidation2At); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
