package stocks

import (
	"github.com/Xanaduxan/tasks-golang/internal/storage"
)

type Service struct {
	stocks *storage.StockStorage
}

func NewService(stocks *storage.StockStorage) *Service {
	return &Service{stocks: stocks}
}

func (s *Service) GetStocks() ([]storage.Stock, error) {
	return s.stocks.GetStocks()
}
