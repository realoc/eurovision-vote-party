package persistence_test

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
)

func setupFirestoreClient(t *testing.T) *firestore.Client {
	t.Helper()

	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST not set, skipping Firestore tests")
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		t.Fatalf("failed to create Firestore client: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	return client
}

func cleanupCollection(t *testing.T, client *firestore.Client, collection string) {
	t.Helper()
	ctx := context.Background()
	docs, err := client.Collection(collection).Documents(ctx).GetAll()
	if err != nil {
		return
	}
	for _, doc := range docs {
		doc.Ref.Delete(ctx)
	}
}

func createTestParty(id, code, adminID string) *models.Party {
	return &models.Party{
		ID:        id,
		Name:      "Test Party",
		Code:      code,
		EventType: models.EventGrandFinal,
		AdminID:   adminID,
		Status:    models.PartyStatusActive,
		CreatedAt: time.Now(),
	}
}

func TestFirestorePartyDAO_Create(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "parties") })

	dao := persistence.NewFirestorePartyDAO(client)
	ctx := context.Background()

	t.Run("creates party successfully", func(t *testing.T) {
		party := createTestParty("party-1", "CODE1", "admin-1")

		err := dao.Create(ctx, party)

		require.NoError(t, err)

		retrieved, err := dao.GetByID(ctx, party.ID)
		require.NoError(t, err)
		assert.Equal(t, party.ID, retrieved.ID)
		assert.Equal(t, party.Code, retrieved.Code)
		assert.Equal(t, party.Name, retrieved.Name)
	})

	t.Run("returns ErrCodeExists when code already exists", func(t *testing.T) {
		party1 := createTestParty("party-2", "DUPCODE", "admin-1")
		party2 := createTestParty("party-3", "DUPCODE", "admin-2")

		err := dao.Create(ctx, party1)
		require.NoError(t, err)

		err = dao.Create(ctx, party2)

		assert.ErrorIs(t, err, persistence.ErrCodeExists)
	})
}

func TestFirestorePartyDAO_GetByID(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "parties") })

	dao := persistence.NewFirestorePartyDAO(client)
	ctx := context.Background()

	t.Run("returns party when exists", func(t *testing.T) {
		party := createTestParty("party-get-1", "GETCODE1", "admin-1")
		err := dao.Create(ctx, party)
		require.NoError(t, err)

		retrieved, err := dao.GetByID(ctx, party.ID)

		require.NoError(t, err)
		assert.Equal(t, party.ID, retrieved.ID)
		assert.Equal(t, party.Code, retrieved.Code)
		assert.Equal(t, party.AdminID, retrieved.AdminID)
	})

	t.Run("returns ErrNotFound when party does not exist", func(t *testing.T) {
		_, err := dao.GetByID(ctx, "nonexistent-id")

		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})
}

func TestFirestorePartyDAO_GetByCode(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "parties") })

	dao := persistence.NewFirestorePartyDAO(client)
	ctx := context.Background()

	t.Run("returns party when code exists", func(t *testing.T) {
		party := createTestParty("party-code-1", "BYCODE1", "admin-1")
		err := dao.Create(ctx, party)
		require.NoError(t, err)

		retrieved, err := dao.GetByCode(ctx, party.Code)

		require.NoError(t, err)
		assert.Equal(t, party.ID, retrieved.ID)
		assert.Equal(t, party.Code, retrieved.Code)
	})

	t.Run("returns ErrNotFound when code does not exist", func(t *testing.T) {
		_, err := dao.GetByCode(ctx, "NONEXISTENT")

		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})
}

func TestFirestorePartyDAO_ListByAdminID(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "parties") })

	dao := persistence.NewFirestorePartyDAO(client)
	ctx := context.Background()

	t.Run("returns parties for admin", func(t *testing.T) {
		adminID := "admin-list-1"
		party1 := createTestParty("party-list-1", "LIST1", adminID)
		party2 := createTestParty("party-list-2", "LIST2", adminID)
		party3 := createTestParty("party-list-3", "LIST3", "other-admin")

		require.NoError(t, dao.Create(ctx, party1))
		require.NoError(t, dao.Create(ctx, party2))
		require.NoError(t, dao.Create(ctx, party3))

		parties, err := dao.ListByAdminID(ctx, adminID)

		require.NoError(t, err)
		assert.Len(t, parties, 2)

		ids := make([]string, len(parties))
		for i, p := range parties {
			ids[i] = p.ID
		}
		assert.Contains(t, ids, party1.ID)
		assert.Contains(t, ids, party2.ID)
	})

	t.Run("returns empty slice when no parties for admin", func(t *testing.T) {
		parties, err := dao.ListByAdminID(ctx, "admin-with-no-parties")

		require.NoError(t, err)
		assert.Empty(t, parties)
		assert.NotNil(t, parties)
	})
}

func TestFirestorePartyDAO_Delete(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "parties") })

	dao := persistence.NewFirestorePartyDAO(client)
	ctx := context.Background()

	t.Run("removes party successfully", func(t *testing.T) {
		party := createTestParty("party-del-1", "DELCODE1", "admin-1")
		err := dao.Create(ctx, party)
		require.NoError(t, err)

		err = dao.Delete(ctx, party.ID)

		require.NoError(t, err)

		_, err = dao.GetByID(ctx, party.ID)
		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})

	t.Run("returns ErrNotFound when party does not exist", func(t *testing.T) {
		err := dao.Delete(ctx, "nonexistent-party")

		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})
}

func TestFirestorePartyDAO_CodeExists(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "parties") })

	dao := persistence.NewFirestorePartyDAO(client)
	ctx := context.Background()

	t.Run("returns true when code exists", func(t *testing.T) {
		party := createTestParty("party-exists-1", "EXISTS1", "admin-1")
		err := dao.Create(ctx, party)
		require.NoError(t, err)

		exists, err := dao.CodeExists(ctx, party.Code)

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false when code does not exist", func(t *testing.T) {
		exists, err := dao.CodeExists(ctx, "NOTEXISTS")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}
