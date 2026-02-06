package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

// mockVoteDAO mocks the VoteDAO interface used by the vote service.
type mockVoteDAO struct {
	createFunc             func(ctx context.Context, vote *models.Vote) error
	getByGuestAndPartyFunc func(ctx context.Context, guestID, partyID string) (*models.Vote, error)
	updateFunc             func(ctx context.Context, vote *models.Vote) error
}

func (m *mockVoteDAO) Create(ctx context.Context, vote *models.Vote) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, vote)
	}
	return nil
}

func (m *mockVoteDAO) GetByGuestAndParty(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
	if m.getByGuestAndPartyFunc != nil {
		return m.getByGuestAndPartyFunc(ctx, guestID, partyID)
	}
	return nil, persistence.ErrNotFound
}

func (m *mockVoteDAO) Update(ctx context.Context, vote *models.Vote) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, vote)
	}
	return nil
}

// mockVotePartyDAO mocks the VotePartyDAO interface used by the vote service.
type mockVotePartyDAO struct {
	getByIDFunc func(ctx context.Context, id string) (*models.Party, error)
}

func (m *mockVotePartyDAO) GetByID(ctx context.Context, id string) (*models.Party, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, persistence.ErrNotFound
}

// mockVoteGuestDAO mocks the VoteGuestDAO interface used by the vote service.
type mockVoteGuestDAO struct {
	getByIDFunc func(ctx context.Context, id string) (*models.Guest, error)
}

func (m *mockVoteGuestDAO) GetByID(ctx context.Context, id string) (*models.Guest, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, persistence.ErrNotFound
}

// mockVoteActsService mocks the VoteActsService interface used by the vote service.
type mockVoteActsService struct {
	listActsFunc func(eventType string) ([]models.Act, error)
}

func (m *mockVoteActsService) ListActs(eventType string) ([]models.Act, error) {
	if m.listActsFunc != nil {
		return m.listActsFunc(eventType)
	}
	return []models.Act{}, nil
}

func testActs() []models.Act {
	acts := make([]models.Act, 10)
	for i := 0; i < 10; i++ {
		acts[i] = models.Act{
			ID:           fmt.Sprintf("act-%d", i+1),
			Country:      fmt.Sprintf("Country %d", i+1),
			Artist:       fmt.Sprintf("Artist %d", i+1),
			Song:         fmt.Sprintf("Song %d", i+1),
			RunningOrder: i + 1,
			EventType:    models.EventGrandFinal,
		}
	}
	return acts
}

func validVotes() map[int]string {
	return map[int]string{
		12: "act-1",
		10: "act-2",
		8:  "act-3",
		7:  "act-4",
		6:  "act-5",
		5:  "act-6",
		4:  "act-7",
		3:  "act-8",
		2:  "act-9",
		1:  "act-10",
	}
}

func TestVoteService_SubmitVote(t *testing.T) {
	t.Run("submits vote successfully as admin", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		var createdVote *models.Vote
		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				return nil, persistence.ErrNotFound
			},
			createFunc: func(ctx context.Context, vote *models.Vote) error {
				createdVote = vote
				return nil
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return approvedGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{
			listActsFunc: func(eventType string) ([]models.Act, error) {
				return testActs(), nil
			},
		}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		require.NoError(t, err)
		assert.NotEmpty(t, vote.ID)
		assert.Equal(t, "guest-1", vote.GuestID)
		assert.Equal(t, "party-1", vote.PartyID)
		assert.Equal(t, validVotes(), vote.Votes)
		assert.False(t, vote.CreatedAt.IsZero())
		assert.Equal(t, createdVote, vote)
	})

	t.Run("submits vote successfully without admin", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				return nil, persistence.ErrNotFound
			},
			createFunc: func(ctx context.Context, vote *models.Vote) error {
				return nil
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return approvedGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{
			listActsFunc: func(eventType string) ([]models.Act, error) {
				return testActs(), nil
			},
		}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		require.NoError(t, err)
		assert.NotEmpty(t, vote.ID)
		assert.Equal(t, "guest-1", vote.GuestID)
		assert.Equal(t, "party-1", vote.PartyID)
	})

	t.Run("returns ErrNotFound when party not found", func(t *testing.T) {
		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return nil, persistence.ErrNotFound
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "nonexistent", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrPartyClosed when party is closed", func(t *testing.T) {
		closedParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusClosed,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return closedParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrPartyClosed)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrUnauthorized when admin does not own party", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "other-admin", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrNotFound when guest not found", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "nonexistent",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrGuestNotApproved when guest is pending", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		pendingGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusPending,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return pendingGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrGuestNotApproved)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrGuestNotApproved when guest belongs to different party", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		guestFromOtherParty := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-2",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return guestFromOtherParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrGuestNotApproved)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrInvalidVotes when act ID not in event type", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		invalidVotes := validVotes()
		invalidVotes[12] = "unknown-act"

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return approvedGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{
			listActsFunc: func(eventType string) ([]models.Act, error) {
				return testActs(), nil
			},
		}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   invalidVotes,
		})

		assert.ErrorIs(t, err, services.ErrInvalidVotes)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrVoteAlreadyExists when vote exists", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		existingVote := &models.Vote{
			ID:        "vote-1",
			GuestID:   "guest-1",
			PartyID:   "party-1",
			Votes:     validVotes(),
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				return existingVote, nil
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return approvedGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{
			listActsFunc: func(eventType string) ([]models.Act, error) {
				return testActs(), nil
			},
		}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrVoteAlreadyExists)
		assert.Nil(t, vote)
	})

	t.Run("propagates DAO create error", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		daoErr := errors.New("firestore unavailable")
		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				return nil, persistence.ErrNotFound
			},
			createFunc: func(ctx context.Context, vote *models.Vote) error {
				return daoErr
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return approvedGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{
			listActsFunc: func(eventType string) ([]models.Act, error) {
				return testActs(), nil
			},
		}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.SubmitVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, daoErr)
		assert.Nil(t, vote)
	})
}

