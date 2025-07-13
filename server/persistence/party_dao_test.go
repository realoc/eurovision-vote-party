package persistence

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	firebase "firebase.google.com/go/v4"
)

// isEmulatorRunning checks if the Firebase emulator is running on the specified port
func isEmulatorRunning() bool {
	conn, err := net.DialTimeout("tcp", "localhost:3001", 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// TestAddAndGet tests the round trip of writing and reading a document using the Firebase emulator
func TestAddAndGet(t *testing.T) {
	// Fail test if emulator is not running
	if !isEmulatorRunning() {
		t.Fatal("Firestore emulator is not running. Start it with 'firebase emulators:start'")
	}

	// Set environment variable to use the emulator
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:3001")
	defer os.Unsetenv("FIRESTORE_EMULATOR_HOST")

	// Create a context
	ctx := context.Background()

	// Initialize Firebase app with project ID (required even for emulator)
	config := &firebase.Config{ProjectID: "test-project"}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		t.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	// Get Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		t.Fatalf("Failed to get Firestore client: %v", err)
	}
	defer client.Close()

	// Create a PartyDAO instance with the emulator client
	dao := &PartyDAO{
		client: client,
		ctx:    ctx,
	}

	// Test data
	id := "test-id-123"
	name := "Test Party"
	password := "test-password"

	// Add a party document
	err = dao.Add(id, name, password)
	if err != nil {
		t.Fatalf("Failed to add party document: %v", err)
	}

	// Get the party document
	party, err := dao.Get(id)
	if err != nil {
		t.Fatalf("Failed to get party document: %v", err)
	}

	// Verify the document was retrieved
	if party == nil {
		t.Fatal("Party document not found")
	}

	// Verify the document data
	if party.ID != id {
		t.Errorf("Expected ID %s, got %s", id, party.ID)
	}
	if party.Name != name {
		t.Errorf("Expected Name %s, got %s", name, party.Name)
	}
	if party.Password != password {
		t.Errorf("Expected Password %s, got %s", password, party.Password)
	}

	// Clean up - delete the test document
	_, err = client.Collection("parties").Doc(id).Delete(ctx)
	if err != nil {
		t.Logf("Warning: Failed to delete test document: %v", err)
	}
}
