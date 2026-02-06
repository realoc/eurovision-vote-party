package models

import (
	"fmt"
	"strings"
	"time"
)

// ValidPointValues defines the Eurovision scoring system point values.
var ValidPointValues = []int{12, 10, 8, 7, 6, 5, 4, 3, 2, 1}

// Vote records the points a guest awards to acts during a party.
type Vote struct {
	ID        string         `firestore:"id" json:"id"`
	GuestID   string         `firestore:"guestId" json:"guestId"`
	PartyID   string         `firestore:"partyId" json:"partyId"`
	Votes     map[int]string `firestore:"votes" json:"votes"` // points -> actID
	CreatedAt time.Time      `firestore:"createdAt" json:"createdAt"`
}

// Validate ensures that the vote capture is well-formed.
func (v Vote) Validate() error {
	if strings.TrimSpace(v.GuestID) == "" {
		return fmt.Errorf("guest id is required")
	}
	if strings.TrimSpace(v.PartyID) == "" {
		return fmt.Errorf("party id is required")
	}
	if len(v.Votes) != len(ValidPointValues) {
		return fmt.Errorf("exactly %d votes required, got %d", len(ValidPointValues), len(v.Votes))
	}
	for _, points := range ValidPointValues {
		actID, ok := v.Votes[points]
		if !ok {
			return fmt.Errorf("missing vote for point value %d", points)
		}
		if strings.TrimSpace(actID) == "" {
			return fmt.Errorf("act id is required for points %d", points)
		}
	}
	seen := make(map[string]bool, len(v.Votes))
	for points, actID := range v.Votes {
		if seen[actID] {
			return fmt.Errorf("duplicate act id %q for points %d", actID, points)
		}
		seen[actID] = true
	}
	if v.CreatedAt.IsZero() {
		return fmt.Errorf("created at timestamp is required")
	}
	return nil
}

// VoteResult represents the aggregated outcome for an act within a party.
type VoteResult struct {
	ActID       string `json:"actId"`
	Country     string `json:"country"`
	Artist      string `json:"artist"`
	Song        string `json:"song"`
	TotalPoints int    `json:"totalPoints"`
	Rank        int    `json:"rank"`
}

// Validate ensures the vote result contains the data required to present rankings.
func (r VoteResult) Validate() error {
	if strings.TrimSpace(r.ActID) == "" {
		return fmt.Errorf("act id is required")
	}
	if strings.TrimSpace(r.Country) == "" {
		return fmt.Errorf("country is required")
	}
	if strings.TrimSpace(r.Artist) == "" {
		return fmt.Errorf("artist is required")
	}
	if strings.TrimSpace(r.Song) == "" {
		return fmt.Errorf("song is required")
	}
	if r.TotalPoints < 0 {
		return fmt.Errorf("total points cannot be negative")
	}
	if r.Rank < 0 {
		return fmt.Errorf("rank cannot be negative")
	}
	return nil
}
