package persistence

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sipgate/eurovision-vote-party/server/models"
)

// GuestDAO defines the persistence operations for guests.
type GuestDAO interface {
	Create(ctx context.Context, guest *models.Guest) error
	GetByID(ctx context.Context, id string) (*models.Guest, error)
	ListByPartyID(ctx context.Context, partyID string) ([]*models.Guest, error)
	ListByPartyIDAndStatus(ctx context.Context, partyID string, status models.GuestStatus) ([]*models.Guest, error)
	UpdateStatus(ctx context.Context, id string, status models.GuestStatus) error
	Delete(ctx context.Context, id string) error
	ExistsByPartyAndUsername(ctx context.Context, partyID, username string) (bool, error)
}

// FirestoreGuestDAO is the Firestore implementation of GuestDAO.
type FirestoreGuestDAO struct {
	client *firestore.Client
}

// NewFirestoreGuestDAO creates a new FirestoreGuestDAO.
func NewFirestoreGuestDAO(client *firestore.Client) *FirestoreGuestDAO {
	return &FirestoreGuestDAO{client: client}
}

const guestsCollection = "guests"

// Create stores a new guest in Firestore.
func (d *FirestoreGuestDAO) Create(ctx context.Context, guest *models.Guest) error {
	_, err := d.client.Collection(guestsCollection).Doc(guest.ID).Set(ctx, guest)
	return err
}

// GetByID retrieves a guest by its ID.
// Returns ErrNotFound if the guest does not exist.
func (d *FirestoreGuestDAO) GetByID(ctx context.Context, id string) (*models.Guest, error) {
	doc, err := d.client.Collection(guestsCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var guest models.Guest
	if err := doc.DataTo(&guest); err != nil {
		return nil, err
	}

	return &guest, nil
}

// ListByPartyID retrieves all guests for a given party.
// Returns an empty slice if no guests are found.
func (d *FirestoreGuestDAO) ListByPartyID(ctx context.Context, partyID string) ([]*models.Guest, error) {
	iter := d.client.Collection(guestsCollection).Where("partyId", "==", partyID).Documents(ctx)
	defer iter.Stop()

	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}

	guests := make([]*models.Guest, 0, len(docs))
	for _, doc := range docs {
		var guest models.Guest
		if err := doc.DataTo(&guest); err != nil {
			return nil, err
		}
		guests = append(guests, &guest)
	}

	return guests, nil
}

// ListByPartyIDAndStatus retrieves all guests for a given party filtered by status.
// Returns an empty slice if no guests are found.
func (d *FirestoreGuestDAO) ListByPartyIDAndStatus(ctx context.Context, partyID string, guestStatus models.GuestStatus) ([]*models.Guest, error) {
	iter := d.client.Collection(guestsCollection).Where("partyId", "==", partyID).Where("status", "==", string(guestStatus)).Documents(ctx)
	defer iter.Stop()

	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}

	guests := make([]*models.Guest, 0, len(docs))
	for _, doc := range docs {
		var guest models.Guest
		if err := doc.DataTo(&guest); err != nil {
			return nil, err
		}
		guests = append(guests, &guest)
	}

	return guests, nil
}

// UpdateStatus updates the status field of an existing guest.
// Returns ErrNotFound if the guest does not exist.
func (d *FirestoreGuestDAO) UpdateStatus(ctx context.Context, id string, guestStatus models.GuestStatus) error {
	_, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = d.client.Collection(guestsCollection).Doc(id).Update(ctx, []firestore.Update{
		{Path: "status", Value: string(guestStatus)},
	})
	return err
}

// Delete removes a guest from Firestore.
// Returns ErrNotFound if the guest does not exist.
func (d *FirestoreGuestDAO) Delete(ctx context.Context, id string) error {
	_, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = d.client.Collection(guestsCollection).Doc(id).Delete(ctx)
	return err
}

// ExistsByPartyAndUsername checks whether a guest with the given partyID and username exists.
func (d *FirestoreGuestDAO) ExistsByPartyAndUsername(ctx context.Context, partyID, username string) (bool, error) {
	iter := d.client.Collection(guestsCollection).Where("partyId", "==", partyID).Where("username", "==", username).Limit(1).Documents(ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err != nil {
		// iterator.Done means no documents found
		return false, nil
	}

	return true, nil
}
