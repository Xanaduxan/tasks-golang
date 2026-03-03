package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/service/tasks"
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Name     string     `json:"name"`
	Deadline *time.Time `json:"deadline"`
}

type CreateTaskResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateTaskRequest struct {
	Name     string     `json:"name"`
	Deadline *time.Time `json:"deadline"`
}

var taskService *tasks.Service

func SetTaskService(s *tasks.Service) { taskService = s }

func userIDFromContext(r *http.Request) (uuid.UUID, bool) {
	idStr, _ := r.Context().Value("id").(string)
	if idStr == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	taskID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := taskService.GetTask(userID, taskID)
	if err != nil {
		handleTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	newTaskID, err := taskService.CreateTask(userID, req.Name, req.Deadline)
	if err != nil {
		handleTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateTaskResponse{ID: newTaskID})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	taskID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := taskService.DeleteTask(userID, taskID); err != nil {
		handleTaskError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	taskID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := taskService.UpdateTask(userID, taskID, req.Name, req.Deadline); err != nil {
		handleTaskError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (req CreateTaskRequest) Validate() error { return validateTask(req.Name, req.Deadline) }
func (req UpdateTaskRequest) Validate() error { return validateTask(req.Name, req.Deadline) }
