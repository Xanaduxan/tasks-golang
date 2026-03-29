package deliveries

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/auth"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/products"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/stocks"
	storage2 "github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	"github.com/google/uuid"
)

type DeliveryInterface interface {
	Create(delivery storage2.Delivery) error
	GetByID(id uuid.UUID) (storage2.Delivery, error)
	Update(delivery storage2.Delivery) (int64, error)
	DeleteByID(id uuid.UUID) (int64, error)
	GetByUserID(userID uuid.UUID) ([]storage2.Delivery, error)
	GetAllNotAccepted() ([]storage2.Delivery, error)
}

type DeliveryItemsInterface interface {
	Create(deliveryItem storage2.DeliveryItem) error
	GetByID(id uuid.UUID) (storage2.DeliveryItem, error)
	Update(deliveryItem storage2.DeliveryItem) (int64, error)
	DeleteByID(id uuid.UUID) (int64, error)
	GetByDeliveryID(deliveryID uuid.UUID) ([]storage2.DeliveryItem, error)
}
type Service struct {
	products      products.ProductInterface
	users         auth.UserInterface
	deliveries    DeliveryInterface
	deliveryItems DeliveryItemsInterface
	stocks        stocks.StockInterface
	notifier      Notifier
}

func NewService(
	products products.ProductInterface,
	users auth.UserInterface,
	deliveries DeliveryInterface,
	deliveryItems DeliveryItemsInterface,
	stocks stocks.StockInterface,
	notifier Notifier,
) *Service {
	return &Service{
		products:      products,
		users:         users,
		deliveries:    deliveries,
		deliveryItems: deliveryItems,
		stocks:        stocks,
		notifier:      notifier,
	}
}

func (s *Service) CreateDelivery(userID uuid.UUID, items []storage2.DeliveryItem) (uuid.UUID, error) {
	if userID == uuid.Nil || len(items) == 0 {
		return uuid.Nil, ErrInvalidInput
	}
	for _, it := range items {
		if it.ProductID == uuid.Nil || it.Quantity <= 0 {
			return uuid.Nil, ErrInvalidInput
		}
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrNotFound
		}
		return uuid.Nil, err
	}

	deliveryID := uuid.New()

	if err := s.deliveries.Create(storage2.Delivery{
		ID:        deliveryID,
		Status:    storage2.StatusAwaiting,
		UserID:    userID,
		UpdatedAt: time.Now(),
	}); err != nil {
		return uuid.Nil, err
	}

	for _, it := range items {
		it.ID = uuid.New()
		it.DeliveryID = deliveryID
		if err := s.deliveryItems.Create(it); err != nil {
			return uuid.Nil, err
		}
	}

	return deliveryID, nil
}

func (s *Service) GetDelivery(userID, deliveryID uuid.UUID) (storage2.Delivery, []storage2.DeliveryItem, error) {
	if userID == uuid.Nil || deliveryID == uuid.Nil {
		return storage2.Delivery{}, nil, ErrInvalidInput
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage2.Delivery{}, nil, ErrNotFound
		}
		return storage2.Delivery{}, nil, err
	}

	d, err := s.deliveries.GetByID(deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage2.Delivery{}, nil, ErrNotFound
		}
		return storage2.Delivery{}, nil, err
	}
	if d.UserID != userID {
		return storage2.Delivery{}, nil, ErrForbidden
	}

	items, err := s.deliveryItems.GetByDeliveryID(deliveryID)
	if err != nil {
		return storage2.Delivery{}, nil, err
	}

	return d, items, nil
}

func (s *Service) GetDeliveriesNotAccepted() ([]storage2.Delivery, error) {
	return s.deliveries.GetAllNotAccepted()
}

func (s *Service) UpdateDelivery(userID, deliveryID uuid.UUID, status storage2.DeliveryStatus) error {
	if userID == uuid.Nil || deliveryID == uuid.Nil {
		return ErrInvalidInput
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	existing, err := s.deliveries.GetByID(deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if existing.UserID != userID {
		return ErrForbidden
	}

	return s.UpdateDeliveryStatus(deliveryID, status)
}

func (s *Service) UpdateDeliveryStatus(deliveryID uuid.UUID, status storage2.DeliveryStatus) error {
	if deliveryID == uuid.Nil {
		return ErrInvalidInput
	}

	existing, err := s.deliveries.GetByID(deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if existing.Status == status {
		return nil
	}

	if !isValidStatusTransition(existing.Status, status) {
		return ErrInvalidInput
	}

	rows, err := s.deliveries.Update(storage2.Delivery{
		ID:        existing.ID,
		Status:    status,
		UserID:    existing.UserID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	if existing.Status != storage2.StatusAccepted && status == storage2.StatusAccepted {
		items, err := s.deliveryItems.GetByDeliveryID(deliveryID)
		if err != nil {
			return err
		}
		for _, it := range items {
			if err := s.stocks.Increase(it.ProductID, it.Quantity); err != nil {
				return err
			}
		}
	}

	if s.notifier != nil {
		err := s.notifier.SendNotification(
			[]string{existing.UserID.String()},
			"delivery.status_updated",
			existing.ID.String(),
			string(status),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DeleteDelivery(userID, deliveryID uuid.UUID) (int64, error) {
	if userID == uuid.Nil || deliveryID == uuid.Nil {
		return 0, ErrInvalidInput
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	d, err := s.deliveries.GetByID(deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	if d.UserID != userID {
		return 0, ErrForbidden
	}

	return s.deliveries.DeleteByID(deliveryID)
}

func (s *Service) GetDeliveryByID(deliveryID uuid.UUID) (storage2.Delivery, error) {
	if deliveryID == uuid.Nil {
		return storage2.Delivery{}, ErrInvalidInput
	}

	d, err := s.deliveries.GetByID(deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage2.Delivery{}, ErrNotFound
		}
		return storage2.Delivery{}, err
	}

	return d, nil
}

func isValidStatusTransition(from, to storage2.DeliveryStatus) bool {
	switch from {
	case storage2.StatusAwaiting:
		return to == storage2.StatusProcessing
	case storage2.StatusProcessing:
		return to == storage2.StatusChecked
	case storage2.StatusChecked:
		return to == storage2.StatusOnPath
	case storage2.StatusOnPath:
		return to == storage2.StatusAccepted
	default:
		return false
	}
}
