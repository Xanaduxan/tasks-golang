package http_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/queue"
	deliveries2 "github.com/Xanaduxan/tasks-golang/task-service/internal/service/deliveries"
	storage2 "github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	"github.com/google/uuid"
)

type CreateDeliveryItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
}

type CreateDeliveryRequest struct {
	Items []CreateDeliveryItemRequest `json:"items"`
}

type CreateDeliveryResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateDeliveryRequest struct {
	Status storage2.DeliveryStatus `json:"status"`
}

type DeliveryResponse struct {
	Delivery storage2.Delivery       `json:"delivery"`
	Items    []storage2.DeliveryItem `json:"items"`
}

var deliveryService *deliveries2.Service
var deliveryQueue *queue.DeliveryQueue

func SetDeliveryService(s *deliveries2.Service) {
	deliveryService = s
}

func SetDeliveryQueue(q *queue.DeliveryQueue) {
	deliveryQueue = q
}

func GetDelivery(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	deliveryID, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	d, items, err := deliveryService.GetDelivery(userID, deliveryID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, DeliveryResponse{
		Delivery: d,
		Items:    items,
	})
}

func CreateDelivery(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	items := make([]storage2.DeliveryItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, storage2.DeliveryItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
		})
	}

	newID, err := deliveryService.CreateDelivery(userID, items)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateDeliveryResponse{ID: newID})
}

func DeleteDelivery(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	deliveryID, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	_, err := deliveryService.DeleteDelivery(userID, deliveryID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	deliveryID, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	var req UpdateDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := deliveryService.UpdateDelivery(userID, deliveryID, req.Status); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (req CreateDeliveryRequest) Validate() error {
	if len(req.Items) == 0 {
		return deliveries2.ErrInvalidInput
	}
	for _, it := range req.Items {
		if it.ProductID == uuid.Nil || it.Quantity <= 0 {
			return deliveries2.ErrInvalidInput
		}
	}
	return nil
}

func (req UpdateDeliveryRequest) Validate() error {
	switch req.Status {
	case storage2.StatusAwaiting, storage2.StatusOnPath, storage2.StatusProcessing, storage2.StatusChecked, storage2.StatusAccepted:
		return nil
	default:
		return deliveries2.ErrInvalidInput
	}
}
