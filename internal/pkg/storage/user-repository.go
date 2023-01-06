package storage

import "github.com/psyb0t/telegram-logger/internal/pkg/types"

// UserRepositoryReader is an interface for reading
// user data stored in the database.
type UserRepositoryReader interface {
	// Get retrieves a user by ID.
	Get(id string) (types.User, error)

	// GetAll retrieves all users from the database.
	GetAll() ([]types.User, error)

	// GetByTelegramChatID retrieves a user by its Telegram chat ID.
	GetByTelegramChatID(chatID int64) (types.User, error)
}

// UserRepositoryWriter is an interface for writing
// user data stored in the database.
type UserRepositoryWriter interface {
	// Create stores a new user in the database.
	Create(user types.User) error

	// Delete removes a user from the database by ID.
	Delete(id string) error

	// DeleteAllByTelegramChatID removes all users from the
	// database matching the given Telegram chat ID.
	DeleteAllByTelegramChatID(chatID int64) error
}
