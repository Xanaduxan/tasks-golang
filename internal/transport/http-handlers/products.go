package http_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/internal/service/products"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductRequest struct {
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price"`
}

type CreateProductResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateProductRequest struct {
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price"`
}

var productService *products.Service

func SetProductService(s *products.Service) { productService = s }

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	p, err := productService.GetProduct(id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	newID, err := productService.CreateProduct(req.Name, req.Price)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateProductResponse{ID: newID})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	if err := productService.DeleteProduct(id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := productService.UpdateProduct(id, req.Name, req.Price); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (req CreateProductRequest) Validate() error {
	if req.Name == "" || req.Price.LessThanOrEqual(decimal.Zero) {
		return products.ErrInvalidInput
	}
	return nil
}

func (req UpdateProductRequest) Validate() error {
	if req.Name == "" || req.Price.LessThanOrEqual(decimal.Zero) {
		return products.ErrInvalidInput
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
