package persistence

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sipgate/eurovision-vote-party/server/models"
)

// PartyDAO defines the persistence operations for parties.
type PartyDAO interface {
	Create(ctx context.Context, party *models.Party) error
	GetByID(ctx context.Context, id string) (*models.Party, error)
	GetByCode(ctx context.Context, code string) (*models.Party, error)
	ListByAdminID(ctx context.Context, adminID string) ([]*models.Party, error)
	Delete(ctx context.Context, id string) error
	CodeExists(ctx context.Context, code string) (bool, error)
	UpdateStatus(ctx context.Context, id string, status models.PartyStatus) error
}

// FirestorePartyDAO is the Firestore implementation of PartyDAO.
type FirestorePartyDAO struct {
	client *firestore.Client
}

// NewFirestorePartyDAO creates a new FirestorePartyDAO.
func NewFirestorePartyDAO(client *firestore.Client) *FirestorePartyDAO {
	return &FirestorePartyDAO{client: client}
}

const partiesCollection = "parties"

// Create stores a new party in Firestore.
// Returns ErrCodeExists if a party with the same code already exists.
func (d *FirestorePartyDAO) Create(ctx context.Context, party *models.Party) error {
	exists, err := d.CodeExists(ctx, party.Code)
	if err != nil {
		return err
	}
	if exists {
		return ErrCodeExists
	}

	_, err = d.client.Collection(partiesCollection).Doc(party.ID).Set(ctx, party)
	return err
}

// GetByID retrieves a party by its ID.
// Returns ErrNotFound if the party does not exist.
func (d *FirestorePartyDAO) GetByID(ctx context.Context, id string) (*models.Party, error) {
	doc, err := d.client.Collection(partiesCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var party models.Party
	if err := doc.DataTo(&party); err != nil {
		return nil, err
	}

	return &party, nil
}

// GetByCode retrieves a party by its unique code.
// Returns ErrNotFound if no party with the given code exists.
func (d *FirestorePartyDAO) GetByCode(ctx context.Context, code string) (*models.Party, error) {
	iter := d.client.Collection(partiesCollection).Where("code", "==", code).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		return nil, ErrNotFound
	}

	var party models.Party
	if err := doc.DataTo(&party); err != nil {
		return nil, err
	}

	return &party, nil
}

// ListByAdminID retrieves all parties created by a given admin.
// Returns an empty slice if no parties are found.
func (d *FirestorePartyDAO) ListByAdminID(ctx context.Context, adminID string) ([]*models.Party, error) {
	iter := d.client.Collection(partiesCollection).Where("adminId", "==", adminID).Documents(ctx)
	defer iter.Stop()

	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}

	parties := make([]*models.Party, 0, len(docs))
	for _, doc := range docs {
		var party models.Party
		if err := doc.DataTo(&party); err != nil {
			return nil, err
		}
		parties = append(parties, &party)
	}

	return parties, nil
}

// Delete removes a party from Firestore.
// Returns ErrNotFound if the party does not exist.
func (d *FirestorePartyDAO) Delete(ctx context.Context, id string) error {
	_, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = d.client.Collection(partiesCollection).Doc(id).Delete(ctx)
	return err
}

// UpdateStatus updates the status of a party.
// Returns ErrNotFound if the party does not exist.
func (d *FirestorePartyDAO) UpdateStatus(ctx context.Context, id string, status models.PartyStatus) error {
	_, err := d.GetByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = d.client.Collection(partiesCollection).Doc(id).Set(ctx, map[string]interface{}{
		"status": status,
	}, firestore.MergeAll)
	return err
}

// CodeExists checks whether a party with the given code exists.
func (d *FirestorePartyDAO) CodeExists(ctx context.Context, code string) (bool, error) {
	iter := d.client.Collection(partiesCollection).Where("code", "==", code).Limit(1).Documents(ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err != nil {
		// iterator.Done means no documents found
		return false, nil
	}

	return true, nil
}
