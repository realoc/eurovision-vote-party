package persistence

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sipgate/eurovision-vote-party/server/models"
)

// UserDAO defines the persistence operations for users.
type UserDAO interface {
	Upsert(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
}

// FirestoreUserDAO is the Firestore implementation of UserDAO.
type FirestoreUserDAO struct {
	client *firestore.Client
}

// NewFirestoreUserDAO creates a new FirestoreUserDAO.
func NewFirestoreUserDAO(client *firestore.Client) *FirestoreUserDAO {
	return &FirestoreUserDAO{client: client}
}

const usersCollection = "users"

// Upsert creates or updates a user in Firestore.
func (d *FirestoreUserDAO) Upsert(ctx context.Context, user *models.User) error {
	_, err := d.client.Collection(usersCollection).Doc(user.ID).Set(ctx, user)
	return err
}

// GetByID retrieves a user by its ID.
// Returns ErrNotFound if the user does not exist.
func (d *FirestoreUserDAO) GetByID(ctx context.Context, id string) (*models.User, error) {
	doc, err := d.client.Collection(usersCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
