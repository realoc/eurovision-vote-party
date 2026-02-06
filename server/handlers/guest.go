package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/sipgate/eurovision-vote-party/server/middleware"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

// GuestService defines the operations needed by the guest handler.
type GuestService interface {
	JoinParty(ctx context.Context, code, username string) (*models.Guest, error)
	ListGuests(ctx context.Context, adminID, partyID string) ([]*models.Guest, error)
	ListGuestsAsGuest(ctx context.Context, guestID, partyID string) ([]*models.Guest, error)
	ListJoinRequests(ctx context.Context, adminID, partyID string) ([]*models.Guest, error)
	ApproveGuest(ctx context.Context, adminID, partyID, guestID string) error
	RejectGuest(ctx context.Context, adminID, partyID, guestID string) error
	RemoveGuest(ctx context.Context, adminID, partyID, guestID string) error
	GetGuestStatus(ctx context.Context, code, guestID string) (*models.Guest, error)
}

// GuestHandler handles HTTP requests for guest management.
type GuestHandler struct {
	service GuestService
}

// NewGuestHandler creates a new GuestHandler.
func NewGuestHandler(service GuestService) *GuestHandler {
	return &GuestHandler{service: service}
}

// joinPartyRequest represents the request body for joining a party.
type joinPartyRequest struct {
	Username string `json:"username"`
}

// ServeHTTP routes requests to the appropriate handler method.
func (h *GuestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/parties/")
	segments := strings.Split(path, "/")

	switch len(segments) {
	case 2:
		switch segments[1] {
		case "join":
			if r.Method == http.MethodPost {
				h.handleJoinParty(w, r, segments[0])
				return
			}
		case "guests":
			if r.Method == http.MethodGet {
				h.handleListGuests(w, r, segments[0])
				return
			}
		case "join-requests":
			if r.Method == http.MethodGet {
				h.handleListJoinRequests(w, r, segments[0])
				return
			}
		case "guest-status":
			if r.Method == http.MethodGet {
				h.handleGetGuestStatus(w, r, segments[0])
				return
			}
		}
	case 3:
		if segments[1] == "guests" && r.Method == http.MethodDelete {
			h.handleRemoveGuest(w, r, segments[0], segments[2])
			return
		}
	case 4:
		if segments[1] == "guests" {
			switch segments[3] {
			case "approve":
				if r.Method == http.MethodPut {
					h.handleApproveGuest(w, r, segments[0], segments[2])
					return
				}
			case "reject":
				if r.Method == http.MethodPut {
					h.handleRejectGuest(w, r, segments[0], segments[2])
					return
				}
			}
		}
	}

	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// handleJoinParty handles POST /api/parties/:code/join.
func (h *GuestHandler) handleJoinParty(w http.ResponseWriter, r *http.Request, code string) {
	var req joinPartyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Username) == "" {
		writeError(w, http.StatusBadRequest)
		return
	}

	guest, err := h.service.JoinParty(r.Context(), code, req.Username)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeError(w, http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrDuplicateUsername) {
			writeError(w, http.StatusConflict)
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, guest)
}

// handleListGuests handles GET /api/parties/:id/guests.
func (h *GuestHandler) handleListGuests(w http.ResponseWriter, r *http.Request, partyID string) {
	// Try admin auth first (OptionalAuthMiddleware may have set it).
	userID, ok := middleware.UserIDFromContext(r.Context())
	if ok {
		guests, err := h.service.ListGuests(r.Context(), userID, partyID)
		if err != nil {
			if errors.Is(err, services.ErrUnauthorized) {
				writeError(w, http.StatusForbidden)
				return
			}
			if errors.Is(err, services.ErrNotFound) {
				writeError(w, http.StatusNotFound)
				return
			}
			writeError(w, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, guests)
		return
	}

	// Fall back to guest auth via query parameter.
	guestID := r.URL.Query().Get("guestId")
	if guestID != "" {
		guests, err := h.service.ListGuestsAsGuest(r.Context(), guestID, partyID)
		if err != nil {
			if errors.Is(err, services.ErrUnauthorized) {
				writeError(w, http.StatusForbidden)
				return
			}
			if errors.Is(err, services.ErrNotFound) {
				writeError(w, http.StatusNotFound)
				return
			}
			writeError(w, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, guests)
		return
	}

	writeError(w, http.StatusUnauthorized)
}

// handleListJoinRequests handles GET /api/parties/:id/join-requests.
func (h *GuestHandler) handleListJoinRequests(w http.ResponseWriter, r *http.Request, partyID string) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	guests, err := h.service.ListJoinRequests(r.Context(), userID, partyID)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			writeError(w, http.StatusForbidden)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			writeError(w, http.StatusNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, guests)
}

// handleApproveGuest handles PUT /api/parties/:id/guests/:guestId/approve.
func (h *GuestHandler) handleApproveGuest(w http.ResponseWriter, r *http.Request, partyID, guestID string) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	err := h.service.ApproveGuest(r.Context(), userID, partyID, guestID)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			writeError(w, http.StatusForbidden)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			writeError(w, http.StatusNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleRejectGuest handles PUT /api/parties/:id/guests/:guestId/reject.
func (h *GuestHandler) handleRejectGuest(w http.ResponseWriter, r *http.Request, partyID, guestID string) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	err := h.service.RejectGuest(r.Context(), userID, partyID, guestID)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			writeError(w, http.StatusForbidden)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			writeError(w, http.StatusNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleRemoveGuest handles DELETE /api/parties/:id/guests/:guestId.
func (h *GuestHandler) handleRemoveGuest(w http.ResponseWriter, r *http.Request, partyID, guestID string) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	err := h.service.RemoveGuest(r.Context(), userID, partyID, guestID)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			writeError(w, http.StatusForbidden)
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			writeError(w, http.StatusNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleGetGuestStatus handles GET /api/parties/:code/guest-status.
func (h *GuestHandler) handleGetGuestStatus(w http.ResponseWriter, r *http.Request, code string) {
	guestID := r.URL.Query().Get("guestId")
	if guestID == "" {
		writeError(w, http.StatusBadRequest)
		return
	}

	guest, err := h.service.GetGuestStatus(r.Context(), code, guestID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeError(w, http.StatusNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, guest)
}
