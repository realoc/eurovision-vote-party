package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sipgate/eurovision-vote-party/server/middleware"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

// UserService defines the operations needed by the user handler.
type UserService interface {
	UpsertProfile(ctx context.Context, userID, email, username string) (*models.User, error)
	GetProfile(ctx context.Context, userID string) (*models.User, error)
}

// UserHandler handles HTTP requests for user profile management.
type UserHandler struct {
	service UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// updateProfileRequest represents the request body for updating a user profile.
type updateProfileRequest struct {
	Username string `json:"username"`
}

// ServeHTTP routes requests to the appropriate handler method.
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetProfile(w, r)
	case http.MethodPut:
		h.handleUpdateProfile(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed)
	}
}

// handleUpdateProfile handles PUT /api/users/profile.
func (h *UserHandler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	email, _ := middleware.UserEmailFromContext(r.Context())

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.UpsertProfile(r.Context(), userID, email, req.Username)
	if err != nil {
		if errors.Is(err, services.ErrInvalidUsername) {
			writeJSONError(w, http.StatusBadRequest, "Invalid username: must be 3-30 alphanumeric characters or underscores")
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// handleGetProfile handles GET /api/users/profile.
func (h *UserHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "Profile not found")
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
