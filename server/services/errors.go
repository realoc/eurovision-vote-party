package services

import "errors"

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrNotFound          = errors.New("party not found")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrInvalidEventType  = errors.New("invalid event type")
	ErrGuestNotApproved  = errors.New("guest not approved")
	ErrPartyClosed       = errors.New("party is not active")
	ErrVoteAlreadyExists = errors.New("vote already exists")
	ErrInvalidVotes      = errors.New("invalid votes")
)
