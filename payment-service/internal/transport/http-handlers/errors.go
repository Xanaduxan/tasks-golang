package http_handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/payment-service/internal/service/shops"
)

func handleError(w http.ResponseWriter, err error) {
	log.Printf("error: %v\n", err)
	switch {
	case errors.Is(err, shops.ErrInvalidInput):
		http.Error(w, "invalid input", http.StatusBadRequest)
	case errors.Is(err, shops.ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, shops.ErrForbidden):
		http.Error(w, "forbidden", http.StatusForbidden)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
