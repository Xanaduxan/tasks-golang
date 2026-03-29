package http_handlers

import (
	"encoding/json"
	"net/http"

	groups2 "github.com/Xanaduxan/tasks-golang/task-service/internal/service/groups"
	"github.com/google/uuid"
)

type CreateGroupRequest struct {
	Name string `json:"name"`
}

type CreateGroupResponse struct {
	ID uuid.UUID `json:"id"`
}

var groupService *groups2.GroupService

func SetGroupService(s *groups2.GroupService) {
	groupService = s
}

func GetGroup(w http.ResponseWriter, r *http.Request) {
	_, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	group, err := groupService.GetGroup(groupID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, group)
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	groupID, err := groupService.CreateGroup(userID, req.Name)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateGroupResponse{ID: groupID})
}

func (req CreateGroupRequest) Validate() error {
	if req.Name == "" {
		return groups2.ErrInvalidInput
	}
	return nil
}
