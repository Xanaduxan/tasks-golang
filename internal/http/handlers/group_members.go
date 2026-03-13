package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/internal/service/group_members"
	"github.com/google/uuid"
)

type CreateGroupMemberRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type UpdateGroupMemberRequest struct {
	Role string `json:"role"`
}

var groupMemberService *group_members.GroupMemberService

func SetGroupMemberService(s *group_members.GroupMemberService) {
	groupMemberService = s
}

func GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	_, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := uuid.Parse(r.PathValue("group_id"))
	if err != nil {
		http.Error(w, "invalid group id", http.StatusBadRequest)
		return
	}

	members, err := groupMemberService.GetMembers(groupID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, members)
}

func CreateGroupMember(w http.ResponseWriter, r *http.Request) {
	_, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := uuid.Parse(r.PathValue("group_id"))
	if err != nil {
		http.Error(w, "invalid group id", http.StatusBadRequest)
		return
	}

	var req CreateGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := groupMemberService.CreateMember(groupID, req.UserID, req.Role); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateGroupMember(w http.ResponseWriter, r *http.Request) {
	_, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := uuid.Parse(r.PathValue("group_id"))
	if err != nil {
		http.Error(w, "invalid group id", http.StatusBadRequest)
		return
	}

	memberUserID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req UpdateGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := groupMemberService.UpdateMember(groupID, memberUserID, req.Role); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteGroupMember(w http.ResponseWriter, r *http.Request) {
	_, ok := userIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := uuid.Parse(r.PathValue("group_id"))
	if err != nil {
		http.Error(w, "invalid group id", http.StatusBadRequest)
		return
	}

	memberUserID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if err := groupMemberService.DeleteMember(groupID, memberUserID); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (req CreateGroupMemberRequest) Validate() error {
	if req.UserID == uuid.Nil {
		return group_members.ErrInvalidInput
	}
	return nil
}

func (req UpdateGroupMemberRequest) Validate() error {
	if req.Role == "" {
		return group_members.ErrInvalidInput
	}
	return nil
}
