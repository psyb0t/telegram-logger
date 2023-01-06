package badgerdb

import (
	"encoding/json"

	"github.com/psyb0t/telegram-logger/internal/pkg/storage"
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

// userRepositoryReader is a struct that implements the
// storage.UserRepositoryReader interface using a badgerDB instance.
type userRepositoryReader struct {
	// db is a pointer to the underlying badgerDB instance.
	db *badgerDB
}

// newUserRepositoryReader creates and returns
// a new userRepositoryReader instance.
func newUserRepositoryReader(db *badgerDB) storage.UserRepositoryReader {
	return userRepositoryReader{db: db}
}

// Get retrieves a user by ID.
func (r userRepositoryReader) Get(id string) (types.User, error) {
	user := types.User{}

	if id == "" {
		return user, storage.ErrEmptyID
	}

	val, err := r.db.get(getUserKey(id))
	if err != nil {
		return user, err
	}

	// Unmarshal the user data into the user struct.
	if err := json.Unmarshal(val, &user); err != nil {
		return user, err
	}

	return user, nil
}

// GetAll retrieves all users from the database.
func (r userRepositoryReader) GetAll() ([]types.User, error) {
	users := []types.User{}

	// Get all user data from the db
	vals, err := r.db.getAllByPrefix([]byte(prefixUserKey))
	if err != nil {
		return nil, err
	}

	// Go through all of the returned values
	for _, val := range vals {
		// Unmarshal the user data into a user struct.
		var user types.User
		if err := json.Unmarshal(val, &user); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetByTelegramChatID retrieves the first user found by its Telegram chat ID.
func (r userRepositoryReader) GetByTelegramChatID(chatID int64) (types.User, error) {
	var user types.User
	if chatID == 0 {
		return user, storage.ErrEmptyTelegramChatID
	}

	// define filter function which unmarshals the value and checks if
	// the Telegram chat ID matches the provided one
	filterFn := func(key, val []byte) bool {
		if err := json.Unmarshal(val, &user); err != nil {
			return false
		}

		if user.TelegramChatID == chatID {
			return true
		}

		return false
	}

	// do get
	_, err := r.db.getByPrefixAndFilterFunc([]byte(prefixUserKey), filterFn, 1)
	if err != nil {
		return user, err
	}

	return user, nil
}
