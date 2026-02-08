//go:build integration

package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
)

func TestCompletePartyFlow(t *testing.T) {
	env := setupTest(t)
	ctx := context.Background()
	adminID := "admin-complete-flow"

	// Step 1: Admin creates a profile
	user, err := env.userService.UpsertProfile(ctx, adminID, "admin@example.com", "AdminUser")
	require.NoError(t, err)
	assert.Equal(t, adminID, user.ID)
	assert.Equal(t, "AdminUser", user.Username)

	// Step 2: Admin creates a party
	party := mustCreateParty(t, env, adminID, "Complete Flow Party")
	assert.Equal(t, models.PartyStatusActive, party.Status)
	assert.Equal(t, models.EventGrandFinal, party.EventType)
	assert.NotEmpty(t, party.Code)
	assert.NotEmpty(t, party.ID)

	// Step 3: Guest joins the party (pending)
	guest := mustJoinParty(t, env, party.Code, "GuestOne")
	assert.Equal(t, models.GuestStatusPending, guest.Status)
	assert.Equal(t, party.ID, guest.PartyID)

	// Step 4: Admin sees join requests
	joinRequests, err := env.guestService.ListJoinRequests(ctx, adminID, party.ID)
	require.NoError(t, err)
	require.Len(t, joinRequests, 1)
	assert.Equal(t, guest.ID, joinRequests[0].ID)
	assert.Equal(t, models.GuestStatusPending, joinRequests[0].Status)

	// Step 5: Admin approves the guest
	mustApproveGuest(t, env, adminID, party.ID, guest.ID)

	// Verify guest is now approved
	guestStatus, err := env.guestService.GetGuestStatus(ctx, party.Code, guest.ID)
	require.NoError(t, err)
	assert.Equal(t, models.GuestStatusApproved, guestStatus.Status)

	// Verify approved guest appears in guest list
	guests, err := env.guestService.ListGuests(ctx, adminID, party.ID)
	require.NoError(t, err)
	require.Len(t, guests, 1)
	assert.Equal(t, guest.ID, guests[0].ID)

	// Step 6: Guest submits a vote
	acts := mustGetGrandFinalActs(t, env)
	votes := validVotesForActs(acts)
	vote := mustSubmitVote(t, env, "", party.ID, guest.ID, votes)
	assert.Equal(t, guest.ID, vote.GuestID)
	assert.Equal(t, party.ID, vote.PartyID)
	assert.Equal(t, votes, vote.Votes)

	// Step 7: Admin ends voting
	closedParty := mustEndVoting(t, env, adminID, party.ID)
	assert.Equal(t, models.PartyStatusClosed, closedParty.Status)

	// Step 8: Verify results
	results, err := env.voteService.GetResults(ctx, adminID, party.ID)
	require.NoError(t, err)
	assert.Equal(t, party.ID, results.PartyID)
	assert.Equal(t, 1, results.TotalVoters)

	// The act that received 12 points should be ranked first
	topActID := actIDAtPoints(votes, 12)
	require.NotEmpty(t, results.Results)
	assert.Equal(t, topActID, results.Results[0].ActID)
	assert.Equal(t, 12, results.Results[0].TotalPoints)
	assert.Equal(t, 1, results.Results[0].Rank)

	// The act that received 1 point should be ranked last among voted acts
	lastActID := actIDAtPoints(votes, 1)
	// Find the result for the act with 1 point
	for _, r := range results.Results {
		if r.ActID == lastActID {
			assert.Equal(t, 1, r.TotalPoints)
			break
		}
	}
}
