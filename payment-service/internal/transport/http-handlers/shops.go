package http_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/payment-service/internal/service/shops"
	"github.com/google/uuid"
)

type CreateShopRequest struct {
	Name string `json:"name"`
}

type CreateShopResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateShopRequest struct {
	Name string `json:"name"`
}

var shopService *shops.Service

func SetShopService(s *shops.Service) { shopService = s }

func GetShops(w http.ResponseWriter, r *http.Request) {
	st, err := shopService.GetShops()
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, st)
}

func GetShop(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	p, err := shopService.GetShop(id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func CreateShop(w http.ResponseWriter, r *http.Request) {
	var req CreateShopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	newID, err := shopService.CreateShop(req.Name)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateShopResponse{ID: newID})
}

func DeleteShop(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	if err := shopService.DeleteShop(id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateShop(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	var req UpdateShopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := shopService.UpdateShop(id, req.Name); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (req CreateShopRequest) Validate() error {
	if req.Name == "" {
		return shops.ErrInvalidInput
	}
	return nil
}

func (req UpdateShopRequest) Validate() error {
	if req.Name == "" {
		return shops.ErrInvalidInput
	}
	return nil
}

func parseUUIDParam(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	idStr := r.PathValue(name)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid "+name, http.StatusBadRequest)
		return uuid.Nil, false
	}

	return id, true
}
