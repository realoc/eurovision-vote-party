package handlers

import (
	"net/http"
	"os"
	"strings"
)

var allowedOrigins []string

func init() {
	origins := os.Getenv("ALLOWED_REQUEST_ORIGINS")
	if origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if len(allowedOrigins) > 0 {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
					break
				}
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// SetupRoutes configures all the HTTP routes for the application
// It requires a custom ServeMux to be provided
func SetupRoutes(mux *http.ServeMux) {
	if mux == nil {
		panic("ServeMux cannot be nil, a custom ServeMux must be provided")
	}

	// Define HTTP routes on provided mux
	mux.HandleFunc("/health", withCORS(HealthHandler))
	mux.HandleFunc("/party/create", withCORS(CreateParty))

	// Serve favicon.ico
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "data/favicon.ico")
	})

	// Serve static files (flags)
	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))
}
