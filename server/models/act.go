package models

import (
	"fmt"
	"strings"
)

// Act represents an entry performing at Eurovision.
type Act struct {
	ID           string    `json:"id"`
	Country      string    `json:"country"`
	Artist       string    `json:"artist"`
	Song         string    `json:"song"`
	RunningOrder int       `json:"runningOrder"`
	EventType    EventType `json:"eventType"`
}

// Validate checks that the act contains essential information.
func (a Act) Validate() error {
	if strings.TrimSpace(a.Country) == "" {
		return fmt.Errorf("country is required")
	}
	if strings.TrimSpace(a.Artist) == "" {
		return fmt.Errorf("artist is required")
	}
	if strings.TrimSpace(a.Song) == "" {
		return fmt.Errorf("song is required")
	}
	if a.RunningOrder <= 0 {
		return fmt.Errorf("running order must be positive")
	}
	if !a.EventType.IsValid() {
		return fmt.Errorf("event type %q is invalid", string(a.EventType))
	}
	return nil
}
