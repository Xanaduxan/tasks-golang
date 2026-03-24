package stocks

import (
	"github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	"github.com/google/uuid"
)

type StockInterface interface {
	Create(stock storage.Stock) error
	GetByID(productId uuid.UUID) (storage.Stock, error)
	GetStocks() ([]storage.Stock, error)
	Update(stock storage.Stock) (int64, error)
	DeleteByID(productId uuid.UUID) (int64, error)
	Increase(productID uuid.UUID, qty int64) error
}
type Service struct {
	stocks StockInterface
}

func NewService(stocks StockInterface) *Service {
	return &Service{stocks: stocks}
}

func (s *Service) GetStocks() ([]storage.Stock, error) {
	return s.stocks.GetStocks()
}
