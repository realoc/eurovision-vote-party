package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
)

func createTestGuest(id, partyID, username string, status models.GuestStatus) *models.Guest {
	return &models.Guest{
		ID:        id,
		PartyID:   partyID,
		Username:  username,
		Status:    status,
		CreatedAt: time.Now(),
	}
}

func TestFirestoreGuestDAO_Create(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "guests") })

	dao := persistence.NewFirestoreGuestDAO(client)
	ctx := context.Background()

	t.Run("creates guest successfully", func(t *testing.T) {
		guest := createTestGuest("guest-create-1", "party-1", "alice", models.GuestStatusPending)

		err := dao.Create(ctx, guest)

		require.NoError(t, err)

		retrieved, err := dao.GetByID(ctx, guest.ID)
		require.NoError(t, err)
		assert.Equal(t, guest.ID, retrieved.ID)
		assert.Equal(t, guest.PartyID, retrieved.PartyID)
		assert.Equal(t, guest.Username, retrieved.Username)
		assert.Equal(t, guest.Status, retrieved.Status)
	})
}

func TestFirestoreGuestDAO_GetByID(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "guests") })

	dao := persistence.NewFirestoreGuestDAO(client)
	ctx := context.Background()

	t.Run("returns guest when exists", func(t *testing.T) {
		guest := createTestGuest("guest-get-1", "party-1", "bob", models.GuestStatusApproved)
		err := dao.Create(ctx, guest)
		require.NoError(t, err)

		retrieved, err := dao.GetByID(ctx, guest.ID)

		require.NoError(t, err)
		assert.Equal(t, guest.ID, retrieved.ID)
		assert.Equal(t, guest.PartyID, retrieved.PartyID)
		assert.Equal(t, guest.Username, retrieved.Username)
		assert.Equal(t, guest.Status, retrieved.Status)
	})

	t.Run("returns ErrNotFound when guest does not exist", func(t *testing.T) {
		_, err := dao.GetByID(ctx, "nonexistent-id")

		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})
}

func TestFirestoreGuestDAO_ListByPartyID(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "guests") })

	dao := persistence.NewFirestoreGuestDAO(client)
	ctx := context.Background()

	t.Run("returns all guests for party", func(t *testing.T) {
		partyID := "party-list-1"
		guest1 := createTestGuest("guest-list-1", partyID, "charlie", models.GuestStatusPending)
		guest2 := createTestGuest("guest-list-2", partyID, "dana", models.GuestStatusApproved)
		guest3 := createTestGuest("guest-list-3", "other-party", "eve", models.GuestStatusPending)

		require.NoError(t, dao.Create(ctx, guest1))
		require.NoError(t, dao.Create(ctx, guest2))
		require.NoError(t, dao.Create(ctx, guest3))

		guests, err := dao.ListByPartyID(ctx, partyID)

		require.NoError(t, err)
		assert.Len(t, guests, 2)

		ids := make([]string, len(guests))
		for i, g := range guests {
			ids[i] = g.ID
		}
		assert.Contains(t, ids, guest1.ID)
		assert.Contains(t, ids, guest2.ID)
	})

	t.Run("returns empty slice when no guests for party", func(t *testing.T) {
		guests, err := dao.ListByPartyID(ctx, "party-with-no-guests")

		require.NoError(t, err)
		assert.Empty(t, guests)
		assert.NotNil(t, guests)
	})
}

func TestFirestoreGuestDAO_ListByPartyIDAndStatus(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "guests") })

	dao := persistence.NewFirestoreGuestDAO(client)
	ctx := context.Background()

	t.Run("filters by status correctly", func(t *testing.T) {
		partyID := "party-status-1"
		guest1 := createTestGuest("guest-status-1", partyID, "frank", models.GuestStatusPending)
		guest2 := createTestGuest("guest-status-2", partyID, "grace", models.GuestStatusApproved)
		guest3 := createTestGuest("guest-status-3", partyID, "heidi", models.GuestStatusPending)
		guest4 := createTestGuest("guest-status-4", partyID, "ivan", models.GuestStatusRejected)

		require.NoError(t, dao.Create(ctx, guest1))
		require.NoError(t, dao.Create(ctx, guest2))
		require.NoError(t, dao.Create(ctx, guest3))
		require.NoError(t, dao.Create(ctx, guest4))

		pending, err := dao.ListByPartyIDAndStatus(ctx, partyID, models.GuestStatusPending)
		require.NoError(t, err)
		assert.Len(t, pending, 2)

		ids := make([]string, len(pending))
		for i, g := range pending {
			ids[i] = g.ID
		}
		assert.Contains(t, ids, guest1.ID)
		assert.Contains(t, ids, guest3.ID)

		approved, err := dao.ListByPartyIDAndStatus(ctx, partyID, models.GuestStatusApproved)
		require.NoError(t, err)
		assert.Len(t, approved, 1)
		assert.Equal(t, guest2.ID, approved[0].ID)

		rejected, err := dao.ListByPartyIDAndStatus(ctx, partyID, models.GuestStatusRejected)
		require.NoError(t, err)
		assert.Len(t, rejected, 1)
		assert.Equal(t, guest4.ID, rejected[0].ID)
	})
}

func TestFirestoreGuestDAO_UpdateStatus(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "guests") })

	dao := persistence.NewFirestoreGuestDAO(client)
	ctx := context.Background()

	t.Run("updates status successfully", func(t *testing.T) {
		guest := createTestGuest("guest-update-1", "party-1", "judy", models.GuestStatusPending)
		err := dao.Create(ctx, guest)
		require.NoError(t, err)

		err = dao.UpdateStatus(ctx, guest.ID, models.GuestStatusApproved)

		require.NoError(t, err)

		retrieved, err := dao.GetByID(ctx, guest.ID)
		require.NoError(t, err)
		assert.Equal(t, models.GuestStatusApproved, retrieved.Status)
	})

	t.Run("returns ErrNotFound when guest does not exist", func(t *testing.T) {
		err := dao.UpdateStatus(ctx, "nonexistent-guest", models.GuestStatusApproved)

		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})
}

func TestFirestoreGuestDAO_Delete(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "guests") })

	dao := persistence.NewFirestoreGuestDAO(client)
	ctx := context.Background()

	t.Run("removes guest successfully", func(t *testing.T) {
		guest := createTestGuest("guest-del-1", "party-1", "karl", models.GuestStatusPending)
		err := dao.Create(ctx, guest)
		require.NoError(t, err)

		err = dao.Delete(ctx, guest.ID)

		require.NoError(t, err)

		_, err = dao.GetByID(ctx, guest.ID)
		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})

	t.Run("returns ErrNotFound when guest does not exist", func(t *testing.T) {
		err := dao.Delete(ctx, "nonexistent-guest")

		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})
}

func TestFirestoreGuestDAO_ExistsByPartyAndUsername(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "guests") })

	dao := persistence.NewFirestoreGuestDAO(client)
	ctx := context.Background()

	t.Run("returns true when guest exists for party and username", func(t *testing.T) {
		guest := createTestGuest("guest-exists-1", "party-exists-1", "laura", models.GuestStatusPending)
		err := dao.Create(ctx, guest)
		require.NoError(t, err)

		exists, err := dao.ExistsByPartyAndUsername(ctx, guest.PartyID, guest.Username)

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false when guest does not exist", func(t *testing.T) {
		exists, err := dao.ExistsByPartyAndUsername(ctx, "party-exists-1", "nonexistent-user")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}
