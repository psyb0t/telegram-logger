package v1

import "errors"

var (
	// ErrUnableToOpenDatabaseConnection is returned when there is an error when opening a database connection.
	ErrUnableToOpenDatabaseConnection = errors.New("error when opening database connection")
	// ErrUnsupportedStorageType is returned when an unsupported storage type is used.
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
)
