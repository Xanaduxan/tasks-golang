package http_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/payment-service/internal/service/payments"
	"github.com/Xanaduxan/tasks-golang/payment-service/internal/storage"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreatePaymentRequest struct {
	ShopID string          `json:"shop_id"`
	Amount decimal.Decimal `json:"amount"`
}

type CreatePaymentResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdatePaymentRequest struct {
	Amount   *decimal.Decimal       `json:"amount"`
	Status   *storage.PaymentStatus `json:"status"`
	Attempts *int64                 `json:"attempts"`
}

var PaymentsService *payments.PaymentsService

func SetPaymentsService(s *payments.PaymentsService) {
	PaymentsService = s
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	p, err := PaymentsService.GetPayment(id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	newID, err := PaymentsService.CreatePayment(req.ShopID, req.Amount)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreatePaymentResponse{ID: newID})
}

func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	var req UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := PaymentsService.UpdatePayment(id, req.Amount, req.Status, req.Attempts); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeletePayment(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	if err := PaymentsService.DeletePayment(id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ClosePayment(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	if err := PaymentsService.ClosePayment(id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetShopPayments(w http.ResponseWriter, r *http.Request) {
	shopID, ok := parseUUIDParam(w, r, "shop_id")
	if !ok {
		return
	}

	p, err := PaymentsService.GetShopPayments(shopID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (req CreatePaymentRequest) Validate() error {
	if req.ShopID == "" || req.Amount.IsNegative() {
		return payments.ErrInvalidInput
	}
	return nil
}

func (req UpdatePaymentRequest) Validate() error {
	if req.Amount != nil && req.Amount.IsNegative() {
		return payments.ErrInvalidInput
	}

	if req.Attempts != nil && *req.Attempts < 0 {
		return payments.ErrInvalidInput
	}

	if req.Status != nil {
		switch *req.Status {
		case storage.PaymentStatusNew,
			storage.PaymentStatusWaitingValidation2,
			storage.PaymentStatusFailed,
			storage.PaymentStatusReadyForClosure,
			storage.PaymentStatusClosed:
		default:
			return payments.ErrInvalidInput
		}
	}

	if req.Amount == nil && req.Status == nil && req.Attempts == nil {
		return payments.ErrInvalidInput
	}

	return nil
}
