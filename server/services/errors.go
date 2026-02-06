package services

import "errors"

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrNotFound          = errors.New("party not found")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrInvalidEventType  = errors.New("invalid event type")
)
