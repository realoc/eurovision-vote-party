package main

import (
    "log"
    "net/http"

    "github.com/sipgate/eurovision-vote-party/server/handlers"
)

func main() {
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

