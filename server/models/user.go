package models

import (
	"fmt"
	"regexp"
	"strings"
)

// User represents an authenticated participant in the system.
type User struct {
	ID       string `firestore:"id" json:"id"`
	Username string `firestore:"username" json:"username"`
	Email    string `firestore:"email" json:"email"`
}

var validUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// ValidateUsername checks that the username is 3-30 characters and contains only
// alphanumeric characters and underscores.
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters")
	}
	if !validUsernameRegex.MatchString(username) {
		return fmt.Errorf("username must contain only alphanumeric characters and underscores")
	}
	return nil
}

// Validate checks that the user record contains mandatory fields.
func (u User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return fmt.Errorf("id is required")
	}
	if err := ValidateUsername(u.Username); err != nil {
		return err
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}
