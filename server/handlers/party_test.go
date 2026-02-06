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
	"github.com/sipgate/eurovision-vote-party/server/middleware"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

type mockPartyService struct {
	createPartyFunc        func(ctx context.Context, adminID string, req services.CreatePartyRequest) (*models.Party, error)
	getPartyByIDFunc       func(ctx context.Context, adminID, partyID string) (*models.Party, error)
	getPartyByCodeFunc     func(ctx context.Context, code string) (*models.Party, error)
	listPartiesByAdminFunc func(ctx context.Context, adminID string) ([]*models.Party, error)
	deletePartyFunc        func(ctx context.Context, adminID, partyID string) error
}

func (m *mockPartyService) CreateParty(ctx context.Context, adminID string, req services.CreatePartyRequest) (*models.Party, error) {
	if m.createPartyFunc != nil {
		return m.createPartyFunc(ctx, adminID, req)
	}
	return nil, nil
}

func (m *mockPartyService) GetPartyByID(ctx context.Context, adminID, partyID string) (*models.Party, error) {
	if m.getPartyByIDFunc != nil {
		return m.getPartyByIDFunc(ctx, adminID, partyID)
	}
	return nil, services.ErrNotFound
}

func (m *mockPartyService) GetPartyByCode(ctx context.Context, code string) (*models.Party, error) {
	if m.getPartyByCodeFunc != nil {
		return m.getPartyByCodeFunc(ctx, code)
	}
	return nil, services.ErrNotFound
}

func (m *mockPartyService) ListPartiesByAdmin(ctx context.Context, adminID string) ([]*models.Party, error) {
	if m.listPartiesByAdminFunc != nil {
		return m.listPartiesByAdminFunc(ctx, adminID)
	}
	return []*models.Party{}, nil
}

func (m *mockPartyService) DeleteParty(ctx context.Context, adminID, partyID string) error {
	if m.deletePartyFunc != nil {
		return m.deletePartyFunc(ctx, adminID, partyID)
	}
	return nil
}

func requestWithUserID(req *http.Request, userID string) *http.Request {
	ctx := middleware.WithUserID(req.Context(), userID)
	return req.WithContext(ctx)
}

// --- Create Party Tests ---

func TestPartyHandler_CreateParty_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	body := `{"name": "Test Party", "eventType": "grandfinal"}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPartyHandler_CreateParty_ReturnsBadRequestWithInvalidJSON(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	req := httptest.NewRequest(http.MethodPost, "/api/parties", bytes.NewBufferString("invalid json"))
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPartyHandler_CreateParty_ReturnsCreatedWithValidRequest(t *testing.T) {
	createdParty := &models.Party{
		ID:        "party-1",
		Name:      "Test Party",
		Code:      "ABC123",
		EventType: models.EventGrandFinal,
		AdminID:   "user-123",
		Status:    models.PartyStatusActive,
		CreatedAt: time.Now(),
	}

	svc := &mockPartyService{
		createPartyFunc: func(ctx context.Context, adminID string, req services.CreatePartyRequest) (*models.Party, error) {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "Test Party", req.Name)
			assert.Equal(t, models.EventGrandFinal, req.EventType)
			return createdParty, nil
		},
	}

	handler := handlers.NewPartyHandler(svc)

	body := `{"name": "Test Party", "eventType": "grandfinal"}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties", bytes.NewBufferString(body))
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.Party
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "party-1", response.ID)
	assert.Equal(t, "Test Party", response.Name)
	assert.Equal(t, "ABC123", response.Code)
}

