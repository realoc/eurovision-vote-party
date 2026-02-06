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

// PartyService defines the operations needed by the handler.
type PartyService interface {
	CreateParty(ctx context.Context, adminID string, req services.CreatePartyRequest) (*models.Party, error)
	GetPartyByID(ctx context.Context, adminID, partyID string) (*models.Party, error)
	GetPartyByCode(ctx context.Context, code string) (*models.Party, error)
	ListPartiesByAdmin(ctx context.Context, adminID string) ([]*models.Party, error)
	DeleteParty(ctx context.Context, adminID, partyID string) error
}

// PartyHandler handles HTTP requests for party management.
type PartyHandler struct {
	service PartyService
}

// NewPartyHandler creates a new PartyHandler.
func NewPartyHandler(service PartyService) *PartyHandler {
	return &PartyHandler{service: service}
}

// createPartyRequest represents the request body for creating a party.
type createPartyRequest struct {
	Name      string           `json:"name"`
	EventType models.EventType `json:"eventType"`
}

// publicPartyResponse represents the public-facing party data.
type publicPartyResponse struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Code      string             `json:"code"`
	EventType models.EventType   `json:"eventType"`
	Status    models.PartyStatus `json:"status"`
}

// ServeHTTP routes requests to the appropriate handler method.
func (h *PartyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Path should be /api/parties or /api/parties/{id_or_code}
	path := strings.TrimPrefix(r.URL.Path, "/api/parties")
	path = strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodPost:
		if path == "" {
			h.handleCreate(w, r)
			return
		}
	case http.MethodGet:
		if path == "" {
			h.handleList(w, r)
			return
		}
		// Distinguish code (6 chars) from UUID (36 chars with dashes)
		if isPartyCode(path) {
			h.handleGetByCode(w, r, path)
			return
		}
		h.handleGetByID(w, r, path)
		return
	case http.MethodDelete:
		if path != "" {
			h.handleDelete(w, r, path)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// isPartyCode checks if the string looks like a party code (6 alphanumeric chars).
func isPartyCode(s string) bool {
	if len(s) != 6 {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= '2' && c <= '9')) {
			return false
		}
	}
	return true
}

// handleCreate handles POST /api/parties.
func (h *PartyHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	var req createPartyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest)
		return
	}

	party, err := h.service.CreateParty(r.Context(), userID, services.CreatePartyRequest{
		Name:      req.Name,
		EventType: req.EventType,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, party)
}

// handleList handles GET /api/parties.
func (h *PartyHandler) handleList(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	parties, err := h.service.ListPartiesByAdmin(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, parties)
}

// handleGetByCode handles GET /api/parties/:code (public endpoint).
func (h *PartyHandler) handleGetByCode(w http.ResponseWriter, r *http.Request, code string) {
	party, err := h.service.GetPartyByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeError(w, http.StatusNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	// Return only public information
	response := publicPartyResponse{
		ID:        party.ID,
		Name:      party.Name,
		Code:      party.Code,
		EventType: party.EventType,
		Status:    party.Status,
	}

	writeJSON(w, http.StatusOK, response)
}

// handleGetByID handles GET /api/parties/:id (authenticated endpoint).
func (h *PartyHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	party, err := h.service.GetPartyByID(r.Context(), userID, id)
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

	writeJSON(w, http.StatusOK, party)
}

// handleDelete handles DELETE /api/parties/:id.
func (h *PartyHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized)
		return
	}

	err := h.service.DeleteParty(r.Context(), userID, id)
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

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an HTTP error response.
func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
