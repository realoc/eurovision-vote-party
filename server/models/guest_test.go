package models

import (
	"testing"
	"time"
)

func TestGuestValidate(t *testing.T) {
	base := Guest{
		ID:        "guest-1",
		PartyID:   "party-1",
		Username:  "Dana International",
		Status:    GuestStatusApproved,
		CreatedAt: time.Now(),
	}

	t.Run("valid guest", func(t *testing.T) {
		if err := base.Validate(); err != nil {
			t.Fatalf("expected validation to succeed, got %v", err)
		}
	})

	tests := map[string]Guest{
		"missing party id": {
			ID:        base.ID,
			Username:  base.Username,
			Status:    base.Status,
			CreatedAt: base.CreatedAt,
		},
		"missing username": {
			ID:        base.ID,
			PartyID:   base.PartyID,
			Status:    base.Status,
			CreatedAt: base.CreatedAt,
		},
		"invalid status": {
			ID:        base.ID,
			PartyID:   base.PartyID,
			Username:  base.Username,
			Status:    GuestStatus("invalid"),
			CreatedAt: base.CreatedAt,
		},
		"zero created at": {
			ID:       base.ID,
			PartyID:  base.PartyID,
			Username: base.Username,
			Status:   base.Status,
		},
	}

	for name, guest := range tests {
		t.Run(name, func(t *testing.T) {
			if err := guest.Validate(); err == nil {
				t.Fatalf("expected validation to fail for %s", name)
			}
		})
	}
}
