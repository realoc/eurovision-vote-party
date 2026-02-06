package models

import "testing"

func TestUserValidate(t *testing.T) {
	base := User{
		ID:       "user-1",
		Username: "mango",
		Email:    "mango@example.com",
	}

	t.Run("valid user", func(t *testing.T) {
		if err := base.Validate(); err != nil {
			t.Fatalf("expected validation to succeed, got %v", err)
		}
	})

	tests := map[string]User{
		"missing id": {
			Username: base.Username,
			Email:    base.Email,
		},
		"missing username": {
			ID:    base.ID,
			Email: base.Email,
		},
		"missing email": {
			ID:       base.ID,
			Username: base.Username,
		},
	}

	for name, user := range tests {
		t.Run(name, func(t *testing.T) {
			if err := user.Validate(); err == nil {
				t.Fatalf("expected validation to fail for %s", name)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	validCases := []struct {
		name     string
		username string
	}{
		{"minimum length (3 chars)", "abc"},
		{"maximum length (30 chars)", "aaaaabbbbbcccccdddddeeeeefffff"},
		{"alphanumeric with underscores", "user_name_123"},
		{"all digits", "123"},
		{"all uppercase", "ABC"},
		{"mixed case", "AbC123"},
	}

	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateUsername(tc.username); err != nil {
				t.Fatalf("expected username %q to be valid, got %v", tc.username, err)
			}
		})
	}

	invalidCases := []struct {
		name     string
		username string
	}{
		{"empty string", ""},
		{"too short (2 chars)", "ab"},
		{"too long (31 chars)", "aaaaabbbbbcccccdddddeeeeeffffffg"},
		{"contains space", "user name"},
		{"contains @", "user@name"},
		{"contains dash", "user-name"},
		{"contains dot", "user.name"},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateUsername(tc.username); err == nil {
				t.Fatalf("expected username %q to be invalid", tc.username)
			}
		})
	}
}
