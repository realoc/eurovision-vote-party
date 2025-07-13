package persistence

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// PartyDAO handles data access operations for parties
type PartyDAO struct {
	client *firestore.Client
	ctx    context.Context
}

// NewPartyDAO creates a new PartyDAO instance
func NewPartyDAO(ctx context.Context) (*PartyDAO, error) {
	// Initialize Firebase app
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json") // Update with actual path
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("Error initializing Firebase app: %v", err)
		return nil, err
	}

	// Get Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("Error getting Firestore client: %v", err)
		return nil, err
	}

	return &PartyDAO{
		client: client,
		ctx:    ctx,
	}, nil
}

// Add creates a new party document in Firestore
func (dao *PartyDAO) Add(id, name, password string) error {
	// Reference to the parties collection
	partiesCol := dao.client.Collection("parties")

	// Check if collection exists, create if not (this is handled automatically by Firestore)

	// Create a new document with the provided ID
	partyDoc := partiesCol.Doc(id)

	// Set the party data
	party := map[string]interface{}{
		"id":       id,
		"name":     name,
		"password": password,
	}

	// Set with merge option
	_, err := partyDoc.Set(dao.ctx, party, firestore.MergeAll)
	return err
}

// Party represents a party document in Firestore
type Party struct {
	ID       string `firestore:"id"`
	Name     string `firestore:"name"`
	Password string `firestore:"password"`
}

// Get retrieves a party document from Firestore by ID
func (dao *PartyDAO) Get(id string) (*Party, error) {
	// Reference to the party document
	partyDoc := dao.client.Collection("parties").Doc(id)

	// Get the document
	docSnap, err := partyDoc.Get(dao.ctx)
	if err != nil {
		return nil, err
	}

	// Check if document exists
	if !docSnap.Exists() {
		return nil, nil
	}

	// Parse document data into Party struct
	var party Party
	if err := docSnap.DataTo(&party); err != nil {
		return nil, err
	}

	return &party, nil
}
