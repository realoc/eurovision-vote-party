package handlers

import (
	"errors"
	"net/http"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

// ActsService defines the operations needed by the acts handler.
type ActsService interface {
	ListActs(eventType string) ([]models.Act, error)
}

// ActsHandler handles HTTP requests for acts.
type ActsHandler struct {
	service ActsService
}

// NewActsHandler creates a new ActsHandler.
func NewActsHandler(service ActsService) *ActsHandler {
	return &ActsHandler{service: service}
}

// actsResponse wraps the acts list for JSON serialization.
// Using a slice initialized to empty ensures JSON output is [] not null.
type actsResponse struct {
	Acts []models.Act `json:"acts"`
}

// ServeHTTP handles requests to /api/acts.
func (h *ActsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	eventType := r.URL.Query().Get("event")

	acts, err := h.service.ListActs(eventType)
	if err != nil {
		if errors.Is(err, services.ErrInvalidEventType) {
			writeError(w, http.StatusBadRequest)
			return
		}
		writeError(w, http.StatusInternalServerError)
		return
	}

	if acts == nil {
		acts = []models.Act{}
	}

	writeJSON(w, http.StatusOK, actsResponse{Acts: acts})
}
