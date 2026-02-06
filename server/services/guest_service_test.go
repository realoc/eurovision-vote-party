package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

// mockGuestDAO mocks the GuestDAO interface used by the guest service.
type mockGuestDAO struct {
	createFunc                func(ctx context.Context, guest *models.Guest) error
	getByIDFunc               func(ctx context.Context, id string) (*models.Guest, error)
	listByPartyIDFunc         func(ctx context.Context, partyID string) ([]*models.Guest, error)
	listByPartyIDAndStatusFunc func(ctx context.Context, partyID string, status models.GuestStatus) ([]*models.Guest, error)
	updateStatusFunc          func(ctx context.Context, id string, status models.GuestStatus) error
	deleteFunc                func(ctx context.Context, id string) error
	existsByPartyAndUsernameFunc func(ctx context.Context, partyID, username string) (bool, error)
}

func (m *mockGuestDAO) Create(ctx context.Context, guest *models.Guest) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, guest)
	}
	return nil
}

func (m *mockGuestDAO) GetByID(ctx context.Context, id string) (*models.Guest, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, persistence.ErrNotFound
}

func (m *mockGuestDAO) ListByPartyID(ctx context.Context, partyID string) ([]*models.Guest, error) {
	if m.listByPartyIDFunc != nil {
		return m.listByPartyIDFunc(ctx, partyID)
	}
	return []*models.Guest{}, nil
}

func (m *mockGuestDAO) ListByPartyIDAndStatus(ctx context.Context, partyID string, status models.GuestStatus) ([]*models.Guest, error) {
	if m.listByPartyIDAndStatusFunc != nil {
		return m.listByPartyIDAndStatusFunc(ctx, partyID, status)
	}
	return []*models.Guest{}, nil
}

func (m *mockGuestDAO) UpdateStatus(ctx context.Context, id string, status models.GuestStatus) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *mockGuestDAO) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockGuestDAO) ExistsByPartyAndUsername(ctx context.Context, partyID, username string) (bool, error) {
	if m.existsByPartyAndUsernameFunc != nil {
		return m.existsByPartyAndUsernameFunc(ctx, partyID, username)
	}
	return false, nil
}

// mockGuestPartyDAO mocks the minimal PartyDAO interface used by the guest service.
type mockGuestPartyDAO struct {
	getByIDFunc   func(ctx context.Context, id string) (*models.Party, error)
	getByCodeFunc func(ctx context.Context, code string) (*models.Party, error)
}

func (m *mockGuestPartyDAO) GetByID(ctx context.Context, id string) (*models.Party, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, persistence.ErrNotFound
}

func (m *mockGuestPartyDAO) GetByCode(ctx context.Context, code string) (*models.Party, error) {
	if m.getByCodeFunc != nil {
		return m.getByCodeFunc(ctx, code)
	}
	return nil, persistence.ErrNotFound
}

