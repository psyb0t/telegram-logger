package badgerdb

import (
	"encoding/json"

	"github.com/psyb0t/telegram-logger/internal/pkg/storage"
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

// userRepositoryWriter is a struct that implements the
// storage.UserRepositoryWriter interface using a badgerDB instance.
type userRepositoryWriter struct {
	// db is a pointer to the underlying badgerDB instance.
	db *badgerDB
}

// newUserRepositoryWriter creates and returns
// a new userRepositoryWriter instance.
func newUserRepositoryWriter(db *badgerDB) storage.UserRepositoryWriter {
	return userRepositoryWriter{db: db}
}

// Create stores a new user in the database.
//
// user is the user to be stored. It must have a non-empty ID field.
func (r userRepositoryWriter) Create(user types.User) error {
	if user.ID == "" {
		return storage.ErrEmptyID
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
		return storage.ErrEmptyID
	}

	// Delete the user data from the database using the prefixed provided ID as the key.
	return r.db.delete(getUserKey(id))
}

// DeleteByTelegramChatID removes a user from the database by Telegram chat ID.
func (r userRepositoryWriter) DeleteByTelegramChatID(chatID int64) error {
	if chatID == 0 {
		return storage.ErrEmptyTelegramChatID
	}

	// define filter function which unmarshals the value and checks if
	// the Telegram chat ID matches the provided one
	filterFn := func(key, val []byte) bool {
		var user types.User
		if err := json.Unmarshal(val, &user); err != nil {
			return false
		}

		if user.TelegramChatID == chatID {
			return true
		}

		return false
	}

	// find key to delete
	key, _, err := r.db.getByPrefixAndFilterFunc([]byte(prefixUserKey), filterFn)
	if err != nil {
		return err
	}

	return r.db.delete(key)
}
