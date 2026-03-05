package deliveries

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	products      *storage.ProductStorage
	users         *storage.UserStorage
	deliveries    *storage.DeliveryStorage
	deliveryItems *storage.DeliveryItemStorage
	stocks        *storage.StockStorage
}

func NewService(
	products *storage.ProductStorage,
	users *storage.UserStorage,
	deliveries *storage.DeliveryStorage,
	deliveryItems *storage.DeliveryItemStorage,
	stocks *storage.StockStorage,
) *Service {
	return &Service{
		products:      products,
		users:         users,
		deliveries:    deliveries,
		deliveryItems: deliveryItems,
		stocks:        stocks,
	}
}

func (s *Service) CreateDelivery(userID uuid.UUID, items []storage.DeliveryItem) (uuid.UUID, error) {
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

	if err := s.deliveries.Create(storage.Delivery{
		ID:        deliveryID,
		Status:    storage.StatusAwaiting,
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

func (s *Service) GetDelivery(userID, deliveryID uuid.UUID) (storage.Delivery, []storage.DeliveryItem, error) {
	if userID == uuid.Nil || deliveryID == uuid.Nil {
		return storage.Delivery{}, nil, ErrInvalidInput
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Delivery{}, nil, ErrNotFound
		}
		return storage.Delivery{}, nil, err
	}

	d, err := s.deliveries.GetByID(deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Delivery{}, nil, ErrNotFound
		}
		return storage.Delivery{}, nil, err
	}
	if d.UserID != userID {
		return storage.Delivery{}, nil, ErrForbidden
	}

	items, err := s.deliveryItems.GetByDeliveryID(deliveryID)
	if err != nil {
		return storage.Delivery{}, nil, err
	}

	return d, items, nil
}

func (s *Service) UpdateDelivery(delivery storage.Delivery) error {
	if delivery.ID == uuid.Nil || delivery.UserID == uuid.Nil {
		return ErrInvalidInput
	}

	if _, err := s.users.GetByID(delivery.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	existing, err := s.deliveries.GetByID(delivery.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if existing.UserID != delivery.UserID {
		return ErrForbidden
	}

	rows, err := s.deliveries.Update(storage.Delivery{
		ID:        delivery.ID,
		Status:    delivery.Status,
		UserID:    delivery.UserID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	if existing.Status != storage.StatusAccepted && delivery.Status == storage.StatusAccepted {
		items, err := s.deliveryItems.GetByDeliveryID(delivery.ID)
		if err != nil {
			return err
		}
		for _, it := range items {
			if err := s.stocks.Increase(it.ProductID, it.Quantity); err != nil {
				return err
			}
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
