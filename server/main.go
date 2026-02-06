package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"

	"github.com/sipgate/eurovision-vote-party/server/handlers"
	"github.com/sipgate/eurovision-vote-party/server/middleware"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

func main() {
	ctx := context.Background()
	app := configureFirebaseAuth(ctx)
	firestoreClient := configureFirestore(ctx, app)
	defer firestoreClient.Close()

	partyDAO := persistence.NewFirestorePartyDAO(firestoreClient)
	partyService := services.NewPartyService(partyDAO)
	partyHandler := handlers.NewPartyHandler(partyService)

	guestDAO := persistence.NewFirestoreGuestDAO(firestoreClient)
	guestService := services.NewGuestService(guestDAO, partyDAO)
	guestHandler := handlers.NewGuestHandler(guestService)

	actsService, err := services.NewActsService("data/acts.json")
	if err != nil {
		log.Fatalf("failed to load acts data: %v", err)
	}
	actsHandler := handlers.NewActsHandler(actsService)

	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/parties/")
		if strings.Contains(path, "/") {
			guestHandler.ServeHTTP(w, r)
			return
		}
		partyHandler.ServeHTTP(w, r)
	})

	mux := http.NewServeMux()
	mux.Handle("/api/health", handlers.NewHealthHandler())
	mux.Handle("/api/acts", actsHandler)
	mux.Handle("/api/parties", middleware.AuthMiddleware(partyHandler))
	mux.Handle("/api/parties/", middleware.OptionalAuthMiddleware(apiHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("server is running on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func configureFirebaseAuth(ctx context.Context) *firebase.App {
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
	return app
}

func configureFirestore(ctx context.Context, app *firebase.App) *firestore.Client {
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("failed to initialise firestore client: %v", err)
	}
	return client
}
