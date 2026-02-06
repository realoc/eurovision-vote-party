package persistence

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/sipgate/eurovision-vote-party/server/models"
)

// VoteDAO defines the persistence operations for votes.
type VoteDAO interface {
	Create(ctx context.Context, vote *models.Vote) error
	GetByGuestAndParty(ctx context.Context, guestID, partyID string) (*models.Vote, error)
	Update(ctx context.Context, vote *models.Vote) error
}

// FirestoreVoteDAO is the Firestore implementation of VoteDAO.
type FirestoreVoteDAO struct {
	client *firestore.Client
}

// NewFirestoreVoteDAO creates a new FirestoreVoteDAO.
func NewFirestoreVoteDAO(client *firestore.Client) *FirestoreVoteDAO {
	return &FirestoreVoteDAO{client: client}
}

const votesCollection = "votes"

// Create stores a new vote in Firestore.
func (d *FirestoreVoteDAO) Create(ctx context.Context, vote *models.Vote) error {
	_, err := d.client.Collection(votesCollection).Doc(vote.ID).Set(ctx, vote)
	return err
}

// GetByGuestAndParty retrieves a vote by guest ID and party ID.
// Returns ErrNotFound if no vote exists for the given guest and party.
func (d *FirestoreVoteDAO) GetByGuestAndParty(ctx context.Context, guestID, partyID string) (*models.Vote, error) {
	iter := d.client.Collection(votesCollection).Where("guestId", "==", guestID).Where("partyId", "==", partyID).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		return nil, ErrNotFound
	}

	var vote models.Vote
	if err := doc.DataTo(&vote); err != nil {
		return nil, err
	}

	return &vote, nil
}

// Update overwrites an existing vote in Firestore.
func (d *FirestoreVoteDAO) Update(ctx context.Context, vote *models.Vote) error {
	_, err := d.client.Collection(votesCollection).Doc(vote.ID).Set(ctx, vote)
	return err
}
