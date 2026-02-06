package persistence_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
)

func TestFirestoreUserDAO_Upsert(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "users") })

	dao := persistence.NewFirestoreUserDAO(client)
	ctx := context.Background()

	t.Run("creates and retrieves user", func(t *testing.T) {
		user := &models.User{
			ID:       "user-1",
			Username: "testuser",
			Email:    "test@example.com",
		}

		err := dao.Upsert(ctx, user)

		require.NoError(t, err)

		retrieved, err := dao.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Username, retrieved.Username)
		assert.Equal(t, user.Email, retrieved.Email)
	})

	t.Run("overwrites existing user", func(t *testing.T) {
		user := &models.User{
			ID:       "user-2",
			Username: "original",
			Email:    "original@example.com",
		}

		err := dao.Upsert(ctx, user)
		require.NoError(t, err)

		user.Username = "updated"
		err = dao.Upsert(ctx, user)
		require.NoError(t, err)

		retrieved, err := dao.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated", retrieved.Username)
	})
}

func TestFirestoreUserDAO_GetByID(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "users") })

	dao := persistence.NewFirestoreUserDAO(client)
	ctx := context.Background()

	t.Run("returns ErrNotFound for missing user", func(t *testing.T) {
		_, err := dao.GetByID(ctx, "nonexistent-id")

		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})
}
