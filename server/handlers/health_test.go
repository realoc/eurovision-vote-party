package handlers_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/sipgate/eurovision-vote-party/server/handlers"
)

func TestHealthHandlerReturnsOKStatus(t *testing.T) {
    handler := handlers.NewHealthHandler()

    req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
    rec := httptest.NewRecorder()

    handler.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
    }

    if got := rec.Header().Get("Content-Type"); got != "application/json" {
        t.Fatalf("expected Content-Type application/json, got %q", got)
    }

    var payload map[string]string
    if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
        t.Fatalf("failed to unmarshal response body: %v", err)
    }

    if payload["status"] != "ok" {
        t.Fatalf("expected status \"ok\", got %q", payload["status"])
    }
}

