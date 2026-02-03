package models

import (
	"testing"
	"time"
)

func TestVoteValidate(t *testing.T) {
	base := Vote{
		ID:      "vote-1",
		GuestID: "guest-1",
		PartyID: "party-1",
		Votes: map[int]string{
			12: "act-1",
			10: "act-2",
			8:  "act-3",
		},
		CreatedAt: time.Now(),
	}

	t.Run("valid vote", func(t *testing.T) {
		if err := base.Validate(); err != nil {
			t.Fatalf("expected validation to succeed, got %v", err)
		}
	})

	tests := map[string]Vote{
		"missing guest id": {
			ID:        base.ID,
			PartyID:   base.PartyID,
			Votes:     base.Votes,
			CreatedAt: base.CreatedAt,
		},
		"missing party id": {
			ID:        base.ID,
			GuestID:   base.GuestID,
			Votes:     base.Votes,
			CreatedAt: base.CreatedAt,
		},
		"empty votes": {
			ID:        base.ID,
			GuestID:   base.GuestID,
			PartyID:   base.PartyID,
			Votes:     map[int]string{},
			CreatedAt: base.CreatedAt,
		},
		"invalid point value": {
			ID:      base.ID,
			GuestID: base.GuestID,
			PartyID: base.PartyID,
			Votes: map[int]string{
				-1: "act-1",
			},
			CreatedAt: base.CreatedAt,
		},
		"missing act id": {
			ID:      base.ID,
			GuestID: base.GuestID,
			PartyID: base.PartyID,
			Votes: map[int]string{
				12: "",
			},
			CreatedAt: base.CreatedAt,
		},
		"zero created at": {
			ID:      base.ID,
			GuestID: base.GuestID,
			PartyID: base.PartyID,
			Votes:   base.Votes,
		},
	}

	for name, vote := range tests {
		t.Run(name, func(t *testing.T) {
			if err := vote.Validate(); err == nil {
				t.Fatalf("expected validation to fail for %s", name)
			}
		})
	}
}

func TestVoteResultValidate(t *testing.T) {
	base := VoteResult{
		ActID:       "act-1",
		Country:     "Sweden",
		Artist:      "Loreen",
		Song:        "Tattoo",
		TotalPoints: 583,
		Rank:        1,
	}

	t.Run("valid vote result", func(t *testing.T) {
		if err := base.Validate(); err != nil {
			t.Fatalf("expected validation to succeed, got %v", err)
		}
	})

	tests := map[string]VoteResult{
		"missing act id": {
			Country:     base.Country,
			Artist:      base.Artist,
			Song:        base.Song,
			TotalPoints: base.TotalPoints,
			Rank:        base.Rank,
		},
		"missing country": {
			ActID:       base.ActID,
			Artist:      base.Artist,
			Song:        base.Song,
			TotalPoints: base.TotalPoints,
			Rank:        base.Rank,
		},
		"missing artist": {
			ActID:       base.ActID,
			Country:     base.Country,
			Song:        base.Song,
			TotalPoints: base.TotalPoints,
			Rank:        base.Rank,
		},
		"missing song": {
			ActID:       base.ActID,
			Country:     base.Country,
			Artist:      base.Artist,
			TotalPoints: base.TotalPoints,
			Rank:        base.Rank,
		},
		"negative total points": {
			ActID:       base.ActID,
			Country:     base.Country,
			Artist:      base.Artist,
			Song:        base.Song,
			TotalPoints: -1,
			Rank:        base.Rank,
		},
		"negative rank": {
			ActID:       base.ActID,
			Country:     base.Country,
			Artist:      base.Artist,
			Song:        base.Song,
			TotalPoints: base.TotalPoints,
			Rank:        -1,
		},
	}

	for name, result := range tests {
		t.Run(name, func(t *testing.T) {
			if err := result.Validate(); err == nil {
				t.Fatalf("expected validation to fail for %s", name)
			}
		})
	}
}
