//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

var (
	firestoreClient *firestore.Client
	actsService     services.ActsService
)

func TestMain(m *testing.M) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		log.Fatal("FIRESTORE_EMULATOR_HOST must be set for integration tests")
	}

	ctx := context.Background()
	var err error
	firestoreClient, err = firestore.NewClient(ctx, "test-project")
	if err != nil {
		log.Fatalf("failed to create Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	actsService, err = services.NewActsService("../data/acts.json")
	if err != nil {
		log.Fatalf("failed to load acts data: %v", err)
	}

	os.Exit(m.Run())
}

// testEnv holds all services and DAOs for a single test.
type testEnv struct {
	partyService services.PartyService
	guestService services.GuestService
	voteService  services.VoteService
	userService  services.UserService
	actsService  services.ActsService
}

// setupTest creates a fresh testEnv with real DAOs and services wired to the
// Firestore emulator. It registers cleanup that deletes all documents from the
// relevant collections after the test completes.
func setupTest(t *testing.T) *testEnv {
	t.Helper()

	partyDAO := persistence.NewFirestorePartyDAO(firestoreClient)
	guestDAO := persistence.NewFirestoreGuestDAO(firestoreClient)
	voteDAO := persistence.NewFirestoreVoteDAO(firestoreClient)
	userDAO := persistence.NewFirestoreUserDAO(firestoreClient)

	partyService := services.NewPartyService(partyDAO)
	guestService := services.NewGuestService(guestDAO, partyDAO)
	voteService := services.NewVoteService(voteDAO, partyDAO, guestDAO, actsService)
	userService := services.NewUserService(userDAO)

	t.Cleanup(func() {
		ctx := context.Background()
		for _, col := range []string{"parties", "guests", "votes", "users"} {
			cleanupCollection(t, ctx, col)
		}
	})

	return &testEnv{
		partyService: partyService,
		guestService: guestService,
		voteService:  voteService,
		userService:  userService,
		actsService:  actsService,
	}
}

// cleanupCollection deletes all documents from a Firestore collection.
func cleanupCollection(t *testing.T, ctx context.Context, collection string) {
	t.Helper()
	docs, err := firestoreClient.Collection(collection).Documents(ctx).GetAll()
	if err != nil {
		return
	}
	for _, doc := range docs {
		if _, err := doc.Ref.Delete(ctx); err != nil {
			t.Logf("warning: failed to delete doc %s/%s: %v", collection, doc.Ref.ID, err)
		}
	}
}

// validVotesForActs builds a valid 10-point Eurovision vote map from the first
// 10 acts in the provided slice. Points are assigned in ValidPointValues order
// (12, 10, 8, 7, 6, 5, 4, 3, 2, 1).
func validVotesForActs(acts []models.Act) map[int]string {
	votes := make(map[int]string, len(models.ValidPointValues))
	for i, points := range models.ValidPointValues {
		votes[points] = acts[i].ID
	}
	return votes
}

// buildRotatedVotes shifts the act assignment by offset positions, producing a
// different vote distribution while still assigning exactly 10 valid acts.
func buildRotatedVotes(acts []models.Act, offset int) map[int]string {
	votes := make(map[int]string, len(models.ValidPointValues))
	for i, points := range models.ValidPointValues {
		idx := (i + offset) % len(acts)
		votes[points] = acts[idx].ID
	}
	return votes
}

// mustCreateParty is a test helper that creates a party and fails the test on error.
func mustCreateParty(t *testing.T, env *testEnv, adminID, name string) *models.Party {
	t.Helper()
	ctx := context.Background()
	party, err := env.partyService.CreateParty(ctx, adminID, services.CreatePartyRequest{
		Name:      name,
		EventType: models.EventGrandFinal,
	})
	if err != nil {
		t.Fatalf("failed to create party: %v", err)
	}
	return party
}

// mustJoinParty is a test helper that joins a guest to a party and fails the test on error.
func mustJoinParty(t *testing.T, env *testEnv, code, username string) *models.Guest {
	t.Helper()
	ctx := context.Background()
	guest, err := env.guestService.JoinParty(ctx, code, username)
	if err != nil {
		t.Fatalf("failed to join party: %v", err)
	}
	return guest
}

// mustApproveGuest is a test helper that approves a guest and fails the test on error.
func mustApproveGuest(t *testing.T, env *testEnv, adminID, partyID, guestID string) {
	t.Helper()
	ctx := context.Background()
	if err := env.guestService.ApproveGuest(ctx, adminID, partyID, guestID); err != nil {
		t.Fatalf("failed to approve guest: %v", err)
	}
}

// mustGetGrandFinalActs returns at least 10 grand final acts, failing the test otherwise.
func mustGetGrandFinalActs(t *testing.T, env *testEnv) []models.Act {
	t.Helper()
	acts, err := env.actsService.ListActs(string(models.EventGrandFinal))
	if err != nil {
		t.Fatalf("failed to list acts: %v", err)
	}
	if len(acts) < len(models.ValidPointValues) {
		t.Fatalf("need at least %d acts, got %d", len(models.ValidPointValues), len(acts))
	}
	return acts
}

// mustSubmitVote submits a vote and fails the test on error.
func mustSubmitVote(t *testing.T, env *testEnv, adminID, partyID, guestID string, votes map[int]string) *models.Vote {
	t.Helper()
	ctx := context.Background()
	vote, err := env.voteService.SubmitVote(ctx, adminID, partyID, services.SubmitVoteRequest{
		GuestID: guestID,
		Votes:   votes,
	})
	if err != nil {
		t.Fatalf("failed to submit vote: %v", err)
	}
	return vote
}

// mustEndVoting ends voting for a party and fails the test on error.
func mustEndVoting(t *testing.T, env *testEnv, adminID, partyID string) *models.Party {
	t.Helper()
	ctx := context.Background()
	party, err := env.voteService.EndVoting(ctx, adminID, partyID)
	if err != nil {
		t.Fatalf("failed to end voting: %v", err)
	}
	return party
}

// actIDAtPoints finds the act ID that received the given points in a vote map.
func actIDAtPoints(votes map[int]string, points int) string {
	return votes[points]
}

// totalPointsForAct sums the points an act received across multiple vote maps.
func totalPointsForAct(actID string, allVotes ...map[int]string) int {
	total := 0
	for _, votes := range allVotes {
		for points, id := range votes {
			if id == actID {
				total += points
			}
		}
	}
	return total
}

// uniqueActIDs collects all distinct act IDs across multiple vote maps.
func uniqueActIDs(allVotes ...map[int]string) map[string]bool {
	seen := make(map[string]bool)
	for _, votes := range allVotes {
		for _, id := range votes {
			seen[id] = true
		}
	}
	return seen
}

// formatVoteSummary returns a human-readable summary of a vote map for debugging.
func formatVoteSummary(votes map[int]string) string {
	s := ""
	for _, points := range models.ValidPointValues {
		if actID, ok := votes[points]; ok {
			s += fmt.Sprintf("  %2d pts -> %s\n", points, actID)
		}
	}
	return s
}
