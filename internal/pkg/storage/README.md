# Storage

This is a Go package for interacting with a storage system

## Usage

```go
import "github.com/psyb0t/telegram-logger/internal/pkg/storage"

func main() {
	// Open a connection to the storage system
	storage, err := storage.Open("storage-dsn")
	if err != nil {
		// Handle error
	}
	defer storage.Close()

	// Check if the storage system is reachable and responding
	err = storage.Ping()
	if err != nil {
		// Handle error
	}

	// Get a reader and writer for user data
	userReader := storage.GetUserRepositoryReader()
	userWriter := storage.GetUserRepositoryWriter()

	// Use the reader and writer to read and write user data
	// ...
}
```

## Users

The `UserRepositoryReader` interface provides the following methods for reading user data:

- `Get(id string) (types.User, error)`: Retrieves a user by ID.
- `GetAll() ([]types.User, error)`: Retrieves all users from the database.
- `GetByTelegramChatID(chatID int64) (types.User, error)`: Retrieves a user by its Telegram chat ID.

The `UserRepositoryWriter` interface provides the following methods for writing user data:

- `Create(user types.User) error`: Stores a new user in the database.
- `Delete(id string) error`: Removes a user from the database by ID.
- `DeleteAllByTelegramChatID(chatID int64) error`: Removes all users from the database matching the given Telegram chat ID.

## Errors

The following errors can be returned by the repository interfaces:

- `ErrEmptyID`: Returned when an ID is empty.
- `ErrEmptyTelegramChatID`: Returned when an Telegram chat ID is empty.
- `ErrNotFound`: Returned when a user is not found.