func TestVoteService_GetVotes(t *testing.T) {
	t.Run("returns vote successfully as admin", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		existingVote := &models.Vote{
			ID:        "vote-1",
			GuestID:   "guest-1",
			PartyID:   "party-1",
			Votes:     validVotes(),
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				if guestID == "guest-1" && partyID == "party-1" {
					return existingVote, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.GetVotes(ctx, "admin-1", "party-1", "guest-1")

		require.NoError(t, err)
		assert.Equal(t, existingVote, vote)
	})

	t.Run("returns vote successfully without admin", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		existingVote := &models.Vote{
			ID:        "vote-1",
			GuestID:   "guest-1",
			PartyID:   "party-1",
			Votes:     validVotes(),
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				if guestID == "guest-1" && partyID == "party-1" {
					return existingVote, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.GetVotes(ctx, "", "party-1", "guest-1")

		require.NoError(t, err)
		assert.Equal(t, existingVote, vote)
	})

	t.Run("returns ErrNotFound when party not found", func(t *testing.T) {
		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return nil, persistence.ErrNotFound
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.GetVotes(ctx, "admin-1", "nonexistent", "guest-1")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrUnauthorized when admin does not own party", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.GetVotes(ctx, "other-admin", "party-1", "guest-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrNotFound when vote not found", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.GetVotes(ctx, "admin-1", "party-1", "guest-1")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, vote)
	})
}

func TestVoteService_UpdateVote(t *testing.T) {
	t.Run("updates vote successfully", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		existingVote := &models.Vote{
			ID:        "vote-1",
			GuestID:   "guest-1",
			PartyID:   "party-1",
			Votes:     validVotes(),
			CreatedAt: time.Now(),
		}

		updatedVotes := map[int]string{
			12: "act-10",
			10: "act-9",
			8:  "act-8",
			7:  "act-7",
			6:  "act-6",
			5:  "act-5",
			4:  "act-4",
			3:  "act-3",
			2:  "act-2",
			1:  "act-1",
		}

		var updatedVote *models.Vote
		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				if guestID == "guest-1" && partyID == "party-1" {
					return existingVote, nil
				}
				return nil, persistence.ErrNotFound
			},
			updateFunc: func(ctx context.Context, vote *models.Vote) error {
				updatedVote = vote
				return nil
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return approvedGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{
			listActsFunc: func(eventType string) ([]models.Act, error) {
				return testActs(), nil
			},
		}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.UpdateVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   updatedVotes,
		})

		require.NoError(t, err)
		assert.Equal(t, "vote-1", vote.ID)
		assert.Equal(t, updatedVotes, vote.Votes)
		assert.Equal(t, updatedVote, vote)
	})

	t.Run("returns ErrPartyClosed when party is closed", func(t *testing.T) {
		closedParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusClosed,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return closedParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{}
		actsService := &mockVoteActsService{}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.UpdateVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrPartyClosed)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrNotFound when no existing vote", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		voteDAO := &mockVoteDAO{
			getByGuestAndPartyFunc: func(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return approvedGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{
			listActsFunc: func(eventType string) ([]models.Act, error) {
				return testActs(), nil
			},
		}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.UpdateVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   validVotes(),
		})

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, vote)
	})

	t.Run("returns ErrInvalidVotes when act ID not in event type", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		invalidVotes := validVotes()
		invalidVotes[12] = "unknown-act"

		voteDAO := &mockVoteDAO{}
		partyDAO := &mockVotePartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}
		guestDAO := &mockVoteGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return approvedGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		actsService := &mockVoteActsService{
			listActsFunc: func(eventType string) ([]models.Act, error) {
				return testActs(), nil
			},
		}

		svc := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
		ctx := context.Background()

		vote, err := svc.UpdateVote(ctx, "admin-1", "party-1", services.SubmitVoteRequest{
			GuestID: "guest-1",
			Votes:   invalidVotes,
		})

		assert.ErrorIs(t, err, services.ErrInvalidVotes)
		assert.Nil(t, vote)
	})
}
