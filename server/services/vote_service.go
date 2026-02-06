package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
)

// VotePartyDAO defines the minimal party persistence operations needed by the vote service.
type VotePartyDAO interface {
	GetByID(ctx context.Context, id string) (*models.Party, error)
}

// VoteGuestDAO defines the minimal guest persistence operations needed by the vote service.
type VoteGuestDAO interface {
	GetByID(ctx context.Context, id string) (*models.Guest, error)
}

// VoteActsService defines the acts service operations needed by the vote service.
type VoteActsService interface {
	ListActs(eventType string) ([]models.Act, error)
}

// VoteDAO defines the persistence operations needed by the vote service.
type VoteDAO interface {
	Create(ctx context.Context, vote *models.Vote) error
	GetByGuestAndParty(ctx context.Context, guestID, partyID string) (*models.Vote, error)
	Update(ctx context.Context, vote *models.Vote) error
}

// SubmitVoteRequest captures the data needed to submit or update a vote.
type SubmitVoteRequest struct {
	GuestID string
	Votes   map[int]string
}

// VoteService defines the business logic operations for votes.
type VoteService interface {
	SubmitVote(ctx context.Context, adminID, partyID string, req SubmitVoteRequest) (*models.Vote, error)
	GetVotes(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error)
	UpdateVote(ctx context.Context, adminID, partyID string, req SubmitVoteRequest) (*models.Vote, error)
}

// voteService is the default implementation.
type voteService struct {
	voteDAO     VoteDAO
	partyDAO    VotePartyDAO
	guestDAO    VoteGuestDAO
	actsService VoteActsService
}

// NewVoteService creates a new VoteService.
func NewVoteService(voteDAO VoteDAO, partyDAO VotePartyDAO, guestDAO VoteGuestDAO, actsService VoteActsService) VoteService {
	return &voteService{
		voteDAO:     voteDAO,
		partyDAO:    partyDAO,
		guestDAO:    guestDAO,
		actsService: actsService,
	}
}

// SubmitVote creates a new vote for a guest in a party.
func (s *voteService) SubmitVote(ctx context.Context, adminID, partyID string, req SubmitVoteRequest) (*models.Vote, error) {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if party.Status != models.PartyStatusActive {
		return nil, ErrPartyClosed
	}

	if adminID != "" {
		if party.AdminID != adminID {
			return nil, ErrUnauthorized
		}
	}

	guest, err := s.guestDAO.GetByID(ctx, req.GuestID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if guest.PartyID != partyID || guest.Status != models.GuestStatusApproved {
		return nil, ErrGuestNotApproved
	}

	acts, err := s.actsService.ListActs(string(party.EventType))
	if err != nil {
		return nil, err
	}

	actIDs := make(map[string]bool, len(acts))
	for _, act := range acts {
		actIDs[act.ID] = true
	}

	for _, actID := range req.Votes {
		if !actIDs[actID] {
			return nil, ErrInvalidVotes
		}
	}

	_, err = s.voteDAO.GetByGuestAndParty(ctx, req.GuestID, partyID)
	if err == nil {
		return nil, ErrVoteAlreadyExists
	}
	if !errors.Is(err, persistence.ErrNotFound) {
		return nil, err
	}

	vote := &models.Vote{
		ID:        uuid.New().String(),
		GuestID:   req.GuestID,
		PartyID:   partyID,
		Votes:     req.Votes,
		CreatedAt: time.Now(),
	}

	if err := vote.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidVotes, err)
	}

	if err := s.voteDAO.Create(ctx, vote); err != nil {
		return nil, err
	}

	return vote, nil
}

// GetVotes retrieves a guest's vote for a party.
func (s *voteService) GetVotes(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error) {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if adminID != "" {
		if party.AdminID != adminID {
			return nil, ErrUnauthorized
		}
	}

	vote, err := s.voteDAO.GetByGuestAndParty(ctx, guestID, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return vote, nil
}

// UpdateVote updates an existing vote for a guest in a party.
func (s *voteService) UpdateVote(ctx context.Context, adminID, partyID string, req SubmitVoteRequest) (*models.Vote, error) {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if party.Status != models.PartyStatusActive {
		return nil, ErrPartyClosed
	}

	if adminID != "" {
		if party.AdminID != adminID {
			return nil, ErrUnauthorized
		}
	}

	guest, err := s.guestDAO.GetByID(ctx, req.GuestID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if guest.PartyID != partyID || guest.Status != models.GuestStatusApproved {
		return nil, ErrGuestNotApproved
	}

	acts, err := s.actsService.ListActs(string(party.EventType))
	if err != nil {
		return nil, err
	}

	actIDs := make(map[string]bool, len(acts))
	for _, act := range acts {
		actIDs[act.ID] = true
	}

	for _, actID := range req.Votes {
		if !actIDs[actID] {
			return nil, ErrInvalidVotes
		}
	}

	existingVote, err := s.voteDAO.GetByGuestAndParty(ctx, req.GuestID, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	existingVote.Votes = req.Votes

	if err := existingVote.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidVotes, err)
	}

	if err := s.voteDAO.Update(ctx, existingVote); err != nil {
		return nil, err
	}

	return existingVote, nil
}
