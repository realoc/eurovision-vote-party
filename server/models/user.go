package models

import (
	"fmt"
	"strings"
)

// User represents an authenticated participant in the system.
type User struct {
	ID       string `firestore:"id" json:"id"`
	Username string `firestore:"username" json:"username"`
	Email    string `firestore:"email" json:"email"`
}

// Validate checks that the user record contains mandatory fields.
func (u User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return fmt.Errorf("id is required")
	}
	if strings.TrimSpace(u.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}
