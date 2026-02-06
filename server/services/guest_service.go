package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
)

// GuestDAO defines the persistence operations needed by the guest service.
type GuestDAO interface {
	Create(ctx context.Context, guest *models.Guest) error
	GetByID(ctx context.Context, id string) (*models.Guest, error)
	ListByPartyID(ctx context.Context, partyID string) ([]*models.Guest, error)
	ListByPartyIDAndStatus(ctx context.Context, partyID string, status models.GuestStatus) ([]*models.Guest, error)
	UpdateStatus(ctx context.Context, id string, status models.GuestStatus) error
	Delete(ctx context.Context, id string) error
	ExistsByPartyAndUsername(ctx context.Context, partyID, username string) (bool, error)
}

// GuestPartyDAO defines the minimal party persistence operations needed by the guest service.
type GuestPartyDAO interface {
	GetByID(ctx context.Context, id string) (*models.Party, error)
	GetByCode(ctx context.Context, code string) (*models.Party, error)
}

// GuestService defines the business logic operations for guests.
type GuestService interface {
	JoinParty(ctx context.Context, code, username string) (*models.Guest, error)
	ListGuests(ctx context.Context, adminID, partyID string) ([]*models.Guest, error)
	ListGuestsAsGuest(ctx context.Context, guestID, partyID string) ([]*models.Guest, error)
	ListJoinRequests(ctx context.Context, adminID, partyID string) ([]*models.Guest, error)
	ApproveGuest(ctx context.Context, adminID, partyID, guestID string) error
	RejectGuest(ctx context.Context, adminID, partyID, guestID string) error
	RemoveGuest(ctx context.Context, adminID, partyID, guestID string) error
	GetGuestStatus(ctx context.Context, code, guestID string) (*models.Guest, error)
}

// guestService is the default implementation.
type guestService struct {
	guestDAO GuestDAO
	partyDAO GuestPartyDAO
}

// NewGuestService creates a new GuestService.
func NewGuestService(guestDAO GuestDAO, partyDAO GuestPartyDAO) GuestService {
	return &guestService{guestDAO: guestDAO, partyDAO: partyDAO}
}

// JoinParty allows a guest to request joining a party by its public code.
func (s *guestService) JoinParty(ctx context.Context, code, username string) (*models.Guest, error) {
	party, err := s.partyDAO.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	exists, err := s.guestDAO.ExistsByPartyAndUsername(ctx, party.ID, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateUsername
	}

	guest := &models.Guest{
		ID:        uuid.New().String(),
		PartyID:   party.ID,
		Username:  username,
		Status:    models.GuestStatusPending,
		CreatedAt: time.Now(),
	}

	if err := guest.Validate(); err != nil {
		return nil, err
	}

	if err := s.guestDAO.Create(ctx, guest); err != nil {
		return nil, err
	}

	return guest, nil
}

// ListGuests returns all approved guests for a party, ensuring the requester is the admin.
func (s *guestService) ListGuests(ctx context.Context, adminID, partyID string) ([]*models.Guest, error) {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if party.AdminID != adminID {
		return nil, ErrUnauthorized
	}

	return s.guestDAO.ListByPartyIDAndStatus(ctx, partyID, models.GuestStatusApproved)
}

// ListGuestsAsGuest returns all approved guests for a party, ensuring the requester is an approved guest.
func (s *guestService) ListGuestsAsGuest(ctx context.Context, guestID, partyID string) ([]*models.Guest, error) {
	guest, err := s.guestDAO.GetByID(ctx, guestID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if guest.PartyID != partyID || guest.Status != models.GuestStatusApproved {
		return nil, ErrUnauthorized
	}

	return s.guestDAO.ListByPartyIDAndStatus(ctx, partyID, models.GuestStatusApproved)
}

// ListJoinRequests returns all pending guests for a party, ensuring the requester is the admin.
func (s *guestService) ListJoinRequests(ctx context.Context, adminID, partyID string) ([]*models.Guest, error) {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if party.AdminID != adminID {
		return nil, ErrUnauthorized
	}

	return s.guestDAO.ListByPartyIDAndStatus(ctx, partyID, models.GuestStatusPending)
}

// ApproveGuest approves a pending guest, ensuring the requester is the admin.
func (s *guestService) ApproveGuest(ctx context.Context, adminID, partyID, guestID string) error {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if party.AdminID != adminID {
		return ErrUnauthorized
	}

	guest, err := s.guestDAO.GetByID(ctx, guestID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if guest.PartyID != partyID || guest.Status != models.GuestStatusPending {
		return ErrNotFound
	}

	return s.guestDAO.UpdateStatus(ctx, guestID, models.GuestStatusApproved)
}

// RejectGuest rejects a pending guest, ensuring the requester is the admin.
func (s *guestService) RejectGuest(ctx context.Context, adminID, partyID, guestID string) error {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if party.AdminID != adminID {
		return ErrUnauthorized
	}

	guest, err := s.guestDAO.GetByID(ctx, guestID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if guest.PartyID != partyID || guest.Status != models.GuestStatusPending {
		return ErrNotFound
	}

	return s.guestDAO.UpdateStatus(ctx, guestID, models.GuestStatusRejected)
}

// RemoveGuest deletes a guest from a party, ensuring the requester is the admin.
func (s *guestService) RemoveGuest(ctx context.Context, adminID, partyID, guestID string) error {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if party.AdminID != adminID {
		return ErrUnauthorized
	}

	guest, err := s.guestDAO.GetByID(ctx, guestID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if guest.PartyID != partyID {
		return ErrNotFound
	}

	return s.guestDAO.Delete(ctx, guestID)
}

// GetGuestStatus retrieves a guest's status by party code and guest ID.
func (s *guestService) GetGuestStatus(ctx context.Context, code, guestID string) (*models.Guest, error) {
	party, err := s.partyDAO.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	guest, err := s.guestDAO.GetByID(ctx, guestID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if guest.PartyID != party.ID {
		return nil, ErrNotFound
	}

	return guest, nil
}
