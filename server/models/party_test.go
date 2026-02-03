package models

import (
	"testing"
	"time"
)

func TestPartyValidate(t *testing.T) {
	base := Party{
		ID:        "party-1",
		Name:      "Eurovision Watch Party",
		Code:      "ABC123",
		EventType: EventGrandFinal,
		AdminID:   "admin-1",
		Status:    PartyStatusActive,
		CreatedAt: time.Now(),
	}

	t.Run("valid party", func(t *testing.T) {
		if err := base.Validate(); err != nil {
			t.Fatalf("expected validation to succeed, got error: %v", err)
		}
	})

	tests := map[string]Party{
		"missing name": {
			ID:        base.ID,
			Code:      base.Code,
			EventType: base.EventType,
			AdminID:   base.AdminID,
			Status:    base.Status,
			CreatedAt: base.CreatedAt,
		},
		"missing code": {
			ID:        base.ID,
			Name:      base.Name,
			EventType: base.EventType,
			AdminID:   base.AdminID,
			Status:    base.Status,
			CreatedAt: base.CreatedAt,
		},
		"invalid event type": {
			ID:        base.ID,
			Name:      base.Name,
			Code:      base.Code,
			EventType: EventType("invalid"),
			AdminID:   base.AdminID,
			Status:    base.Status,
			CreatedAt: base.CreatedAt,
		},
		"invalid status": {
			ID:        base.ID,
			Name:      base.Name,
			Code:      base.Code,
			EventType: base.EventType,
			AdminID:   base.AdminID,
			Status:    PartyStatus("invalid"),
			CreatedAt: base.CreatedAt,
		},
		"missing admin id": {
			ID:        base.ID,
			Name:      base.Name,
			Code:      base.Code,
			EventType: base.EventType,
			Status:    base.Status,
			CreatedAt: base.CreatedAt,
		},
		"zero created at": {
			ID:        base.ID,
			Name:      base.Name,
			Code:      base.Code,
			EventType: base.EventType,
			AdminID:   base.AdminID,
			Status:    base.Status,
		},
	}

	for name, party := range tests {
		t.Run(name, func(t *testing.T) {
			if err := party.Validate(); err == nil {
				t.Fatalf("expected validation to fail for %s", name)
			}
		})
	}
}
