package http_handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/internal/service/auth"
)

type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

var authService *auth.Service

func SetAuthService(s *auth.Service) { authService = s }

func Registration(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	token, err := authService.Registration(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		case errors.Is(err, auth.ErrEmailTaken):
			http.Error(w, "email already taken", http.StatusConflict)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, http.StatusCreated, RegistrationResponse{AccessToken: token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	token, err := authService.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		case errors.Is(err, auth.ErrInvalidCredentials):
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, http.StatusOK, LoginResponse{AccessToken: token})
}

func (req RegistrationRequest) Validate() error {
	if req.Email == "" || req.Password == "" {
		return auth.ErrInvalidInput
	}
	return nil
}

func (req LoginRequest) Validate() error {
	if req.Email == "" || req.Password == "" {
		return auth.ErrInvalidInput
	}
	return nil
}
