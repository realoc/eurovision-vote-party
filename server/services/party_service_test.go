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

type mockPartyDAO struct {
	createFunc      func(ctx context.Context, party *models.Party) error
	getByIDFunc     func(ctx context.Context, id string) (*models.Party, error)
	getByCodeFunc   func(ctx context.Context, code string) (*models.Party, error)
	listByAdminFunc func(ctx context.Context, adminID string) ([]*models.Party, error)
	deleteFunc      func(ctx context.Context, id string) error
	codeExistsFunc  func(ctx context.Context, code string) (bool, error)
}

func (m *mockPartyDAO) Create(ctx context.Context, party *models.Party) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, party)
	}
	return nil
}

func (m *mockPartyDAO) GetByID(ctx context.Context, id string) (*models.Party, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, persistence.ErrNotFound
}

func (m *mockPartyDAO) GetByCode(ctx context.Context, code string) (*models.Party, error) {
	if m.getByCodeFunc != nil {
		return m.getByCodeFunc(ctx, code)
	}
	return nil, persistence.ErrNotFound
}

func (m *mockPartyDAO) ListByAdminID(ctx context.Context, adminID string) ([]*models.Party, error) {
	if m.listByAdminFunc != nil {
		return m.listByAdminFunc(ctx, adminID)
	}
	return []*models.Party{}, nil
}

func (m *mockPartyDAO) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockPartyDAO) CodeExists(ctx context.Context, code string) (bool, error) {
	if m.codeExistsFunc != nil {
		return m.codeExistsFunc(ctx, code)
	}
	return false, nil
}

