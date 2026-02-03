package models

import (
	"fmt"
	"strings"
	"time"
)

// GuestStatus enumerates allowed guest approval states.
type GuestStatus string

const (
	// GuestStatusPending indicates the guest has requested access and awaits approval.
	GuestStatusPending GuestStatus = "pending"
	// GuestStatusApproved indicates the guest has been approved to join the party.
	GuestStatusApproved GuestStatus = "approved"
	// GuestStatusRejected indicates the guest request was rejected.
	GuestStatusRejected GuestStatus = "rejected"
)

// IsValid reports whether the status is supported.
func (s GuestStatus) IsValid() bool {
	switch s {
	case GuestStatusPending, GuestStatusApproved, GuestStatusRejected:
		return true
	default:
		return false
	}
}

// Guest captures a party participant.
type Guest struct {
	ID        string      `firestore:"id" json:"id"`
	PartyID   string      `firestore:"partyId" json:"partyId"`
	Username  string      `firestore:"username" json:"username"`
	Status    GuestStatus `firestore:"status" json:"status"`
	CreatedAt time.Time   `firestore:"createdAt" json:"createdAt"`
}

// Validate ensures guest data adheres to expected constraints.
func (g Guest) Validate() error {
	if strings.TrimSpace(g.PartyID) == "" {
		return fmt.Errorf("party id is required")
	}
	if strings.TrimSpace(g.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if !g.Status.IsValid() {
		return fmt.Errorf("guest status %q is invalid", string(g.Status))
	}
	if g.CreatedAt.IsZero() {
		return fmt.Errorf("created at timestamp is required")
	}
	return nil
}
