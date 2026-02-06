package persistence

import "errors"

var (
	ErrNotFound   = errors.New("entity not found")
	ErrCodeExists = errors.New("party code already exists")
)
