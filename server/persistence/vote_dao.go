package persistence

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/sipgate/eurovision-vote-party/server/models"
)

// VoteDAO defines the persistence operations for votes.
type VoteDAO interface {
	Create(ctx context.Context, vote *models.Vote) error
	GetByGuestAndParty(ctx context.Context, guestID, partyID string) (*models.Vote, error)
	Update(ctx context.Context, vote *models.Vote) error
	ListByPartyID(ctx context.Context, partyID string) ([]*models.Vote, error)
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

// firestoreVote is the Firestore-compatible representation of a Vote.
// Firestore requires map keys to be strings, so the int-keyed Votes map
// is stored as map[string]string.
type firestoreVote struct {
	ID        string            `firestore:"id"`
	GuestID   string            `firestore:"guestId"`
	PartyID   string            `firestore:"partyId"`
	Votes     map[string]string `firestore:"votes"`
	CreatedAt time.Time         `firestore:"createdAt"`
}

func toFirestoreVote(v *models.Vote) *firestoreVote {
	votes := make(map[string]string, len(v.Votes))
	for points, actID := range v.Votes {
		votes[strconv.Itoa(points)] = actID
	}
	return &firestoreVote{
		ID:        v.ID,
		GuestID:   v.GuestID,
		PartyID:   v.PartyID,
		Votes:     votes,
		CreatedAt: v.CreatedAt,
	}
}

func fromFirestoreVote(fv *firestoreVote) (*models.Vote, error) {
	votes := make(map[int]string, len(fv.Votes))
	for k, actID := range fv.Votes {
		points, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("invalid vote key %q: %w", k, err)
		}
		votes[points] = actID
	}
	return &models.Vote{
		ID:        fv.ID,
		GuestID:   fv.GuestID,
		PartyID:   fv.PartyID,
		Votes:     votes,
		CreatedAt: fv.CreatedAt,
	}, nil
}

// Create stores a new vote in Firestore.
func (d *FirestoreVoteDAO) Create(ctx context.Context, vote *models.Vote) error {
	_, err := d.client.Collection(votesCollection).Doc(vote.ID).Set(ctx, toFirestoreVote(vote))
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

	var fv firestoreVote
	if err := doc.DataTo(&fv); err != nil {
		return nil, err
	}

	return fromFirestoreVote(&fv)
}

// Update overwrites an existing vote in Firestore.
func (d *FirestoreVoteDAO) Update(ctx context.Context, vote *models.Vote) error {
	_, err := d.client.Collection(votesCollection).Doc(vote.ID).Set(ctx, toFirestoreVote(vote))
	return err
}

// ListByPartyID retrieves all votes for a given party.
func (d *FirestoreVoteDAO) ListByPartyID(ctx context.Context, partyID string) ([]*models.Vote, error) {
	iter := d.client.Collection(votesCollection).Where("partyId", "==", partyID).Documents(ctx)
	defer iter.Stop()

	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}

	votes := make([]*models.Vote, 0, len(docs))
	for _, doc := range docs {
		var fv firestoreVote
		if err := doc.DataTo(&fv); err != nil {
			return nil, err
		}
		vote, err := fromFirestoreVote(&fv)
		if err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}

	return votes, nil
}
