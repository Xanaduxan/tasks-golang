package http_handlers

import (
	"net/http"

	"github.com/Xanaduxan/tasks-golang/internal/service/stocks"
)

var stockService *stocks.Service

func SetStockService(s *stocks.Service) { stockService = s }

func GetStocks(w http.ResponseWriter, r *http.Request) {
	st, err := stockService.GetStocks()
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, st)
}
