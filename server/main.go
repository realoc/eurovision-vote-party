package main

import (
	"context"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"

	"github.com/sipgate/eurovision-vote-party/server/handlers"
	"github.com/sipgate/eurovision-vote-party/server/middleware"
)

func main() {
	ctx := context.Background()
	configureFirebaseAuth(ctx)

	mux := http.NewServeMux()
	mux.Handle("/api/health", handlers.NewHealthHandler())

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("server is running on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func configureFirebaseAuth(ctx context.Context) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	var cfg *firebase.Config
	if projectID != "" {
		cfg = &firebase.Config{
			ProjectID: projectID,
		}
	}

	app, err := firebase.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialise firebase app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("failed to initialise firebase auth client: %v", err)
	}

	middleware.SetTokenVerifier(authClient)
}
