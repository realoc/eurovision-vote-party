package models

import (
	"fmt"
	"strings"
	"time"
)

// PartyStatus describes the lifecycle state of a party.
type PartyStatus string

const (
	// PartyStatusActive indicates that a party is currently open for interactions.
	PartyStatusActive PartyStatus = "active"
	// PartyStatusClosed indicates that a party has concluded and is no longer active.
	PartyStatusClosed PartyStatus = "closed"
)

// IsValid reports whether the status is one of the supported values.
func (s PartyStatus) IsValid() bool {
	switch s {
	case PartyStatusActive, PartyStatusClosed:
		return true
	default:
		return false
	}
}

// EventType outlines which part of Eurovision the party refers to.
type EventType string

const (
	// EventSemifinal1 represents the first semifinal.
	EventSemifinal1 EventType = "semifinal1"
	// EventSemifinal2 represents the second semifinal.
	EventSemifinal2 EventType = "semifinal2"
	// EventGrandFinal represents the grand final.
	EventGrandFinal EventType = "grandfinal"
)

// IsValid reports whether the event type is recognised.
func (e EventType) IsValid() bool {
	switch e {
	case EventSemifinal1, EventSemifinal2, EventGrandFinal:
		return true
	default:
		return false
	}
}

// Party captures metadata about a Eurovision watch party.
type Party struct {
	ID        string      `firestore:"id" json:"id"`
	Name      string      `firestore:"name" json:"name"`
	Code      string      `firestore:"code" json:"code"`
	EventType EventType   `firestore:"eventType" json:"eventType"`
	AdminID   string      `firestore:"adminId" json:"adminId"`
	Status    PartyStatus `firestore:"status" json:"status"`
	CreatedAt time.Time   `firestore:"createdAt" json:"createdAt"`
}

// Validate ensures the party contains the required data.
func (p Party) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("party name is required")
	}
	if strings.TrimSpace(p.Code) == "" {
		return fmt.Errorf("party code is required")
	}
	if !p.EventType.IsValid() {
		return fmt.Errorf("event type %q is invalid", string(p.EventType))
	}
	if strings.TrimSpace(p.AdminID) == "" {
		return fmt.Errorf("admin id is required")
	}
	if !p.Status.IsValid() {
		return fmt.Errorf("party status %q is invalid", string(p.Status))
	}
	if p.CreatedAt.IsZero() {
		return fmt.Errorf("created at timestamp is required")
	}
	return nil
}
