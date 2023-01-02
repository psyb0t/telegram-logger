package v1

import "errors"

var (
	// ErrUnableToOpenDatabaseConnection is returned when there is an error when opening a database connection.
	ErrUnableToOpenDatabaseConnection = errors.New("error when opening database connection")
	// ErrUnsupportedStorageType is returned when an unsupported storage type is used.
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
	// ErrUnauthorizedToUseTelegramBotCommand is returned when the trying to use the telegram bot command is not a superuser
	ErrUnauthorizedToUseTelegramBotCommand = errors.New("unauthorized to use command")
)
