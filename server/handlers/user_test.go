package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/handlers"
	"github.com/sipgate/eurovision-vote-party/server/middleware"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

type mockUserService struct {
	upsertProfileFunc func(ctx context.Context, userID, email, username string) (*models.User, error)
	getProfileFunc    func(ctx context.Context, userID string) (*models.User, error)
}

func (m *mockUserService) UpsertProfile(ctx context.Context, userID, email, username string) (*models.User, error) {
	if m.upsertProfileFunc != nil {
		return m.upsertProfileFunc(ctx, userID, email, username)
	}
	return nil, nil
}

func (m *mockUserService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	if m.getProfileFunc != nil {
		return m.getProfileFunc(ctx, userID)
	}
	return nil, services.ErrNotFound
}

func requestWithUserIDAndEmail(req *http.Request, userID, email string) *http.Request {
	ctx := middleware.WithUserID(req.Context(), userID)
	ctx = middleware.WithUserEmail(ctx, email)
	return req.WithContext(ctx)
}

// --- PUT /api/users/profile Tests ---

func TestUserHandler_UpdateProfile_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})

	body := `{"username":"validuser"}`
	req := httptest.NewRequest(http.MethodPut, "/api/users/profile", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUserHandler_UpdateProfile_ReturnsBadRequestWithInvalidJSON(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})

	req := httptest.NewRequest(http.MethodPut, "/api/users/profile", bytes.NewBufferString("invalid json"))
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestUserHandler_UpdateProfile_ReturnsBadRequestWithInvalidUsername(t *testing.T) {
	svc := &mockUserService{
		upsertProfileFunc: func(ctx context.Context, userID, email, username string) (*models.User, error) {
			return nil, services.ErrInvalidUsername
		},
	}

	handler := handlers.NewUserHandler(svc)

	body := `{"username":"ab"}`
	req := httptest.NewRequest(http.MethodPut, "/api/users/profile", bytes.NewBufferString(body))
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_UpdateProfile_ReturnsOKWithValidRequest(t *testing.T) {
	returnedUser := &models.User{
		ID:       "user-123",
		Email:    "user@example.com",
		Username: "validuser",
	}

	svc := &mockUserService{
		upsertProfileFunc: func(ctx context.Context, userID, email, username string) (*models.User, error) {
			assert.Equal(t, "user-123", userID)
			assert.Equal(t, "user@example.com", email)
			assert.Equal(t, "validuser", username)
			return returnedUser, nil
		},
	}

	handler := handlers.NewUserHandler(svc)

	body := `{"username":"validuser"}`
	req := httptest.NewRequest(http.MethodPut, "/api/users/profile", bytes.NewBufferString(body))
	req = requestWithUserIDAndEmail(req, "user-123", "user@example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.User
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user-123", response.ID)
	assert.Equal(t, "validuser", response.Username)
	assert.Equal(t, "user@example.com", response.Email)
}

func TestUserHandler_UpdateProfile_PassesEmailFromContextToService(t *testing.T) {
	var capturedEmail string

	svc := &mockUserService{
		upsertProfileFunc: func(ctx context.Context, userID, email, username string) (*models.User, error) {
			capturedEmail = email
			return &models.User{
				ID:       userID,
				Email:    email,
				Username: username,
			}, nil
		},
	}

	handler := handlers.NewUserHandler(svc)

	body := `{"username":"validuser"}`
	req := httptest.NewRequest(http.MethodPut, "/api/users/profile", bytes.NewBufferString(body))
	req = requestWithUserIDAndEmail(req, "user-123", "user@example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "user@example.com", capturedEmail)
}

// --- GET /api/users/profile Tests ---

func TestUserHandler_GetProfile_ReturnsUnauthorizedWithoutAuth(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})

	req := httptest.NewRequest(http.MethodGet, "/api/users/profile", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUserHandler_GetProfile_ReturnsOKWhenProfileExists(t *testing.T) {
	returnedUser := &models.User{
		ID:       "user-123",
		Email:    "user@example.com",
		Username: "testuser",
	}

	svc := &mockUserService{
		getProfileFunc: func(ctx context.Context, userID string) (*models.User, error) {
			assert.Equal(t, "user-123", userID)
			return returnedUser, nil
		},
	}

	handler := handlers.NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/users/profile", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response models.User
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user-123", response.ID)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "user@example.com", response.Email)
}

func TestUserHandler_GetProfile_ReturnsNotFoundWhenProfileNotFound(t *testing.T) {
	svc := &mockUserService{
		getProfileFunc: func(ctx context.Context, userID string) (*models.User, error) {
			return nil, services.ErrNotFound
		},
	}

	handler := handlers.NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/users/profile", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Profile not found", response["error"])
}

func TestUserHandler_GetProfile_ReturnsInternalServerErrorOnServiceError(t *testing.T) {
	svc := &mockUserService{
		getProfileFunc: func(ctx context.Context, userID string) (*models.User, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := handlers.NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/users/profile", nil)
	req = requestWithUserID(req, "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- Method Routing Tests ---

func TestUserHandler_ReturnsMethodNotAllowedForPost(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})

	req := httptest.NewRequest(http.MethodPost, "/api/users/profile", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestUserHandler_ReturnsMethodNotAllowedForDelete(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/users/profile", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestUserHandler_ReturnsMethodNotAllowedForPatch(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})

	req := httptest.NewRequest(http.MethodPatch, "/api/users/profile", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
