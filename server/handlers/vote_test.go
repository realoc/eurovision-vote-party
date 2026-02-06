package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/handlers"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

type mockVoteService struct {
	submitVoteFunc func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error)
	getVotesFunc   func(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error)
	updateVoteFunc func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error)
}

func (m *mockVoteService) SubmitVote(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
	if m.submitVoteFunc != nil {
		return m.submitVoteFunc(ctx, adminID, partyID, req)
	}
	return nil, nil
}

func (m *mockVoteService) GetVotes(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error) {
	if m.getVotesFunc != nil {
		return m.getVotesFunc(ctx, adminID, partyID, guestID)
	}
	return nil, nil
}

func (m *mockVoteService) UpdateVote(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
	if m.updateVoteFunc != nil {
		return m.updateVoteFunc(ctx, adminID, partyID, req)
	}
	return nil, nil
}

func validVotes() map[int]string {
	return map[int]string{
		12: "act-1",
		10: "act-2",
		8:  "act-3",
		7:  "act-4",
		6:  "act-5",
		5:  "act-6",
		4:  "act-7",
		3:  "act-8",
		2:  "act-9",
		1:  "act-10",
	}
}

// --- SubmitVote Tests ---

func TestVoteHandler_SubmitVote_ReturnsCreatedWithValidRequest(t *testing.T) {
	createdVote := &models.Vote{
		ID:        "vote-1",
		GuestID:   "guest-1",
		PartyID:   "party-1",
		Votes:     validVotes(),
		CreatedAt: time.Now(),
	}

	svc := &mockVoteService{
		submitVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "party-1", partyID)
			assert.Equal(t, "guest-1", req.GuestID)
			assert.Equal(t, validVotes(), req.Votes)
			return createdVote, nil
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.Vote
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "vote-1", response.ID)
	assert.Equal(t, "guest-1", response.GuestID)
	assert.Equal(t, "party-1", response.PartyID)
}

func TestVoteHandler_SubmitVote_ReturnsCreatedWithoutAdminAuth(t *testing.T) {
	createdVote := &models.Vote{
		ID:        "vote-1",
		GuestID:   "guest-1",
		PartyID:   "party-1",
		Votes:     validVotes(),
		CreatedAt: time.Now(),
	}

	svc := &mockVoteService{
		submitVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			assert.Equal(t, "", adminID)
			assert.Equal(t, "party-1", partyID)
			assert.Equal(t, "guest-1", req.GuestID)
			return createdVote, nil
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response models.Vote
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "vote-1", response.ID)
}

func TestVoteHandler_SubmitVote_ReturnsBadRequestWithInvalidJSON(t *testing.T) {
	handler := handlers.NewVoteHandler(&mockVoteService{})

	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVoteHandler_SubmitVote_ReturnsBadRequestWithEmptyGuestID(t *testing.T) {
	handler := handlers.NewVoteHandler(&mockVoteService{})

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVoteHandler_SubmitVote_ReturnsNotFoundWhenPartyNotFound(t *testing.T) {
	svc := &mockVoteService{
		submitVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestVoteHandler_SubmitVote_ReturnsForbiddenWhenPartyClosed(t *testing.T) {
	svc := &mockVoteService{
		submitVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			return nil, services.ErrPartyClosed
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestVoteHandler_SubmitVote_ReturnsForbiddenWhenUnauthorized(t *testing.T) {
	svc := &mockVoteService{
		submitVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			return nil, services.ErrUnauthorized
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestVoteHandler_SubmitVote_ReturnsForbiddenWhenGuestNotApproved(t *testing.T) {
	svc := &mockVoteService{
		submitVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			return nil, services.ErrGuestNotApproved
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestVoteHandler_SubmitVote_ReturnsConflictWhenVoteAlreadyExists(t *testing.T) {
	svc := &mockVoteService{
		submitVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			return nil, services.ErrVoteAlreadyExists
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestVoteHandler_SubmitVote_ReturnsBadRequestWhenInvalidVotes(t *testing.T) {
	svc := &mockVoteService{
		submitVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			return nil, services.ErrInvalidVotes
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- GetVotes Tests ---

func TestVoteHandler_GetVotes_ReturnsVoteDataWithAdminAuth(t *testing.T) {
	vote := &models.Vote{
		ID:        "vote-1",
		GuestID:   "guest-1",
		PartyID:   "party-1",
		Votes:     validVotes(),
		CreatedAt: time.Now(),
	}

	svc := &mockVoteService{
		getVotesFunc: func(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error) {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "party-1", partyID)
			assert.Equal(t, "guest-1", guestID)
			return vote, nil
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/votes/guest-1", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.Vote
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "vote-1", response.ID)
	assert.Equal(t, "guest-1", response.GuestID)
	assert.Equal(t, "party-1", response.PartyID)
}

func TestVoteHandler_GetVotes_ReturnsVoteDataWithoutAdminAuth(t *testing.T) {
	vote := &models.Vote{
		ID:        "vote-1",
		GuestID:   "guest-1",
		PartyID:   "party-1",
		Votes:     validVotes(),
		CreatedAt: time.Now(),
	}

	svc := &mockVoteService{
		getVotesFunc: func(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error) {
			assert.Equal(t, "", adminID)
			assert.Equal(t, "party-1", partyID)
			assert.Equal(t, "guest-1", guestID)
			return vote, nil
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/votes/guest-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Vote
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "vote-1", response.ID)
}

func TestVoteHandler_GetVotes_ReturnsNotFoundWhenNotFound(t *testing.T) {
	svc := &mockVoteService{
		getVotesFunc: func(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/votes/guest-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestVoteHandler_GetVotes_ReturnsForbiddenWhenUnauthorized(t *testing.T) {
	svc := &mockVoteService{
		getVotesFunc: func(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error) {
			return nil, services.ErrUnauthorized
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/votes/guest-1", nil)
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

// --- UpdateVote Tests ---

func TestVoteHandler_UpdateVote_ReturnsOKWithValidRequest(t *testing.T) {
	updatedVote := &models.Vote{
		ID:        "vote-1",
		GuestID:   "guest-1",
		PartyID:   "party-1",
		Votes:     validVotes(),
		CreatedAt: time.Now(),
	}

	svc := &mockVoteService{
		updateVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "party-1", partyID)
			assert.Equal(t, "guest-1", req.GuestID)
			assert.Equal(t, validVotes(), req.Votes)
			return updatedVote, nil
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.Vote
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "vote-1", response.ID)
	assert.Equal(t, "guest-1", response.GuestID)
}

func TestVoteHandler_UpdateVote_ReturnsNotFoundWhenVoteNotFound(t *testing.T) {
	svc := &mockVoteService{
		updateVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestVoteHandler_UpdateVote_ReturnsBadRequestWhenInvalidVotes(t *testing.T) {
	svc := &mockVoteService{
		updateVoteFunc: func(ctx context.Context, adminID, partyID string, req services.SubmitVoteRequest) (*models.Vote, error) {
			return nil, services.ErrInvalidVotes
		},
	}

	handler := handlers.NewVoteHandler(svc)

	body, _ := json.Marshal(map[string]interface{}{
		"guestId": "guest-1",
		"votes":   validVotes(),
	})
	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/votes", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- Routing Tests ---

func TestVoteHandler_ReturnsMethodNotAllowedForUnsupportedMethod(t *testing.T) {
	handler := handlers.NewVoteHandler(&mockVoteService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/parties/party-1/votes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
