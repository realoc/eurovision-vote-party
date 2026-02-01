package handlers

import (
    "encoding/json"
    "net/http"
)

type HealthHandler struct{}

// NewHealthHandler constructs an http.Handler that reports API health.
func NewHealthHandler() http.Handler {
    return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    _ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

