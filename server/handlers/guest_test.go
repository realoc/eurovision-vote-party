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

type mockGuestService struct {
	joinPartyFunc        func(ctx context.Context, code, username string) (*models.Guest, error)
	listGuestsFunc       func(ctx context.Context, adminID, partyID string) ([]*models.Guest, error)
	listGuestsAsGuestFunc func(ctx context.Context, guestID, partyID string) ([]*models.Guest, error)
	listJoinRequestsFunc func(ctx context.Context, adminID, partyID string) ([]*models.Guest, error)
	approveGuestFunc     func(ctx context.Context, adminID, partyID, guestID string) error
	rejectGuestFunc      func(ctx context.Context, adminID, partyID, guestID string) error
	removeGuestFunc      func(ctx context.Context, adminID, partyID, guestID string) error
	getGuestStatusFunc   func(ctx context.Context, code, guestID string) (*models.Guest, error)
}

func (m *mockGuestService) JoinParty(ctx context.Context, code, username string) (*models.Guest, error) {
	if m.joinPartyFunc != nil {
		return m.joinPartyFunc(ctx, code, username)
	}
	return nil, nil
}

func (m *mockGuestService) ListGuests(ctx context.Context, adminID, partyID string) ([]*models.Guest, error) {
	if m.listGuestsFunc != nil {
		return m.listGuestsFunc(ctx, adminID, partyID)
	}
	return nil, nil
}

func (m *mockGuestService) ListGuestsAsGuest(ctx context.Context, guestID, partyID string) ([]*models.Guest, error) {
	if m.listGuestsAsGuestFunc != nil {
		return m.listGuestsAsGuestFunc(ctx, guestID, partyID)
	}
	return nil, nil
}

func (m *mockGuestService) ListJoinRequests(ctx context.Context, adminID, partyID string) ([]*models.Guest, error) {
	if m.listJoinRequestsFunc != nil {
		return m.listJoinRequestsFunc(ctx, adminID, partyID)
	}
	return nil, nil
}

func (m *mockGuestService) ApproveGuest(ctx context.Context, adminID, partyID, guestID string) error {
	if m.approveGuestFunc != nil {
		return m.approveGuestFunc(ctx, adminID, partyID, guestID)
	}
	return nil
}

func (m *mockGuestService) RejectGuest(ctx context.Context, adminID, partyID, guestID string) error {
	if m.rejectGuestFunc != nil {
		return m.rejectGuestFunc(ctx, adminID, partyID, guestID)
	}
	return nil
}

func (m *mockGuestService) RemoveGuest(ctx context.Context, adminID, partyID, guestID string) error {
	if m.removeGuestFunc != nil {
		return m.removeGuestFunc(ctx, adminID, partyID, guestID)
	}
	return nil
}

func (m *mockGuestService) GetGuestStatus(ctx context.Context, code, guestID string) (*models.Guest, error) {
	if m.getGuestStatusFunc != nil {
		return m.getGuestStatusFunc(ctx, code, guestID)
	}
	return nil, services.ErrNotFound
}

// --- Join Party Tests ---