func TestGuestService_JoinParty(t *testing.T) {
	t.Run("joins party successfully", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		var createdGuest *models.Guest
		guestDAO := &mockGuestDAO{
			existsByPartyAndUsernameFunc: func(ctx context.Context, partyID, username string) (bool, error) {
				return false, nil
			},
			createFunc: func(ctx context.Context, guest *models.Guest) error {
				createdGuest = guest
				return nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				if code == "ABC123" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guest, err := svc.JoinParty(ctx, "ABC123", "alice")

		require.NoError(t, err)
		assert.NotEmpty(t, guest.ID)
		assert.Equal(t, "party-1", guest.PartyID)
		assert.Equal(t, "alice", guest.Username)
		assert.Equal(t, models.GuestStatusPending, guest.Status)
		assert.False(t, guest.CreatedAt.IsZero())
		assert.Equal(t, createdGuest, guest)
	})

	t.Run("returns ErrNotFound when party code not found", func(t *testing.T) {
		guestDAO := &mockGuestDAO{}
		partyDAO := &mockGuestPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guest, err := svc.JoinParty(ctx, "NONEXISTENT", "alice")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, guest)
	})

	t.Run("returns ErrDuplicateUsername when username exists in party", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		guestDAO := &mockGuestDAO{
			existsByPartyAndUsernameFunc: func(ctx context.Context, partyID, username string) (bool, error) {
				return true, nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				if code == "ABC123" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guest, err := svc.JoinParty(ctx, "ABC123", "alice")

		assert.ErrorIs(t, err, services.ErrDuplicateUsername)
		assert.Nil(t, guest)
	})
}

func TestGuestService_ListGuests(t *testing.T) {
	t.Run("returns approved guests for admin", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		approvedGuests := []*models.Guest{
			{
				ID:        "guest-1",
				PartyID:   "party-1",
				Username:  "alice",
				Status:    models.GuestStatusApproved,
				CreatedAt: time.Now(),
			},
			{
				ID:        "guest-2",
				PartyID:   "party-1",
				Username:  "bob",
				Status:    models.GuestStatusApproved,
				CreatedAt: time.Now(),
			},
		}

		guestDAO := &mockGuestDAO{
			listByPartyIDAndStatusFunc: func(ctx context.Context, partyID string, status models.GuestStatus) ([]*models.Guest, error) {
				if partyID == "party-1" && status == models.GuestStatusApproved {
					return approvedGuests, nil
				}
				return []*models.Guest{}, nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guests, err := svc.ListGuests(ctx, "admin-1", "party-1")

		require.NoError(t, err)
		assert.Len(t, guests, 2)
		assert.Equal(t, approvedGuests, guests)
	})

	t.Run("returns ErrUnauthorized for non-owner", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		guestDAO := &mockGuestDAO{}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guests, err := svc.ListGuests(ctx, "other-admin", "party-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, guests)
	})

	t.Run("returns ErrNotFound when party not found", func(t *testing.T) {
		guestDAO := &mockGuestDAO{}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guests, err := svc.ListGuests(ctx, "admin-1", "nonexistent")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, guests)
	})
}

func TestGuestService_ListGuestsAsGuest(t *testing.T) {
	t.Run("returns approved guests for approved guest", func(t *testing.T) {
		existingGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		approvedGuests := []*models.Guest{
			existingGuest,
			{
				ID:        "guest-2",
				PartyID:   "party-1",
				Username:  "bob",
				Status:    models.GuestStatusApproved,
				CreatedAt: time.Now(),
			},
		}

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return existingGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
			listByPartyIDAndStatusFunc: func(ctx context.Context, partyID string, status models.GuestStatus) ([]*models.Guest, error) {
				if partyID == "party-1" && status == models.GuestStatusApproved {
					return approvedGuests, nil
				}
				return []*models.Guest{}, nil
			},
		}
		partyDAO := &mockGuestPartyDAO{}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guests, err := svc.ListGuestsAsGuest(ctx, "guest-1", "party-1")

		require.NoError(t, err)
		assert.Len(t, guests, 2)
		assert.Equal(t, approvedGuests, guests)
	})

	t.Run("returns ErrUnauthorized when guest is not approved", func(t *testing.T) {
		existingGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusPending,
			CreatedAt: time.Now(),
		}

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return existingGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockGuestPartyDAO{}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guests, err := svc.ListGuestsAsGuest(ctx, "guest-1", "party-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, guests)
	})

	t.Run("returns ErrUnauthorized when guest belongs to different party", func(t *testing.T) {
		existingGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-2",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return existingGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockGuestPartyDAO{}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guests, err := svc.ListGuestsAsGuest(ctx, "guest-1", "party-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, guests)
	})
}

