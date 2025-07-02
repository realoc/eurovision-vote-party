package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler handles the health check endpoint
func HealthHandler(responseWriter http.ResponseWriter, r *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	fmt.Fprintf(responseWriter, `{"status": "healthy"}`)
}
