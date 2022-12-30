package badgerdb

import (
	"encoding/json"

	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

// UserRepositoryWriter is an interface for writing
// user data stored in the database.
type UserRepositoryWriter interface {
	// Create stores a new user in the database.
	Create(user types.User) error

	// Delete removes a user from the database by ID.
	Delete(id string) error
}

// userRepositoryWriter is a struct that implements the
// UserRepositoryWriter interface using a badgerDB instance.
type userRepositoryWriter struct {
	// db is a pointer to the underlying badgerDB instance.
	db *badgerDB
}

// NewUserRepositoryWriter creates and returns
// a new userRepositoryWriter instance.
func NewUserRepositoryWriter(db *badgerDB) UserRepositoryWriter {
	return userRepositoryWriter{db: db}
}

// Create stores a new user in the database.
//
// user is the user to be stored. It must have a non-empty ID field.
func (r userRepositoryWriter) Create(user types.User) error {
	if user.ID == "" {
		return ErrEmptyID
	}

	// Convert the user struct to a byte slice.
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Create the user
	return r.db.create(getUserKey(user.ID), val)
}

// Delete removes a user from the database by ID.
func (r userRepositoryWriter) Delete(id string) error {
	if id == "" {
		return ErrEmptyID
	}

	// Delete the user data from the database using the prefixed provided ID as the key.
	return r.db.delete(getUserKey(id))
}
