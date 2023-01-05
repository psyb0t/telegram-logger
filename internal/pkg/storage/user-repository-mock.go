package storage

import (
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
	"github.com/stretchr/testify/mock"
)

// UserRepositoryReaderMock is a mock implementation of UserRepositoryReader.
type UserRepositoryReaderMock struct {
	mock.Mock
}

// Get retrieves a user by ID.
func (r *UserRepositoryReaderMock) Get(id string) (types.User, error) {
	args := r.Called(id)
	return args.Get(0).(types.User), args.Error(1)
}

// GetAll retrieves all users from the database.
func (r *UserRepositoryReaderMock) GetAll() ([]types.User, error) {
	args := r.Called()
	return args.Get(0).([]types.User), args.Error(1)
}

// GetByTelegramChatID retrieves a user by its Telegram chat ID.
func (r *UserRepositoryReaderMock) GetByTelegramChatID(chatID int64) (types.User, error) {
	args := r.Called(chatID)
	return args.Get(0).(types.User), args.Error(1)
}

// UserRepositoryWriterMock is a mock implementation of UserRepositoryWriter.
type UserRepositoryWriterMock struct {
	mock.Mock
}

// Create stores a new user in the database.
func (r *UserRepositoryWriterMock) Create(user types.User) error {
	args := r.Called(user)
	return args.Error(0)
}

// Delete removes a user from the database by ID.
func (r *UserRepositoryWriterMock) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}

// DeleteAllByTelegramChatID removes all users from the
// database matching the given Telegram chat ID.
func (r *UserRepositoryWriterMock) DeleteAllByTelegramChatID(chatID int64) error {
	args := r.Called(chatID)
	return args.Error(0)
}
