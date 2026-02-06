package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

type mockUserDAO struct {
	upsertFunc  func(ctx context.Context, user *models.User) error
	getByIDFunc func(ctx context.Context, id string) (*models.User, error)
}

func (m *mockUserDAO) Upsert(ctx context.Context, user *models.User) error {
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, user)
	}
	return nil
}

func (m *mockUserDAO) GetByID(ctx context.Context, id string) (*models.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, persistence.ErrNotFound
}

func TestUserService_UpsertProfile(t *testing.T) {
	t.Run("successfully upserts profile", func(t *testing.T) {
		var capturedUser *models.User
		dao := &mockUserDAO{
			upsertFunc: func(ctx context.Context, user *models.User) error {
				capturedUser = user
				return nil
			},
		}

		svc := services.NewUserService(dao)
		ctx := context.Background()

		user, err := svc.UpsertProfile(ctx, "user-1", "user@example.com", "validuser")

		require.NoError(t, err)
		assert.Equal(t, "user-1", user.ID)
		assert.Equal(t, "user@example.com", user.Email)
		assert.Equal(t, "validuser", user.Username)
		assert.Equal(t, capturedUser, user)
	})

	t.Run("returns ErrInvalidUsername for invalid username", func(t *testing.T) {
		dao := &mockUserDAO{}

		svc := services.NewUserService(dao)
		ctx := context.Background()

		user, err := svc.UpsertProfile(ctx, "user-1", "user@example.com", "ab")

		assert.ErrorIs(t, err, services.ErrInvalidUsername)
		assert.Nil(t, user)
	})

	t.Run("propagates DAO error", func(t *testing.T) {
		daoErr := errors.New("database error")
		dao := &mockUserDAO{
			upsertFunc: func(ctx context.Context, user *models.User) error {
				return daoErr
			},
		}

		svc := services.NewUserService(dao)
		ctx := context.Background()

		user, err := svc.UpsertProfile(ctx, "user-1", "user@example.com", "validuser")

		assert.ErrorIs(t, err, daoErr)
		assert.Nil(t, user)
	})
}

func TestUserService_GetProfile(t *testing.T) {
	t.Run("returns user profile", func(t *testing.T) {
		existingUser := &models.User{
			ID:       "user-1",
			Email:    "user@example.com",
			Username: "validuser",
		}

		dao := &mockUserDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				if id == "user-1" {
					return existingUser, nil
				}
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewUserService(dao)
		ctx := context.Background()

		user, err := svc.GetProfile(ctx, "user-1")

		require.NoError(t, err)
		assert.Equal(t, existingUser, user)
	})

	t.Run("maps ErrNotFound", func(t *testing.T) {
		dao := &mockUserDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return nil, persistence.ErrNotFound
			},
		}

		svc := services.NewUserService(dao)
		ctx := context.Background()

		user, err := svc.GetProfile(ctx, "nonexistent")

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, user)
	})

	t.Run("propagates other errors", func(t *testing.T) {
		daoErr := errors.New("database error")
		dao := &mockUserDAO{
			getByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return nil, daoErr
			},
		}

		svc := services.NewUserService(dao)
		ctx := context.Background()

		user, err := svc.GetProfile(ctx, "user-1")

		assert.ErrorIs(t, err, daoErr)
		assert.Nil(t, user)
	})
}
