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

// VoteServiceHandler defines the operations needed by the vote handler.
type VoteServiceHandler interface {
	SubmitVote(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error)
	GetVotes(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error)
	UpdateVote(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error)
	EndVoting(ctx context.Context, adminID, partyID string) (*models.Party, error)
	GetResults(ctx context.Context, adminID, partyID string) (*services.PartyResults, error)
}

// VoteHandler handles HTTP requests for vote management.
type VoteHandler struct {
	service VoteServiceHandler
}

// NewVoteHandler creates a new VoteHandler.
func NewVoteHandler(service VoteServiceHandler) *VoteHandler {
	return &VoteHandler{service: service}
}

// submitVoteRequest represents the request body for submitting or updating a vote.
type submitVoteRequest struct {
	GuestID string         `json:"guestId"`
	Votes   map[int]string `json:"votes"`
}

// ServeHTTP routes requests to the appropriate handler method.
func (h *VoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/parties/")
	segments := strings.Split(path, "/")

	if len(segments) < 2 {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	partyID := segments[0]

	switch segments[1] {
	case "votes":
		switch len(segments) {
		case 2: // /api/parties/{partyID}/votes
			switch r.Method {
			case http.MethodPost:
				h.handleSubmitVote(w, r, partyID)
				return
			case http.MethodPut:
				h.handleUpdateVote(w, r, partyID)
				return
			}
		case 3: // /api/parties/{partyID}/votes/{guestID}
			if r.Method == http.MethodGet {
				h.handleGetVotes(w, r, partyID, segments[2])
				return
			}
		}
	case "end-voting":
		if r.Method == http.MethodPost {
			h.handleEndVoting(w, r, partyID)
			return
		}
	case "results":
		if r.Method == http.MethodGet {
			h.handleGetResults(w, r, partyID)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// handleSubmitVote handles POST /api/parties/:partyID/votes.
func (h *VoteHandler) handleSubmitVote(w http.ResponseWriter, r *http.Request, partyID string) {
	adminID, _ := middleware.UserIDFromContext(r.Context())

	var req submitVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.GuestID) == "" {
		writeError(w, http.StatusBadRequest)
		return
	}

	vote, err := h.service.SubmitVote(r.Context(), adminID, partyID, services.SubmitVoteRequest{
		GuestID: req.GuestID,
		Votes:   req.Votes,
	})
	if err != nil {
		mapVoteError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, vote)
}

// handleGetVotes handles GET /api/parties/:partyID/votes/:guestID.
func (h *VoteHandler) handleGetVotes(w http.ResponseWriter, r *http.Request, partyID, guestID string) {
	adminID, _ := middleware.UserIDFromContext(r.Context())

	vote, err := h.service.GetVotes(r.Context(), adminID, partyID, guestID)
	if err != nil {
		mapVoteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, vote)
}

// handleUpdateVote handles PUT /api/parties/:partyID/votes.
func (h *VoteHandler) handleUpdateVote(w http.ResponseWriter, r *http.Request, partyID string) {
	adminID, _ := middleware.UserIDFromContext(r.Context())

	var req submitVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.GuestID) == "" {
		writeError(w, http.StatusBadRequest)
		return
	}

	vote, err := h.service.UpdateVote(r.Context(), adminID, partyID, services.SubmitVoteRequest{
		GuestID: req.GuestID,
		Votes:   req.Votes,
	})
	if err != nil {
		mapVoteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, vote)
}

// handleEndVoting handles POST /api/parties/:partyID/end-voting.
func (h *VoteHandler) handleEndVoting(w http.ResponseWriter, r *http.Request, partyID string) {
	adminID, _ := middleware.UserIDFromContext(r.Context())
	if adminID == "" {
		writeError(w, http.StatusUnauthorized)
		return
	}

	party, err := h.service.EndVoting(r.Context(), adminID, partyID)
	if err != nil {
		mapVoteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":     party.ID,
		"status": party.Status,
	})
}

// handleGetResults handles GET /api/parties/:partyID/results.
func (h *VoteHandler) handleGetResults(w http.ResponseWriter, r *http.Request, partyID string) {
	adminID, _ := middleware.UserIDFromContext(r.Context())

	results, err := h.service.GetResults(r.Context(), adminID, partyID)
	if err != nil {
		mapVoteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, results)
}

// mapVoteError maps service errors to HTTP status codes.
func mapVoteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		writeError(w, http.StatusNotFound)
	case errors.Is(err, services.ErrUnauthorized):
		writeError(w, http.StatusForbidden)
	case errors.Is(err, services.ErrGuestNotApproved):
		writeError(w, http.StatusForbidden)
	case errors.Is(err, services.ErrPartyClosed):
		writeError(w, http.StatusForbidden)
	case errors.Is(err, services.ErrVoteAlreadyExists):
		writeError(w, http.StatusConflict)
	case errors.Is(err, services.ErrInvalidVotes):
		writeError(w, http.StatusBadRequest)
	case errors.Is(err, services.ErrVotingNotEnded):
		writeError(w, http.StatusForbidden)
	default:
		writeError(w, http.StatusInternalServerError)
	}
}
