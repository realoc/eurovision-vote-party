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

func createTestVote(id, guestID, partyID string) *models.Vote {
	return &models.Vote{
		ID:      id,
		GuestID: guestID,
		PartyID: partyID,
		Votes: map[int]string{
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
		},
		CreatedAt: time.Now(),
	}
}

func TestFirestoreVoteDAO_Create(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "votes") })

	dao := persistence.NewFirestoreVoteDAO(client)
	ctx := context.Background()

	t.Run("creates vote successfully", func(t *testing.T) {
		vote := createTestVote("vote-create-1", "guest-1", "party-1")

		err := dao.Create(ctx, vote)

		require.NoError(t, err)

		retrieved, err := dao.GetByGuestAndParty(ctx, vote.GuestID, vote.PartyID)
		require.NoError(t, err)
		assert.Equal(t, vote.ID, retrieved.ID)
		assert.Equal(t, vote.GuestID, retrieved.GuestID)
		assert.Equal(t, vote.PartyID, retrieved.PartyID)
		assert.Equal(t, vote.Votes, retrieved.Votes)
	})
}

func TestFirestoreVoteDAO_GetByGuestAndParty(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "votes") })

	dao := persistence.NewFirestoreVoteDAO(client)
	ctx := context.Background()

	t.Run("returns vote when exists", func(t *testing.T) {
		vote := createTestVote("vote-get-1", "guest-get-1", "party-get-1")
		err := dao.Create(ctx, vote)
		require.NoError(t, err)

		retrieved, err := dao.GetByGuestAndParty(ctx, vote.GuestID, vote.PartyID)

		require.NoError(t, err)
		assert.Equal(t, vote.ID, retrieved.ID)
		assert.Equal(t, vote.GuestID, retrieved.GuestID)
		assert.Equal(t, vote.PartyID, retrieved.PartyID)
		assert.Equal(t, vote.Votes, retrieved.Votes)
	})

	t.Run("returns ErrNotFound when vote does not exist", func(t *testing.T) {
		_, err := dao.GetByGuestAndParty(ctx, "nonexistent-guest", "nonexistent-party")

		assert.ErrorIs(t, err, persistence.ErrNotFound)
	})
}

func TestFirestoreVoteDAO_Update(t *testing.T) {
	client := setupFirestoreClient(t)
	t.Cleanup(func() { cleanupCollection(t, client, "votes") })

	dao := persistence.NewFirestoreVoteDAO(client)
	ctx := context.Background()

	t.Run("updates vote successfully", func(t *testing.T) {
		vote := createTestVote("vote-update-1", "guest-update-1", "party-update-1")
		err := dao.Create(ctx, vote)
		require.NoError(t, err)

		vote.Votes = map[int]string{
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

		err = dao.Update(ctx, vote)

		require.NoError(t, err)

		retrieved, err := dao.GetByGuestAndParty(ctx, vote.GuestID, vote.PartyID)
		require.NoError(t, err)
		assert.Equal(t, vote.ID, retrieved.ID)
		assert.Equal(t, "act-10", retrieved.Votes[12])
		assert.Equal(t, "act-9", retrieved.Votes[10])
		assert.Equal(t, "act-1", retrieved.Votes[1])
	})
}
