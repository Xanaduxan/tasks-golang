package payments

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Xanaduxan/tasks-golang/payment-service/internal/storage"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentInterface interface {
	Create(payment storage.Payment) error
	GetByID(id uuid.UUID) (storage.Payment, error)
	Update(payment storage.Payment) (int64, error)
	Close(id uuid.UUID) (int64, error)
	DeleteByID(id uuid.UUID) (int64, error)
	GetShopPayments(shopId uuid.UUID) ([]storage.Payment, error)
	GetPaymentsWaitingForValidation2() ([]storage.Payment, error)
	GetPaymentsWaitingValidation2Expired() ([]storage.Payment, error)
}

type ShopInterface interface {
	GetByID(id uuid.UUID) (storage.Shop, error)
}

type PaymentsService struct {
	payments PaymentInterface
	shops    ShopInterface
}

func NewPaymentsService(
	payments PaymentInterface,
	shops ShopInterface,

) *PaymentsService {
	return &PaymentsService{
		payments: payments,
		shops:    shops,
	}
}

func (s *PaymentsService) GetPayment(id uuid.UUID) (storage.Payment, error) {
	if id == uuid.Nil {
		return storage.Payment{}, ErrInvalidInput
	}

	p, err := s.payments.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Payment{}, ErrNotFound
		}
		return storage.Payment{}, err
	}

	return p, nil
}

func (s *PaymentsService) CreatePayment(shopId string, amount decimal.Decimal) (uuid.UUID, error) {
	if shopId == "" {
		log.Printf("invalid shopId: %s", shopId)
		return uuid.Nil, ErrInvalidInput
	}

	if !validation1(amount) {
		log.Printf("invalid amount: %s", amount)
		return uuid.Nil, ErrInvalidInput
	}

	shopUUID, err := uuid.Parse(shopId)
	if err != nil {
		return uuid.Nil, ErrInvalidInput
	}

	_, err = s.shops.GetByID(shopUUID)
	if err != nil {
		return uuid.Nil, err
	}

	p := storage.Payment{
		ID:     uuid.New(),
		Status: storage.PaymentStatusWaitingValidation2,
		ShopID: shopUUID,
		Amount: amount,
	}

	if err := s.payments.Create(p); err != nil {
		return uuid.Nil, err
	}

	return p.ID, nil
}

func (s *PaymentsService) DeletePayment(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidInput
	}

	rows, err := s.payments.DeleteByID(id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PaymentsService) UpdatePayment(
	id uuid.UUID,
	amount *decimal.Decimal,
	status *storage.PaymentStatus,
	attempts *int64,
) error {
	if id == uuid.Nil {
		return ErrInvalidInput
	}

	p, err := s.GetPayment(id)
	if err != nil {
		log.Println(err)
		return err
	}

	if amount != nil {
		if amount.IsNegative() {
			return ErrInvalidInput
		}
		p.Amount = *amount
	}

	if attempts != nil {
		if *attempts < 0 {
			return ErrInvalidInput
		}
		p.Attempts = *attempts
	}

	if status != nil && p.Status != *status {
		if !isValidStatusTransition(p.Status, *status) {
			log.Printf("invalid transition: %s -> %s", p.Status, *status)
			return ErrInvalidInput
		}
		p.Status = *status
	}

	rows, err := s.payments.Update(p)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PaymentsService) ClosePayment(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidInput
	}

	p, err := s.GetPayment(id)
	if err != nil {
		log.Println(err)
		return err
	}

	if !isValidStatusTransition(p.Status, storage.PaymentStatusClosed) {
		log.Printf("invalid transition: %s -> %s", p.Status, storage.PaymentStatusClosed)
		return ErrInvalidInput
	}

	p.Status = storage.PaymentStatusClosed

	rows, err := s.payments.Update(p)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PaymentsService) GetShopPayments(shopId uuid.UUID) ([]storage.Payment, error) {
	if shopId == uuid.Nil {
		return nil, ErrInvalidInput
	}

	payments, err := s.payments.GetShopPayments(shopId)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func validation1(amount decimal.Decimal) bool {
	if amount.IsNegative() || amount.IsZero() {
		return false
	}

	maxVal := decimal.NewFromInt(10000)
	if amount.GreaterThan(maxVal) {
		return false
	}

	return true
}
