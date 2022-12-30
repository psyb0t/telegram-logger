package badgerdb

import "errors"

var (
	// ErrEmptyID is returned when an ID is empty.
	ErrEmptyID = errors.New("empty ID")

	// ErrNotFound is returned when a user is not found.
	ErrNotFound = errors.New("not found")
)
