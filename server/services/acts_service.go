package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sipgate/eurovision-vote-party/server/models"
)

// ActsService provides access to Eurovision acts.
type ActsService interface {
	ListActs(eventType string) ([]models.Act, error)
}

type actsService struct {
	acts []models.Act
}

// actsFileWrapper represents the JSON structure of the acts file.
type actsFileWrapper struct {
	Acts []models.Act `json:"acts"`
}

// NewActsService loads acts from a JSON file and returns an ActsService.
// The JSON file must have the structure: {"acts": [...]}
func NewActsService(filePath string) (ActsService, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading acts file: %w", err)
	}

	var wrapper actsFileWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing acts file: %w", err)
	}

	return &actsService{acts: wrapper.Acts}, nil
}

func (s *actsService) ListActs(eventType string) ([]models.Act, error) {
	if eventType == "" {
		return s.acts, nil
	}

	et := models.EventType(eventType)
	if !et.IsValid() {
		return nil, ErrInvalidEventType
	}

	result := make([]models.Act, 0)
	for _, act := range s.acts {
		if act.EventType == et {
			result = append(result, act)
		}
	}

	return result, nil
}
