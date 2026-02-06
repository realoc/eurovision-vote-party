package services

import (
	"context"
	"errors"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/persistence"
)

// UserDAO defines the persistence operations needed by the user service.
type UserDAO interface {
	Upsert(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
}

// UserService defines the business logic operations for users.
type UserService interface {
	UpsertProfile(ctx context.Context, userID, email, username string) (*models.User, error)
	GetProfile(ctx context.Context, userID string) (*models.User, error)
}

// userService is the default implementation.
type userService struct {
	dao UserDAO
}

// NewUserService creates a new UserService.
func NewUserService(dao UserDAO) UserService {
	return &userService{dao: dao}
}

// UpsertProfile creates or updates a user profile.
func (s *userService) UpsertProfile(ctx context.Context, userID, email, username string) (*models.User, error) {
	if err := models.ValidateUsername(username); err != nil {
		return nil, ErrInvalidUsername
	}

	user := &models.User{
		ID:       userID,
		Email:    email,
		Username: username,
	}

	if err := s.dao.Upsert(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetProfile retrieves a user profile by ID.
func (s *userService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.dao.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}
