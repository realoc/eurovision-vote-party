package services

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
)

const (
	codeAlphabet   = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	codeLength     = 6
	maxCodeRetries = 10
)

// CreatePartyRequest holds the data for creating a new party.
type CreatePartyRequest struct {
	Name      string
	EventType models.EventType
}

// PartyDAO defines the persistence operations needed by the service.
type PartyDAO interface {
	Create(ctx context.Context, party *models.Party) error
	GetByID(ctx context.Context, id string) (*models.Party, error)
	GetByCode(ctx context.Context, code string) (*models.Party, error)
	ListByAdminID(ctx context.Context, adminID string) ([]*models.Party, error)
	Delete(ctx context.Context, id string) error
	CodeExists(ctx context.Context, code string) (bool, error)
}

// PartyService defines the business logic operations for parties.
type PartyService interface {
	CreateParty(ctx context.Context, adminID string, req CreatePartyRequest) (*models.Party, error)
	GetPartyByID(ctx context.Context, adminID, partyID string) (*models.Party, error)
	GetPartyByCode(ctx context.Context, code string) (*models.Party, error)
	ListPartiesByAdmin(ctx context.Context, adminID string) ([]*models.Party, error)
	DeleteParty(ctx context.Context, adminID, partyID string) error
}

// partyService is the default implementation.
type partyService struct {
	dao PartyDAO
}

// NewPartyService creates a new PartyService.
func NewPartyService(dao PartyDAO) PartyService {
	return &partyService{dao: dao}
}

// generatePartyCode generates a random 6-character party code using a secure alphabet.
func generatePartyCode() (string, error) {
	b := make([]byte, codeLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = codeAlphabet[int(b[i])%len(codeAlphabet)]
	}
	return string(b), nil
}

// CreateParty creates a new party with a unique code.
func (s *partyService) CreateParty(ctx context.Context, adminID string, req CreatePartyRequest) (*models.Party, error) {
	var code string
	var err error

	// Try to generate a unique code with retry logic
	for i := 0; i < maxCodeRetries; i++ {
		code, err = generatePartyCode()
		if err != nil {
			return nil, err
		}

		exists, err := s.dao.CodeExists(ctx, code)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}

		// If this was the last attempt and code still exists, fail
		if i == maxCodeRetries-1 {
			return nil, errors.New("failed to generate unique party code after maximum retries")
		}
	}

	party := &models.Party{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Code:      code,
		EventType: req.EventType,
		AdminID:   adminID,
		Status:    models.PartyStatusActive,
		CreatedAt: time.Now(),
	}

	if err := party.Validate(); err != nil {
		return nil, err
	}

	if err := s.dao.Create(ctx, party); err != nil {
		return nil, err
	}

	return party, nil
}

// GetPartyByID retrieves a party by ID, ensuring the requester is the owner.
func (s *partyService) GetPartyByID(ctx context.Context, adminID, partyID string) (*models.Party, error) {
	party, err := s.dao.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if party.AdminID != adminID {
		return nil, ErrUnauthorized
	}

	return party, nil
}

// GetPartyByCode retrieves a party by its public code.
func (s *partyService) GetPartyByCode(ctx context.Context, code string) (*models.Party, error) {
	party, err := s.dao.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return party, nil
}

// ListPartiesByAdmin lists all parties owned by the given admin.
func (s *partyService) ListPartiesByAdmin(ctx context.Context, adminID string) ([]*models.Party, error) {
	return s.dao.ListByAdminID(ctx, adminID)
}

// DeleteParty deletes a party, ensuring the requester is the owner.
func (s *partyService) DeleteParty(ctx context.Context, adminID, partyID string) error {
	party, err := s.dao.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if party.AdminID != adminID {
		return ErrUnauthorized
	}

	if err := s.dao.Delete(ctx, partyID); err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}