func TestGuestService_ListJoinRequests(t *testing.T) {
	t.Run("returns pending guests for admin", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		pendingGuests := []*models.Guest{
			{
				ID:        "guest-1",
				PartyID:   "party-1",
				Username:  "alice",
				Status:    models.GuestStatusPending,
				CreatedAt: time.Now(),
			},
		}

		guestDAO := &mockGuestDAO{
			listByPartyIDAndStatusFunc: func(ctx context.Context, partyID string, status models.GuestStatus) ([]*models.Guest, error) {
				if partyID == "party-1" && status == models.GuestStatusPending {
					return pendingGuests, nil
				}
				return []*models.Guest{}, nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guests, err := svc.ListJoinRequests(ctx, "admin-1", "party-1")

		require.NoError(t, err)
		assert.Len(t, guests, 1)
		assert.Equal(t, pendingGuests, guests)
	})

	t.Run("returns ErrUnauthorized for non-owner", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		guestDAO := &mockGuestDAO{}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guests, err := svc.ListJoinRequests(ctx, "other-admin", "party-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, guests)
	})
}

func TestGuestService_ApproveGuest(t *testing.T) {
	t.Run("approves pending guest", func(t *testing.T) {
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

		updateStatusCalled := false
		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return pendingGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
			updateStatusFunc: func(ctx context.Context, id string, status models.GuestStatus) error {
				updateStatusCalled = true
				assert.Equal(t, "guest-1", id)
				assert.Equal(t, models.GuestStatusApproved, status)
				return nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.ApproveGuest(ctx, "admin-1", "party-1", "guest-1")

		require.NoError(t, err)
		assert.True(t, updateStatusCalled)
	})

	t.Run("returns ErrUnauthorized for non-owner", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		updateStatusCalled := false
		guestDAO := &mockGuestDAO{
			updateStatusFunc: func(ctx context.Context, id string, status models.GuestStatus) error {
				updateStatusCalled = true
				return nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.ApproveGuest(ctx, "other-admin", "party-1", "guest-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.False(t, updateStatusCalled)
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

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.ApproveGuest(ctx, "admin-1", "party-1", "nonexistent")

		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("returns ErrNotFound when guest belongs to different party", func(t *testing.T) {
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
			Status:    models.GuestStatusPending,
			CreatedAt: time.Now(),
		}

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return guestFromOtherParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.ApproveGuest(ctx, "admin-1", "party-1", "guest-1")

		assert.ErrorIs(t, err, services.ErrNotFound)
	})
}

func TestGuestService_RejectGuest(t *testing.T) {
	t.Run("rejects pending guest", func(t *testing.T) {
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

		updateStatusCalled := false
		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return pendingGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
			updateStatusFunc: func(ctx context.Context, id string, status models.GuestStatus) error {
				updateStatusCalled = true
				assert.Equal(t, "guest-1", id)
				assert.Equal(t, models.GuestStatusRejected, status)
				return nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.RejectGuest(ctx, "admin-1", "party-1", "guest-1")

		require.NoError(t, err)
		assert.True(t, updateStatusCalled)
	})

	t.Run("returns ErrUnauthorized for non-owner", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		updateStatusCalled := false
		guestDAO := &mockGuestDAO{
			updateStatusFunc: func(ctx context.Context, id string, status models.GuestStatus) error {
				updateStatusCalled = true
				return nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.RejectGuest(ctx, "other-admin", "party-1", "guest-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.False(t, updateStatusCalled)
	})
}

func TestGuestService_RemoveGuest(t *testing.T) {
	t.Run("removes guest", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		existingGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusApproved,
			CreatedAt: time.Now(),
		}

		deleteCalled := false
		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return existingGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
			deleteFunc: func(ctx context.Context, id string) error {
				deleteCalled = true
				assert.Equal(t, "guest-1", id)
				return nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.RemoveGuest(ctx, "admin-1", "party-1", "guest-1")

		require.NoError(t, err)
		assert.True(t, deleteCalled)
	})

	t.Run("returns ErrUnauthorized for non-owner", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		deleteCalled := false
		guestDAO := &mockGuestDAO{
			deleteFunc: func(ctx context.Context, id string) error {
				deleteCalled = true
				return nil
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.RemoveGuest(ctx, "other-admin", "party-1", "guest-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.False(t, deleteCalled)
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

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		err := svc.RemoveGuest(ctx, "admin-1", "party-1", "nonexistent")

		assert.ErrorIs(t, err, services.ErrNotFound)
	})
}

func TestGuestService_GetGuestStatus(t *testing.T) {
	t.Run("returns guest status", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		existingGuest := &models.Guest{
			ID:        "guest-1",
			PartyID:   "party-1",
			Username:  "alice",
			Status:    models.GuestStatusPending,
			CreatedAt: time.Now(),
		}

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return existingGuest, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				if code == "ABC123" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guest, err := svc.GetGuestStatus(ctx, "ABC123", "guest-1")

		require.NoError(t, err)
		assert.Equal(t, existingGuest, guest)
	})

	t.Run("returns ErrNotFound when party code not found", func(t *testing.T) {
		guestDAO := &mockGuestDAO{}
		partyDAO := &mockGuestPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guest, err := svc.GetGuestStatus(ctx, "NONEXISTENT", "guest-1")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, guest)
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

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				if code == "ABC123" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guest, err := svc.GetGuestStatus(ctx, "ABC123", "nonexistent")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, guest)
	})

	t.Run("returns ErrNotFound when guest belongs to different party", func(t *testing.T) {
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

		guestDAO := &mockGuestDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Guest, error) {
				if id == "guest-1" {
					return guestFromOtherParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}
		partyDAO := &mockGuestPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				if code == "ABC123" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewGuestService(guestDAO, partyDAO)
		ctx := context.Background()

		guest, err := svc.GetGuestStatus(ctx, "ABC123", "guest-1")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, guest)
	})
}
