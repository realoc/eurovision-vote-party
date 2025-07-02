package handlers

import (
	"encoding/json"
	"eurovision-app/service"
	"log"
	"net/http"
)

// Party represents a party with a name
type Party struct {
	Name string `json:"party_name"`
}

// PartyResponse represents the response for a created party
type PartyResponse struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// CreateParty handles the party creation endpoint
func CreateParty(responseWriter http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request body
	var party Party
	err := json.NewDecoder(r.Body).Decode(&party)
	if err != nil {
		http.Error(responseWriter, "Failed to parse request body", http.StatusBadRequest)
		log.Printf("Error parsing request body: %v", err)
		return
	}

	// Log the party name
	log.Printf("Received party: %s", party.Name)

	// Call the Create function from the service package with the party name
	id, password, err := service.Create(party.Name)
	if err != nil {
		http.Error(responseWriter, "Failed to create party", http.StatusInternalServerError)
		log.Printf("Error creating party: %v", err)
		return
	}

	// Create the response
	response := PartyResponse{
		ID:       id,
		Password: password,
	}

	// Return success response with id and password
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(response)
}
