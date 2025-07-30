package storage

import "errors"

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record already exists")
	ErrInvalid  = errors.New("invalid data")
)
