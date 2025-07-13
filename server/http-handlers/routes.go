package handlers

import (
	"net/http"
)

// SetupRoutes configures all the HTTP routes for the application
// It requires a custom ServeMux to be provided
func SetupRoutes(mux *http.ServeMux) {
	if mux == nil {
		panic("ServeMux cannot be nil, a custom ServeMux must be provided")
	}

	// Define HTTP routes on provided mux
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/party/create", CreateParty)

	// Serve favicon.ico
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "data/favicon.ico")
	})

	// Serve static files (flags)
	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))
}