func TestGuestHandler_JoinParty_ReturnsCreatedWithValidRequest(t *testing.T) {
	createdGuest := &models.Guest{
		ID:        "guest-1",
		PartyID:   "party-1",
		Username:  "alice",
		Status:    models.GuestStatusPending,
		CreatedAt: time.Now(),
	}

	svc := &mockGuestService{
		joinPartyFunc: func(ctx context.Context, code, username string) (*models.Guest, error) {
			assert.Equal(t, "ABC234", code)
			assert.Equal(t, "alice", username)
			return createdGuest, nil
		},
	}

	handler := handlers.NewGuestHandler(svc)

	body := `{"username": "alice"}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties/ABC234/join", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.Guest
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "guest-1", response.ID)
	assert.Equal(t, "alice", response.Username)
	assert.Equal(t, models.GuestStatusPending, response.Status)
}

func TestGuestHandler_JoinParty_ReturnsBadRequestWithMissingUsername(t *testing.T) {
	handler := handlers.NewGuestHandler(&mockGuestService{})

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties/ABC234/join", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGuestHandler_JoinParty_ReturnsBadRequestWithBlankUsername(t *testing.T) {
	handler := handlers.NewGuestHandler(&mockGuestService{})

	body := `{"username": "   "}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties/ABC234/join", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGuestHandler_JoinParty_ReturnsNotFoundWhenPartyNotFound(t *testing.T) {
	svc := &mockGuestService{
		joinPartyFunc: func(ctx context.Context, code, username string) (*models.Guest, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewGuestHandler(svc)

	body := `{"username": "alice"}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties/ABC234/join", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGuestHandler_JoinParty_ReturnsConflictOnDuplicateUsername(t *testing.T) {
	svc := &mockGuestService{
		joinPartyFunc: func(ctx context.Context, code, username string) (*models.Guest, error) {
			return nil, services.ErrDuplicateUsername
		},
	}

	handler := handlers.NewGuestHandler(svc)

	body := `{"username": "alice"}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties/ABC234/join", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

// --- List Guests Tests ---

func TestGuestHandler_ListGuests_ReturnsGuestsForAdmin(t *testing.T) {
	guests := []*models.Guest{
		{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		},
		{
			ID:        "guest-2",
			PartyID:   "party-1",
			Username:  "bob",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		},
	}

	svc := &mockGuestService{
		listGuestsFunc: func(ctx context.Context, adminID, partyID string) ([]*models.Guest, error) {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "party-1", partyID)
			return guests, nil
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/guests", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response []*models.Guest
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "guest-1", response[0].ID)
	assert.Equal(t, "guest-2", response[1].ID)
}

func TestGuestHandler_ListGuests_ReturnsGuestsForApprovedGuest(t *testing.T) {
	guests := []*models.Guest{
		{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		},
	}

	svc := &mockGuestService{
		listGuestsAsGuestFunc: func(ctx context.Context, guestID, partyID string) ([]*models.Guest, error) {
			assert.Equal(t, "guest-1", guestID)
			assert.Equal(t, "party-1", partyID)
			return guests, nil
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/guests?guestId=guest-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response []*models.Guest
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "guest-1", response[0].ID)
}

func TestGuestHandler_ListGuests_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewGuestHandler(&mockGuestService{})

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/guests", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGuestHandler_ListGuests_ReturnsForbiddenForNonOwner(t *testing.T) {
	svc := &mockGuestService{
		listGuestsFunc: func(ctx context.Context, adminID, partyID string) ([]*models.Guest, error) {
			return nil, services.ErrUnauthorized
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/guests", nil)
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

// --- List Join Requests Tests ---

func TestGuestHandler_ListJoinRequests_ReturnsPendingGuests(t *testing.T) {
	guests := []*models.Guest{
		{
			ID:        "guest-3",
			PartyID:   "party-1",
			Username:  "charlie",
			Status:    models.GuestStatusPending,
			CreatedAt: time.Now(),
		},
	}

	svc := &mockGuestService{
		listJoinRequestsFunc: func(ctx context.Context, adminID, partyID string) ([]*models.Guest, error) {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "party-1", partyID)
			return guests, nil
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/join-requests", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response []*models.Guest
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "guest-3", response[0].ID)
	assert.Equal(t, models.GuestStatusPending, response[0].Status)
}

func TestGuestHandler_ListJoinRequests_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewGuestHandler(&mockGuestService{})

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/join-requests", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGuestHandler_ListJoinRequests_ReturnsForbiddenForNonOwner(t *testing.T) {
	svc := &mockGuestService{
		listJoinRequestsFunc: func(ctx context.Context, adminID, partyID string) ([]*models.Guest, error) {
			return nil, services.ErrUnauthorized
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/party-1/join-requests", nil)
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

// --- Approve Guest Tests ---

func TestGuestHandler_ApproveGuest_ReturnsOKOnSuccess(t *testing.T) {
	svc := &mockGuestService{
		approveGuestFunc: func(ctx context.Context, adminID, partyID, guestID string) error {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "party-1", partyID)
			assert.Equal(t, "guest-1", guestID)
			return nil
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/guests/guest-1/approve", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGuestHandler_ApproveGuest_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewGuestHandler(&mockGuestService{})

	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/guests/guest-1/approve", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGuestHandler_ApproveGuest_ReturnsForbiddenForNonOwner(t *testing.T) {
	svc := &mockGuestService{
		approveGuestFunc: func(ctx context.Context, adminID, partyID, guestID string) error {
			return services.ErrUnauthorized
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/guests/guest-1/approve", nil)
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGuestHandler_ApproveGuest_ReturnsNotFoundWhenNotExists(t *testing.T) {
	svc := &mockGuestService{
		approveGuestFunc: func(ctx context.Context, adminID, partyID, guestID string) error {
			return services.ErrNotFound
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/guests/guest-1/approve", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// --- Reject Guest Tests ---

func TestGuestHandler_RejectGuest_ReturnsOKOnSuccess(t *testing.T) {
	svc := &mockGuestService{
		rejectGuestFunc: func(ctx context.Context, adminID, partyID, guestID string) error {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "party-1", partyID)
			assert.Equal(t, "guest-1", guestID)
			return nil
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/guests/guest-1/reject", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGuestHandler_RejectGuest_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewGuestHandler(&mockGuestService{})

	req := httptest.NewRequest(http.MethodPut, "/api/parties/party-1/guests/guest-1/reject", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// --- Remove Guest Tests ---

func TestGuestHandler_RemoveGuest_ReturnsNoContentOnSuccess(t *testing.T) {
	svc := &mockGuestService{
		removeGuestFunc: func(ctx context.Context, adminID, partyID, guestID string) error {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "party-1", partyID)
			assert.Equal(t, "guest-1", guestID)
			return nil
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/parties/party-1/guests/guest-1", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestGuestHandler_RemoveGuest_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewGuestHandler(&mockGuestService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/parties/party-1/guests/guest-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGuestHandler_RemoveGuest_ReturnsNotFoundWhenNotExists(t *testing.T) {
	svc := &mockGuestService{
		removeGuestFunc: func(ctx context.Context, adminID, partyID, guestID string) error {
			return services.ErrNotFound
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/parties/party-1/guests/guest-1", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// --- Guest Status Tests ---

func TestGuestHandler_GetGuestStatus_ReturnsGuestStatus(t *testing.T) {
	guest := &models.Guest{
		ID:        "guest-1",
		PartyID:   "party-1",
		Username:  "alice",
		Status:    models.GuestStatusApproved,
		CreatedAt: time.Now(),
	}

	svc := &mockGuestService{
		getGuestStatusFunc: func(ctx context.Context, code, guestID string) (*models.Guest, error) {
			assert.Equal(t, "ABC234", code)
			assert.Equal(t, "guest-1", guestID)
			return guest, nil
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/ABC234/guest-status?guestId=guest-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.Guest
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "guest-1", response.ID)
	assert.Equal(t, "alice", response.Username)
	assert.Equal(t, models.GuestStatusApproved, response.Status)
}

func TestGuestHandler_GetGuestStatus_ReturnsBadRequestWithoutGuestId(t *testing.T) {
	handler := handlers.NewGuestHandler(&mockGuestService{})

	req := httptest.NewRequest(http.MethodGet, "/api/parties/ABC234/guest-status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGuestHandler_GetGuestStatus_ReturnsNotFoundWhenNotExists(t *testing.T) {
	svc := &mockGuestService{
		getGuestStatusFunc: func(ctx context.Context, code, guestID string) (*models.Guest, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewGuestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/ABC234/guest-status?guestId=guest-1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