func TestPartyService_CreateParty(t *testing.T) {
	t.Run("creates party with generated code", func(t *testing.T) {
		var createdParty *models.Party
		dao := &mockPartyDAO{
			codeExistsFunc: func(ctx context.Context, code string) (bool, error) {
				return false, nil
			},
			createFunc: func(ctx context.Context, party *models.Party) error {
				createdParty = party
				return nil
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.CreateParty(ctx, "admin-1", services.CreatePartyRequest{
			Name:      "My Party",
			EventType: models.EventGrandFinal,
		})

		require.NoError(t, err)
		assert.NotEmpty(t, party.ID)
		assert.NotEmpty(t, party.Code)
		assert.Len(t, party.Code, 6)
		assert.Equal(t, "My Party", party.Name)
		assert.Equal(t, models.EventGrandFinal, party.EventType)
		assert.Equal(t, "admin-1", party.AdminID)
		assert.Equal(t, models.PartyStatusActive, party.Status)
		assert.False(t, party.CreatedAt.IsZero())
		assert.Equal(t, createdParty, party)
	})

	t.Run("retries code generation on collision", func(t *testing.T) {
		codeExistsCalls := 0
		dao := &mockPartyDAO{
			codeExistsFunc: func(ctx context.Context, code string) (bool, error) {
				codeExistsCalls++
				// First two calls return true (collision), third returns false
				return codeExistsCalls < 3, nil
			},
			createFunc: func(ctx context.Context, party *models.Party) error {
				return nil
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.CreateParty(ctx, "admin-1", services.CreatePartyRequest{
			Name:      "My Party",
			EventType: models.EventGrandFinal,
		})

		require.NoError(t, err)
		assert.NotNil(t, party)
		assert.Equal(t, 3, codeExistsCalls)
	})

	t.Run("fails after max retries", func(t *testing.T) {
		codeExistsCalls := 0
		dao := &mockPartyDAO{
			codeExistsFunc: func(ctx context.Context, code string) (bool, error) {
				codeExistsCalls++
				return true, nil // Always collision
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.CreateParty(ctx, "admin-1", services.CreatePartyRequest{
			Name:      "My Party",
			EventType: models.EventGrandFinal,
		})

		assert.Error(t, err)
		assert.Nil(t, party)
		assert.Contains(t, err.Error(), "failed to generate unique party code")
		assert.Equal(t, 10, codeExistsCalls)
	})

	t.Run("validates party name is required", func(t *testing.T) {
		dao := &mockPartyDAO{
			codeExistsFunc: func(ctx context.Context, code string) (bool, error) {
				return false, nil
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.CreateParty(ctx, "admin-1", services.CreatePartyRequest{
			Name:      "",
			EventType: models.EventGrandFinal,
		})

		assert.Error(t, err)
		assert.Nil(t, party)
		assert.Contains(t, err.Error(), "party name is required")
	})

	t.Run("validates event type is required", func(t *testing.T) {
		dao := &mockPartyDAO{
			codeExistsFunc: func(ctx context.Context, code string) (bool, error) {
				return false, nil
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.CreateParty(ctx, "admin-1", services.CreatePartyRequest{
			Name:      "My Party",
			EventType: "",
		})

		assert.Error(t, err)
		assert.Nil(t, party)
		assert.Contains(t, err.Error(), "event type")
	})
}

func TestPartyService_GetPartyByID(t *testing.T) {
	t.Run("returns party for owner", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		dao := &mockPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.GetPartyByID(ctx, "admin-1", "party-1")

		require.NoError(t, err)
		assert.Equal(t, existingParty, party)
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

		dao := &mockPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.GetPartyByID(ctx, "other-admin", "party-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, party)
	})

	t.Run("returns ErrNotFound when party does not exist", func(t *testing.T) {
		dao := &mockPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.GetPartyByID(ctx, "admin-1", "nonexistent")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, party)
	})
}

func TestPartyService_GetPartyByCode(t *testing.T) {
	t.Run("returns party", func(t *testing.T) {
		existingParty := &models.Party{
			ID:        "party-1",
			Name:      "Test Party",
			Code:      "ABC123",
			EventType: models.EventGrandFinal,
			AdminID:   "admin-1",
			Status:    models.PartyStatusActive,
			CreatedAt: time.Now(),
		}

		dao := &mockPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				if code == "ABC123" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.GetPartyByCode(ctx, "ABC123")

		require.NoError(t, err)
		assert.Equal(t, existingParty, party)
	})

	t.Run("returns ErrNotFound when party does not exist", func(t *testing.T) {
		dao := &mockPartyDAO{
			getByCodeFunc: func(ctx context.Context, code string) (*models.Party, error) {
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		party, err := svc.GetPartyByCode(ctx, "NONEXISTENT")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, party)
	})
}

func TestPartyService_ListPartiesByAdmin(t *testing.T) {
	t.Run("returns parties", func(t *testing.T) {
		existingParties := []*models.Party{
			{
				ID:        "party-1",
				Name:      "Party 1",
				Code:      "CODE01",
				EventType: models.EventGrandFinal,
				AdminID:   "admin-1",
				Status:    models.PartyStatusActive,
				CreatedAt: time.Now(),
			},
			{
				ID:        "party-2",
				Name:      "Party 2",
				Code:      "CODE02",
				EventType: models.EventSemifinal1,
				AdminID:   "admin-1",
				Status:    models.PartyStatusActive,
				CreatedAt: time.Now(),
			},
		}

		dao := &mockPartyDAO{
			listByAdminFunc: func(ctx context.Context, adminID string) ([]*models.Party, error) {
				if adminID == "admin-1" {
					return existingParties, nil
				}
				return []*models.Party{}, nil
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		parties, err := svc.ListPartiesByAdmin(ctx, "admin-1")

		require.NoError(t, err)
		assert.Len(t, parties, 2)
		assert.Equal(t, existingParties, parties)
	})
}

func TestPartyService_DeleteParty(t *testing.T) {
	t.Run("deletes party for owner", func(t *testing.T) {
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
		dao := &mockPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				if id == "party-1" {
					return existingParty, nil
				}
				return nil, persistence.ErrNotFound
			},
			deleteFunc: func(ctx context.Context, id string) error {
				deleteCalled = true
				return nil
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		err := svc.DeleteParty(ctx, "admin-1", "party-1")

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
		dao := &mockPartyDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.Party, error) {
				return existingParty, nil
			},
			deleteFunc: func(ctx context.Context, id string) error {
				deleteCalled = true
				return nil
			},
		}

		svc := services.NewPartyService(dao)
		ctx := context.Background()

		err := svc.DeleteParty(ctx, "other-admin", "party-1")

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.False(t, deleteCalled)
	})
}