func TestPartyHandler_CreateParty_ReturnsBadRequestWithMissingName(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	body := `{"name": "", "eventType": "grandfinal"}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties", bytes.NewBufferString(body))
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPartyHandler_CreateParty_ReturnsBadRequestWithWhitespaceName(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	body := `{"name": "   ", "eventType": "grandfinal"}`
	req := httptest.NewRequest(http.MethodPost, "/api/parties", bytes.NewBufferString(body))
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- List Parties Tests ---

func TestPartyHandler_ListParties_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	req := httptest.NewRequest(http.MethodGet, "/api/parties", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPartyHandler_ListParties_ReturnsPartiesList(t *testing.T) {
	parties := []*models.Party{
		{
			ID:        "party-1",
			Name:      "Party 1",
			Code:      "CODE01",
			EventType: models.EventGrandFinal,
			AdminID:   "user-123",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		},
		{
			ID:        "party-2",
			Name:      "Party 2",
			Code:      "CODE02",
			EventType: models.EventSemifinal1,
			AdminID:   "user-123",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		},
	}

	svc := &mockPartyService{
		listPartiesByAdminFunc: func(ctx context.Context, adminID string) ([]*models.Party, error) {
			assert.Equal(t, "user-123", adminID)
			return parties, nil
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response []*models.Party
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "party-1", response[0].ID)
	assert.Equal(t, "party-2", response[1].ID)
}

func TestPartyHandler_ListParties_ReturnsEmptyList(t *testing.T) {
	svc := &mockPartyService{
		listPartiesByAdminFunc: func(ctx context.Context, adminID string) ([]*models.Party, error) {
			return []*models.Party{}, nil
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []*models.Party
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 0)
}

// --- Get By Code Tests ---

func TestPartyHandler_GetByCode_ReturnsPartyWithoutAuth(t *testing.T) {
	party := &models.Party{
		ID:        "party-1",
		Name:      "Test Party",
		Code:      "ABC234",
		EventType: models.EventGrandFinal,
		AdminID:   "admin-1",
		Status:    models.PartyStatusActive,
		CreatedAt: time.Now(),
	}

	svc := &mockPartyService{
		getPartyByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
			assert.Equal(t, "ABC234", code)
			return party, nil
		},
	}

	handler := handlers.NewPartyHandler(svc)

	// No auth header - public endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/parties/ABC234", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should return public info only (no adminId or createdAt)
	assert.Equal(t, "party-1", response["id"])
	assert.Equal(t, "Test Party", response["name"])
	assert.Equal(t, "ABC234", response["code"])
	assert.Equal(t, "grandfinal", response["eventType"])
	assert.Equal(t, "active", response["status"])
	assert.NotContains(t, response, "adminId")
	assert.NotContains(t, response, "createdAt")
}

func TestPartyHandler_GetByCode_ReturnsNotFoundWhenNotExists(t *testing.T) {
	svc := &mockPartyService{
		getPartyByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/XYZ789", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// --- Get By ID Tests ---

func TestPartyHandler_GetByID_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	req := httptest.NewRequest(http.MethodGet, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPartyHandler_GetByID_ReturnsPartyForOwner(t *testing.T) {
	party := &models.Party{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Name:      "Test Party",
		Code:      "ABC123",
		EventType: models.EventGrandFinal,
		AdminID:   "user-123",
		Status:    models.PartyStatusActive,
		CreatedAt: time.Now(),
	}

	svc := &mockPartyService{
		getPartyByIDFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", partyID)
			return party, nil
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.Party
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response.ID)
	assert.Equal(t, "user-123", response.AdminID)
}

func TestPartyHandler_GetByID_ReturnsForbiddenWhenNotOwner(t *testing.T) {
	svc := &mockPartyService{
		getPartyByIDFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			return nil, services.ErrUnauthorized
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestPartyHandler_GetByID_ReturnsNotFoundWhenNotExists(t *testing.T) {
	svc := &mockPartyService{
		getPartyByIDFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// --- Delete Tests ---

func TestPartyHandler_Delete_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPartyHandler_Delete_ReturnsNoContentOnSuccess(t *testing.T) {
	deleteCalled := false
	svc := &mockPartyService{
		deletePartyFunc: func(ctx context.Context, adminID, partyID string) error {
			assert.Equal(t, "user-123", adminID)
			assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", partyID)
			deleteCalled = true
			return nil
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.True(t, deleteCalled)
}

func TestPartyHandler_Delete_ReturnsForbiddenWhenNotOwner(t *testing.T) {
	svc := &mockPartyService{
		deletePartyFunc: func(ctx context.Context, adminID, partyID string) error {
			return services.ErrUnauthorized
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	req = requestWithUserID(req, "other-user")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestPartyHandler_Delete_ReturnsNotFoundWhenNotExists(t *testing.T) {
	svc := &mockPartyService{
		deletePartyFunc: func(ctx context.Context, adminID, partyID string) error {
			return services.ErrNotFound
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// --- Routing Tests ---

func TestPartyHandler_ReturnsMethodNotAllowedForUnsupportedMethod(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	req := httptest.NewRequest(http.MethodPut, "/api/parties", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestPartyHandler_ReturnsMethodNotAllowedForDeleteWithoutID(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/parties", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestPartyHandler_ReturnsMethodNotAllowedForPostWithID(t *testing.T) {
	handler := handlers.NewPartyHandler(&mockPartyService{})

	req := httptest.NewRequest(http.MethodPost, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

// --- Code Detection Tests ---

func TestPartyHandler_DistinguishesBetweenCodeAndUUID(t *testing.T) {
	codeParty := &models.Party{
		ID:        "party-1",
		Name:      "Code Party",
		Code:      "ABC234",
		EventType: models.EventGrandFinal,
		AdminID:   "admin-1",
		Status:    models.PartyStatusActive,
		CreatedAt: time.Now(),
	}

	idParty := &models.Party{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Name:      "ID Party",
		Code:      "XYZ789",
		EventType: models.EventSemifinal1,
		AdminID:   "user-123",
		Status:    models.PartyStatusActive,
		CreatedAt: time.Now(),
	}

	svc := &mockPartyService{
		getPartyByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
			if code == "ABC234" {
				return codeParty, nil
			}
			return nil, services.ErrNotFound
		},
		getPartyByIDFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			if partyID == "550e8400-e29b-41d4-a716-446655440000" {
				return idParty, nil
			}
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewPartyHandler(svc)

	// Test code lookup (6 chars, uppercase letters and digits 2-9)
	t.Run("code lookup", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/parties/ABC234", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "Code Party", response["name"])
	})

	// Test UUID lookup
	t.Run("uuid lookup", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/parties/550e8400-e29b-41d4-a716-446655440000", nil)
		req = requestWithUserID(req, "user-123")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response models.Party
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "ID Party", response.Name)
	})
}

func TestPartyHandler_CodeWithLowercaseIsNotConsideredCode(t *testing.T) {
	svc := &mockPartyService{
		getPartyByIDFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			// Lowercase 6-char string should be treated as ID, not code
			assert.Equal(t, "abc123", partyID)
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/abc123", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should have tried ID lookup and got 404
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPartyHandler_CodeWithInvalidCharsIsNotConsideredCode(t *testing.T) {
	svc := &mockPartyService{
		getPartyByIDFunc: func(ctx context.Context, adminID, partyID string) (*models.Party, error) {
			// Code with invalid chars (0, 1, I, O, L) should be treated as ID
			assert.Equal(t, "ABC10L", partyID)
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewPartyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/parties/ABC10L", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
