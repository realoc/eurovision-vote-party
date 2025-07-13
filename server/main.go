package main

import (
	"eurovision-app/http-handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create a custom ServeMux
	mux := http.NewServeMux()

	// Setup routes with the custom ServeMux
	handlers.SetupRoutes(mux)

	// Start the server with the custom ServeMux
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
