package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"eurovision-app/persistence"
	"log"
)

// Create generates a unique party ID and password and stores it in Firestore
// Parameters:
// - name: the name of the party
// Returns:
// - id: a unique ID with 8 characters length
// - password: a generated password with 16 characters length
func Create(name string) (id string, password string, err error) {
	// Generate random bytes for ID (8 characters = 6 bytes when base64 encoded)
	idBytes := make([]byte, 6)
	_, err = rand.Read(idBytes)
	if err != nil {
		return "", "", err
	}

	// Generate random bytes for password (16 characters = 12 bytes when base64 encoded)
	passwordBytes := make([]byte, 12)
	_, err = rand.Read(passwordBytes)
	if err != nil {
		return "", "", err
	}

	// Encode to base64 and trim to required length
	id = base64.URLEncoding.EncodeToString(idBytes)[:8]
	password = base64.URLEncoding.EncodeToString(passwordBytes)[:16]

	// Store the party in Firestore
	ctx := context.Background()
	partyDAO, err := persistence.NewPartyDAO(ctx)
	if err != nil {
		log.Printf("Error creating PartyDAO: %v", err)
		return id, password, err
	}

	// Use the provided name
	err = partyDAO.Add(id, name, password)
	if err != nil {
		log.Printf("Error adding party to Firestore: %v", err)
		return id, password, err
	}

	return id, password, nil
}
