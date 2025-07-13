package service

import (
	"testing"
)

func TestCreate(t *testing.T) {
	// Call the Create function
	id, password, err := Create("Test Party")

	// Check for errors
	if err != nil {
		t.Fatalf("Create() returned an error: %v", err)
	}

	// Check ID length
	if len(id) != 8 {
		t.Errorf("Expected ID length to be 8, got %d", len(id))
	}

	// Check password length
	if len(password) != 16 {
		t.Errorf("Expected password length to be 16, got %d", len(password))
	}

	// Call Create again to verify uniqueness
	id2, password2, err := Create("Test Party")
	if err != nil {
		t.Fatalf("Second Create() call returned an error: %v", err)
	}

	// Check that IDs are different
	if id == id2 {
		t.Errorf("Expected different IDs, but got the same ID twice: %s", id)
	}

	// Check that passwords are different
	if password == password2 {
		t.Errorf("Expected different passwords, but got the same password twice: %s", password)
	}
}
