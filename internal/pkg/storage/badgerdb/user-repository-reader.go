package badgerdb

import (
	"encoding/json"

	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

// UserRepositoryReader is an interface for reading
// user data stored in the database.
type UserRepositoryReader interface {
	// Get retrieves a user by ID.
	Get(id string) (types.User, error)

	// GetAll retrieves all users from the database.
	GetAll() ([]types.User, error)
}

// userRepositoryReader is a struct that implements the
// UserRepositoryReader interface using a badgerDB instance.
type userRepositoryReader struct {
	// db is a pointer to the underlying badgerDB instance.
	db *badgerDB
}

// NewUserRepositoryReader creates and returns
// a new userRepositoryReader instance.
func NewUserRepositoryReader(db *badgerDB) UserRepositoryReader {
	return userRepositoryReader{db: db}
}

// Get retrieves a user by ID.
func (r userRepositoryReader) Get(id string) (types.User, error) {
	user := types.User{}

	if id == "" {
		return user, ErrEmptyID
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
