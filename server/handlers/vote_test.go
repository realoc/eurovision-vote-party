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
	endVotingFunc  func(ctx context.Context, adminID, partyID string) (*models.Party, error)
	getResultsFunc func(ctx context.Context, adminID, partyID string) (*services.PartyResults, error)
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

func (m *mockVoteService) EndVoting(ctx context.Context, adminID, partyID string) (*models.Party, error) {
	if m.endVotingFunc != nil {
		return m.endVotingFunc(ctx, adminID, partyID)
	}
	return nil, nil
}

func (m *mockVoteService) GetResults(ctx context.Context, adminID, partyID string) (*services.PartyResults, error) {
	if m.getResultsFunc != nil {
		return m.getResultsFunc(ctx, adminID, partyID)
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

// --- EndVoting Tests ---

func TestVoteHandler_EndVoting_ReturnsOKOnSuccess(t *testing.T) {
	svc := &mockVoteService{
		endVotingFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			assert.Equal(t, "admin-1", adminID)
			assert.Equal(t, "party-1", partyID)
			return &models.Party{
				ID:     "party-1",
				Status: models.PartyStatusClosed,
			}, nil
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/end-voting", nil)
	req = requestWithUserID(req, "admin-1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "party-1", response["id"])
	assert.Equal(t, "closed", response["status"])
}

func TestVoteHandler_EndVoting_ReturnsUnauthorizedWhenNoAuth(t *testing.T) {
	handler := handlers.NewVoteHandler(&mockVoteService{})

	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/end-voting", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestVoteHandler_EndVoting_ReturnsForbiddenWhenUnauthorized(t *testing.T) {
	svc := &mockVoteService{
		endVotingFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			return nil, services.ErrUnauthorized
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/end-voting", nil)
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestVoteHandler_EndVoting_ReturnsForbiddenWhenPartyClosed(t *testing.T) {
	svc := &mockVoteService{
		endVotingFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			return nil, services.ErrPartyClosed
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/end-voting", nil)
	req = requestWithUserID(req, "admin-1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestVoteHandler_EndVoting_ReturnsNotFoundWhenPartyNotFound(t *testing.T) {
	svc := &mockVoteService{
		endVotingFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/parties/party-1/end-voting", nil)
	req = requestWithUserID(req, "admin-1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestVoteHandler_EndVoting_ReturnsMethodNotAllowedForWrongMethod(t *testing.T) {
	handler := handlers.NewVoteHandler(&mockVoteService{})

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/end-voting", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

// --- GetResults Tests ---

func TestVoteHandler_GetResults_ReturnsOKWithResults(t *testing.T) {
	svc := &mockVoteService{
		getResultsFunc: func(ctx context.Context, adminID, partyID string) (*services.PartyResults, error) {
			assert.Equal(t, "party-1", partyID)
			return &services.PartyResults{
				PartyID:     "party-1",
				PartyName:   "Test Party",
				TotalVoters: 2,
				Results: []models.VoteResult{
					{
						ActID:       "act-1",
						Country:     "Country 1",
						Artist:      "Artist 1",
						Song:        "Song 1",
						TotalPoints: 24,
						Rank:        1,
					},
				},
			}, nil
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/results", nil)
	req = requestWithUserID(req, "admin-1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response services.PartyResults
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "party-1", response.PartyID)
	assert.Equal(t, "Test Party", response.PartyName)
	assert.Equal(t, 2, response.TotalVoters)
	require.Len(t, response.Results, 1)
	assert.Equal(t, "act-1", response.Results[0].ActID)
	assert.Equal(t, "Country 1", response.Results[0].Country)
	assert.Equal(t, "Artist 1", response.Results[0].Artist)
	assert.Equal(t, "Song 1", response.Results[0].Song)
	assert.Equal(t, 24, response.Results[0].TotalPoints)
	assert.Equal(t, 1, response.Results[0].Rank)
}

func TestVoteHandler_GetResults_ReturnsForbiddenWhenVotingNotEnded(t *testing.T) {
	svc := &mockVoteService{
		getResultsFunc: func(ctx context.Context, adminID, partyID string) (*services.PartyResults, error) {
			return nil, services.ErrVotingNotEnded
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/results", nil)
	req = requestWithUserID(req, "admin-1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestVoteHandler_GetResults_ReturnsNotFoundWhenPartyNotFound(t *testing.T) {
	svc := &mockVoteService{
		getResultsFunc: func(ctx context.Context, adminID, partyID string) (*services.PartyResults, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/results", nil)
	req = requestWithUserID(req, "admin-1")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestVoteHandler_GetResults_ReturnsForbiddenWhenUnauthorized(t *testing.T) {
	svc := &mockVoteService{
		getResultsFunc: func(ctx context.Context, adminID, partyID string) (*services.PartyResults, error) {
			return nil, services.ErrUnauthorized
		},
	}

	handler := handlers.NewVoteHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/results", nil)
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
