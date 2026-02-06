package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
)

// VotePartyDAO defines the minimal party persistence operations needed by the vote service.
type VotePartyDAO interface {
	GetByID(ctx context.Context, id string) (*models.Party, error)
	UpdateStatus(ctx context.Context, id string, status models.PartyStatus) error
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
	ListByPartyID(ctx context.Context, partyID string) ([]*models.Vote, error)
}

// SubmitVoteRequest captures the data needed to submit or update a vote.
type SubmitVoteRequest struct {
	GuestID string
	Votes   map[int]string
}

// PartyResults contains the aggregated vote results for a party.
type PartyResults struct {
	PartyID     string              `json:"partyId"`
	PartyName   string              `json:"partyName"`
	TotalVoters int                 `json:"totalVoters"`
	Results     []models.VoteResult `json:"results"`
}

// VoteService defines the business logic operations for votes.
type VoteService interface {
	SubmitVote(ctx context.Context, adminID, partyID string, req SubmitVoteRequest) (*models.Vote, error)
	GetVotes(ctx context.Context, adminID, partyID, guestID string) (*models.Vote, error)
	UpdateVote(ctx context.Context, adminID, partyID string, req SubmitVoteRequest) (*models.Vote, error)
	EndVoting(ctx context.Context, adminID, partyID string) (*models.Party, error)
	GetResults(ctx context.Context, adminID, partyID string) (*PartyResults, error)
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

// EndVoting closes voting for a party.
func (s *voteService) EndVoting(ctx context.Context, adminID, partyID string) (*models.Party, error) {
	party, err := s.partyDAO.GetByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if adminID == "" {
		return nil, ErrUnauthorized
	}

	if party.AdminID != adminID {
		return nil, ErrUnauthorized
	}

	if party.Status != models.PartyStatusActive {
		return nil, ErrPartyClosed
	}

	if err := s.partyDAO.UpdateStatus(ctx, partyID, models.PartyStatusClosed); err != nil {
		return nil, err
	}

	party.Status = models.PartyStatusClosed
	return party, nil
}

// GetResults returns the aggregated vote results for a closed party.
func (s *voteService) GetResults(ctx context.Context, adminID, partyID string) (*PartyResults, error) {
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

	if party.Status != models.PartyStatusClosed {
		return nil, ErrVotingNotEnded
	}

	votes, err := s.voteDAO.ListByPartyID(ctx, partyID)
	if err != nil {
		return nil, err
	}

	acts, err := s.actsService.ListActs(string(party.EventType))
	if err != nil {
		return nil, err
	}

	// Sum points per act across all votes
	pointsByAct := make(map[string]int, len(acts))
	for _, act := range acts {
		pointsByAct[act.ID] = 0
	}
	for _, vote := range votes {
		for points, actID := range vote.Votes {
			pointsByAct[actID] += points
		}
	}

	// Build results with act metadata
	actMap := make(map[string]models.Act, len(acts))
	for _, act := range acts {
		actMap[act.ID] = act
	}

	results := make([]models.VoteResult, 0, len(acts))
	for _, act := range acts {
		results = append(results, models.VoteResult{
			ActID:       act.ID,
			Country:     act.Country,
			Artist:      act.Artist,
			Song:        act.Song,
			TotalPoints: pointsByAct[act.ID],
		})
	}

	// Sort by total points descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalPoints > results[j].TotalPoints
	})

	// Assign ranks with standard competition ranking (1,2,2,4)
	for i := range results {
		if i == 0 {
			results[i].Rank = 1
		} else if results[i].TotalPoints == results[i-1].TotalPoints {
			results[i].Rank = results[i-1].Rank
		} else {
			results[i].Rank = i + 1
		}
	}

	return &PartyResults{
		PartyID:     party.ID,
		PartyName:   party.Name,
		TotalVoters: len(votes),
		Results:     results,
	}, nil
}
