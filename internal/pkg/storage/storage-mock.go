package storage

import "github.com/stretchr/testify/mock"

// Mock is a mock implementation of Storage.
type Mock struct {
	mock.Mock

	userRepositoryReader UserRepositoryReader
	userRepositoryWriter UserRepositoryWriter
}

// NewMock returns a new instance of Mock.
func NewMock() *Mock {
	return &Mock{
		userRepositoryReader: &UserRepositoryReaderMock{},
		userRepositoryWriter: &UserRepositoryWriterMock{},
	}
}

// Open establishes a connection to the storage system using the given DSN.
func (db *Mock) Open(dsn string) error {
	args := db.Called(dsn)
	return args.Error(0)
}

// Close closes the connection to the storage system.
func (db *Mock) Close() error {
	args := db.Called()
	return args.Error(0)
}

// Ping checks if the database is reachable and responding.
// It returns nil if the database is reachable and responding, or an error otherwise.
func (db *Mock) Ping() error {
	args := db.Called()
	return args.Error(0)
}

// GetUserRepositoryReader returns a repository for reading user data from the database
func (db *Mock) GetUserRepositoryReader() UserRepositoryReader {
	return db.userRepositoryReader
}

// GetUserRepositoryWriter returns a repository for writing user data from the database
func (db *Mock) GetUserRepositoryWriter() UserRepositoryWriter {
	return db.userRepositoryWriter
}
