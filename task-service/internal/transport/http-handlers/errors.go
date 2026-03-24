package http_handlers

import (
	"errors"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/tasks"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, tasks.ErrInvalidInput):
		http.Error(w, "invalid input", http.StatusBadRequest)
	case errors.Is(err, tasks.ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, tasks.ErrForbidden):
		http.Error(w, "forbidden", http.StatusForbidden)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
