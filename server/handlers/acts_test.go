package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/handlers"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

type mockActsService struct {
	listActsFunc func(eventType string) ([]models.Act, error)
}

func (m *mockActsService) ListActs(eventType string) ([]models.Act, error) {
	if m.listActsFunc != nil {
		return m.listActsFunc(eventType)
	}
	return []models.Act{}, nil
}

func TestActsHandler_GET_Returns200WithActsList(t *testing.T) {
	acts := []models.Act{
		{
			ID:           "act-1",
			Country:      "Sweden",
			Artist:       "ABBA",
			Song:         "Waterloo",
			RunningOrder: 1,
			EventType:    models.EventGrandFinal,
		},
		{
			ID:           "act-2",
			Country:      "Italy",
			Artist:       "Måneskin",
			Song:         "Zitti e buoni",
			RunningOrder: 2,
			EventType:    models.EventGrandFinal,
		},
	}

	svc := &mockActsService{
		listActsFunc: func(eventType string) ([]models.Act, error) {
			return acts, nil
		},
	}

	handler := handlers.NewActsHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/acts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response struct {
		Acts []models.Act `json:"acts"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response.Acts, 2)
	assert.Equal(t, "act-1", response.Acts[0].ID)
	assert.Equal(t, "Sweden", response.Acts[0].Country)
	assert.Equal(t, "act-2", response.Acts[1].ID)
	assert.Equal(t, "Italy", response.Acts[1].Country)
}

func TestActsHandler_GET_WithEventQueryParam(t *testing.T) {
	var capturedEventType string

	svc := &mockActsService{
		listActsFunc: func(eventType string) ([]models.Act, error) {
			capturedEventType = eventType
			return []models.Act{}, nil
		},
	}

	handler := handlers.NewActsHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/acts?event=grandfinal", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "grandfinal", capturedEventType)
}

func TestActsHandler_GET_EmptyResultsReturnsEmptyArray(t *testing.T) {
	svc := &mockActsService{
		listActsFunc: func(eventType string) ([]models.Act, error) {
			return []models.Act{}, nil
		},
	}

	handler := handlers.NewActsHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/acts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response struct {
		Acts []models.Act `json:"acts"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotNil(t, response.Acts)
	assert.Len(t, response.Acts, 0)
}

func TestActsHandler_GET_InvalidEventReturns400(t *testing.T) {
	svc := &mockActsService{
		listActsFunc: func(eventType string) ([]models.Act, error) {
			return nil, services.ErrInvalidEventType
		},
	}

	handler := handlers.NewActsHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/acts?event=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestActsHandler_POST_Returns405(t *testing.T) {
	handler := handlers.NewActsHandler(&mockActsService{})

	req := httptest.NewRequest(http.MethodPost, "/api/acts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestActsHandler_DELETE_Returns405(t *testing.T) {
	handler := handlers.NewActsHandler(&mockActsService{})

	req := httptest.NewRequest(http.MethodDelete, "/api/acts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestActsHandler_GET_ServiceErrorReturns500(t *testing.T) {
	svc := &mockActsService{
		listActsFunc: func(eventType string) ([]models.Act, error) {
			return nil, errors.New("something")
		},
	}

	handler := handlers.NewActsHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/acts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
