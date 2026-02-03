package models

import (
	"testing"
)

func TestActValidate(t *testing.T) {
	base := Act{
		ID:           "act-1",
		Country:      "Sweden",
		Artist:       "Loreen",
		Song:         "Tattoo",
		RunningOrder: 10,
		EventType:    EventGrandFinal,
	}

	t.Run("valid act", func(t *testing.T) {
		if err := base.Validate(); err != nil {
			t.Fatalf("expected validation to succeed, got %v", err)
		}
	})

	tests := map[string]Act{
		"missing country": {
			ID:           base.ID,
			Artist:       base.Artist,
			Song:         base.Song,
			RunningOrder: base.RunningOrder,
			EventType:    base.EventType,
		},
		"missing artist": {
			ID:           base.ID,
			Country:      base.Country,
			Song:         base.Song,
			RunningOrder: base.RunningOrder,
			EventType:    base.EventType,
		},
		"missing song": {
			ID:           base.ID,
			Country:      base.Country,
			Artist:       base.Artist,
			RunningOrder: base.RunningOrder,
			EventType:    base.EventType,
		},
		"non-positive running order": {
			ID:           base.ID,
			Country:      base.Country,
			Artist:       base.Artist,
			Song:         base.Song,
			RunningOrder: 0,
			EventType:    base.EventType,
		},
		"invalid event type": {
			ID:           base.ID,
			Country:      base.Country,
			Artist:       base.Artist,
			Song:         base.Song,
			RunningOrder: base.RunningOrder,
			EventType:    EventType("invalid"),
		},
	}

	for name, act := range tests {
		t.Run(name, func(t *testing.T) {
			if err := act.Validate(); err == nil {
				t.Fatalf("expected validation to fail for %s", name)
			}
		})
	}
}
