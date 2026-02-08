//go:build integration

package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/services"
)

func TestEdgeCases(t *testing.T) {
	t.Run("rejected guest cannot vote", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-rejected"

		party := mustCreateParty(t, env, adminID, "Rejected Guest Party")
		acts := mustGetGrandFinalActs(t, env)
		guest := mustJoinParty(t, env, party.Code, "RejectedGuest")

		// Reject the guest
		err := env.guestService.RejectGuest(ctx, adminID, party.ID, guest.ID)
		require.NoError(t, err)

		// Attempt to vote
		votes := validVotesForActs(acts)
		_, err = env.voteService.SubmitVote(ctx, "", party.ID, services.SubmitVoteRequest{
			GuestID: guest.ID,
			Votes:   votes,
		})
		assert.ErrorIs(t, err, services.ErrGuestNotApproved)
	})

	t.Run("pending guest cannot vote", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-pending"

		party := mustCreateParty(t, env, adminID, "Pending Guest Party")
		acts := mustGetGrandFinalActs(t, env)
		guest := mustJoinParty(t, env, party.Code, "PendingGuest")

		// Guest is still pending - attempt to vote
		votes := validVotesForActs(acts)
		_, err := env.voteService.SubmitVote(ctx, "", party.ID, services.SubmitVoteRequest{
			GuestID: guest.ID,
			Votes:   votes,
		})
		assert.ErrorIs(t, err, services.ErrGuestNotApproved)
	})

	t.Run("duplicate username rejected", func(t *testing.T) {
		env := setupTest(t)
		adminID := "admin-dup"

		party := mustCreateParty(t, env, adminID, "Dup Username Party")
		mustJoinParty(t, env, party.Code, "SameName")

		// Try to join with the same username
		_, err := env.guestService.JoinParty(context.Background(), party.Code, "SameName")
		assert.ErrorIs(t, err, services.ErrDuplicateUsername)
	})

	t.Run("invalid act ID rejected", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-invalid-act"

		party := mustCreateParty(t, env, adminID, "Invalid Act Party")
		acts := mustGetGrandFinalActs(t, env)
		guest := mustJoinParty(t, env, party.Code, "ValidGuest")
		mustApproveGuest(t, env, adminID, party.ID, guest.ID)

		// Build votes with one invalid act ID
		votes := validVotesForActs(acts)
		votes[12] = "nonexistent-act-id"

		_, err := env.voteService.SubmitVote(ctx, "", party.ID, services.SubmitVoteRequest{
			GuestID: guest.ID,
			Votes:   votes,
		})
		assert.ErrorIs(t, err, services.ErrInvalidVotes)
	})

	t.Run("duplicate vote rejected", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-dup-vote"

		party := mustCreateParty(t, env, adminID, "Dup Vote Party")
		acts := mustGetGrandFinalActs(t, env)
		guest := mustJoinParty(t, env, party.Code, "VoterDup")
		mustApproveGuest(t, env, adminID, party.ID, guest.ID)

		votes := validVotesForActs(acts)
		mustSubmitVote(t, env, "", party.ID, guest.ID, votes)

		// Try to submit again
		_, err := env.voteService.SubmitVote(ctx, "", party.ID, services.SubmitVoteRequest{
			GuestID: guest.ID,
			Votes:   votes,
		})
		assert.ErrorIs(t, err, services.ErrVoteAlreadyExists)
	})

	t.Run("vote after close rejected", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-closed"

		party := mustCreateParty(t, env, adminID, "Closed Party")
		acts := mustGetGrandFinalActs(t, env)
		guest := mustJoinParty(t, env, party.Code, "LateVoter")
		mustApproveGuest(t, env, adminID, party.ID, guest.ID)

		// End voting first
		mustEndVoting(t, env, adminID, party.ID)

		// Try to vote
		votes := validVotesForActs(acts)
		_, err := env.voteService.SubmitVote(ctx, "", party.ID, services.SubmitVoteRequest{
			GuestID: guest.ID,
			Votes:   votes,
		})
		assert.ErrorIs(t, err, services.ErrPartyClosed)
	})

	t.Run("update vote after close rejected", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-update-closed"

		party := mustCreateParty(t, env, adminID, "Update After Close Party")
		acts := mustGetGrandFinalActs(t, env)
		guest := mustJoinParty(t, env, party.Code, "UpdateVoter")
		mustApproveGuest(t, env, adminID, party.ID, guest.ID)

		// Submit vote while active
		votes := validVotesForActs(acts)
		mustSubmitVote(t, env, "", party.ID, guest.ID, votes)

		// End voting
		mustEndVoting(t, env, adminID, party.ID)

		// Try to update vote
		newVotes := buildRotatedVotes(acts, 1)
		_, err := env.voteService.UpdateVote(ctx, "", party.ID, services.SubmitVoteRequest{
			GuestID: guest.ID,
			Votes:   newVotes,
		})
		assert.ErrorIs(t, err, services.ErrPartyClosed)
	})

	t.Run("results not available while active", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-active-results"

		party := mustCreateParty(t, env, adminID, "Active Results Party")

		_, err := env.voteService.GetResults(ctx, adminID, party.ID)
		assert.ErrorIs(t, err, services.ErrVotingNotEnded)
	})

	t.Run("non-admin cannot end voting", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-end"
		nonAdminID := "not-the-admin"

		party := mustCreateParty(t, env, adminID, "End Voting Party")

		_, err := env.voteService.EndVoting(ctx, nonAdminID, party.ID)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
	})

	t.Run("non-admin cannot approve guest", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-approve"
		nonAdminID := "not-the-admin"

		party := mustCreateParty(t, env, adminID, "Approve Party")
		guest := mustJoinParty(t, env, party.Code, "WaitingGuest")

		err := env.guestService.ApproveGuest(ctx, nonAdminID, party.ID, guest.ID)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
	})

	t.Run("non-admin cannot list join requests", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-list-req"
		nonAdminID := "not-the-admin"

		party := mustCreateParty(t, env, adminID, "List Requests Party")
		mustJoinParty(t, env, party.Code, "RequestGuest")

		_, err := env.guestService.ListJoinRequests(ctx, nonAdminID, party.ID)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
	})

	t.Run("vote update succeeds while active", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-update"

		party := mustCreateParty(t, env, adminID, "Update Vote Party")
		acts := mustGetGrandFinalActs(t, env)
		guest := mustJoinParty(t, env, party.Code, "UpdatingVoter")
		mustApproveGuest(t, env, adminID, party.ID, guest.ID)

		// Submit initial vote
		votes := validVotesForActs(acts)
		mustSubmitVote(t, env, "", party.ID, guest.ID, votes)

		// Update vote while still active
		newVotes := buildRotatedVotes(acts, 1)
		updatedVote, err := env.voteService.UpdateVote(ctx, "", party.ID, services.SubmitVoteRequest{
			GuestID: guest.ID,
			Votes:   newVotes,
		})
		require.NoError(t, err)
		assert.Equal(t, newVotes, updatedVote.Votes)
	})

	t.Run("end voting twice returns error", func(t *testing.T) {
		env := setupTest(t)
		ctx := context.Background()
		adminID := "admin-double-end"

		party := mustCreateParty(t, env, adminID, "Double End Party")
		mustEndVoting(t, env, adminID, party.ID)

		// Try to end again
		_, err := env.voteService.EndVoting(ctx, adminID, party.ID)
		assert.ErrorIs(t, err, services.ErrPartyClosed)
	})
}
