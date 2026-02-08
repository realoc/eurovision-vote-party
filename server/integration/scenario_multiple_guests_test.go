//go:build integration

package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
)

func TestMultipleGuestsVoting(t *testing.T) {
	env := setupTest(t)
	ctx := context.Background()
	adminID := "admin-multi-guest"

	// Create party
	party := mustCreateParty(t, env, adminID, "Multi Guest Party")
	acts := mustGetGrandFinalActs(t, env)

	// Join and approve 3 guests
	guest1 := mustJoinParty(t, env, party.Code, "Alice")
	guest2 := mustJoinParty(t, env, party.Code, "Bob")
	guest3 := mustJoinParty(t, env, party.Code, "Charlie")

	mustApproveGuest(t, env, adminID, party.ID, guest1.ID)
	mustApproveGuest(t, env, adminID, party.ID, guest2.ID)
	mustApproveGuest(t, env, adminID, party.ID, guest3.ID)

	// Build different vote distributions:
	// Guest 1: acts[0] gets 12 pts (offset=0)
	// Guest 2: acts[0] gets 12 pts (same as guest1 but rotated differently for other points)
	// Guest 3: acts[1] gets 12 pts (offset=1)
	votes1 := validVotesForActs(acts)           // 12->acts[0], 10->acts[1], ...
	votes2 := buildRotatedVotes(acts, 10)        // 12->acts[10], 10->acts[11], ... (different first place)
	votes2[12] = acts[0].ID                      // Override: guest2 also gives 12 to acts[0]
	votes2[10] = acts[10].ID                     // 10 pts to acts[10]
	// Ensure no duplicate act IDs in votes2
	usedActs := make(map[string]bool)
	usedActs[acts[0].ID] = true
	usedActs[acts[10].ID] = true
	nextAvail := 11
	for i, points := range models.ValidPointValues {
		if points == 12 || points == 10 {
			continue
		}
		// Find an unused act
		for usedActs[acts[nextAvail].ID] {
			nextAvail++
			if nextAvail >= len(acts) {
				nextAvail = 0
			}
		}
		_ = i
		votes2[points] = acts[nextAvail].ID
		usedActs[acts[nextAvail].ID] = true
		nextAvail++
		if nextAvail >= len(acts) {
			nextAvail = 0
		}
	}

	votes3 := buildRotatedVotes(acts, 1)         // 12->acts[1], 10->acts[2], 8->acts[3], ...

	// Submit votes
	mustSubmitVote(t, env, "", party.ID, guest1.ID, votes1)
	mustSubmitVote(t, env, "", party.ID, guest2.ID, votes2)
	mustSubmitVote(t, env, "", party.ID, guest3.ID, votes3)

	// End voting
	mustEndVoting(t, env, adminID, party.ID)

	// Get results
	results, err := env.voteService.GetResults(ctx, adminID, party.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, results.TotalVoters)

	// acts[0] should have the most points: 12 (guest1) + 12 (guest2) + 10 (guest3) = 34 (guest3 gives 10 to acts[1], not acts[0])
	// Actually let's verify the exact computation:
	// guest1: 12->acts[0]
	// guest2: 12->acts[0]
	// guest3: offset=1 means 12->acts[1], so acts[0] gets nothing from guest3
	// So acts[0] total = 12 + 12 = 24
	act0Points := totalPointsForAct(acts[0].ID, votes1, votes2, votes3)
	act1Points := totalPointsForAct(acts[1].ID, votes1, votes2, votes3)

	// Verify the result for acts[0]
	var act0Result, act1Result *models.VoteResult
	for i := range results.Results {
		if results.Results[i].ActID == acts[0].ID {
			act0Result = &results.Results[i]
		}
		if results.Results[i].ActID == acts[1].ID {
			act1Result = &results.Results[i]
		}
	}

	require.NotNil(t, act0Result, "acts[0] should appear in results")
	assert.Equal(t, act0Points, act0Result.TotalPoints)

	require.NotNil(t, act1Result, "acts[1] should appear in results")
	assert.Equal(t, act1Points, act1Result.TotalPoints)

	// Results should be sorted by total points descending
	for i := 1; i < len(results.Results); i++ {
		assert.GreaterOrEqual(t, results.Results[i-1].TotalPoints, results.Results[i].TotalPoints,
			"results should be sorted by total points descending")
	}

	// Rank 1 should be the act with the highest total points
	assert.Equal(t, 1, results.Results[0].Rank)

	// Verify ranking uses standard competition ranking (1, 2, 2, 4)
	for i := 1; i < len(results.Results); i++ {
		if results.Results[i].TotalPoints == results.Results[i-1].TotalPoints {
			assert.Equal(t, results.Results[i-1].Rank, results.Results[i].Rank,
				"tied acts should have the same rank")
		} else {
			assert.Greater(t, results.Results[i].Rank, results.Results[i-1].Rank,
				"lower-scored act should have a higher rank number")
		}
	}

	// Verify all acts from the grand final appear in results (even those with 0 points)
	allGrandFinalActs, _ := env.actsService.ListActs(string(models.EventGrandFinal))
	assert.Equal(t, len(allGrandFinalActs), len(results.Results),
		"results should include all grand final acts")
}
