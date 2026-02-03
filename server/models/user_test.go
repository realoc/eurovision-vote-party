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
