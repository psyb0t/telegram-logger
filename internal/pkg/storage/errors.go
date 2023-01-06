package storage

import "errors"

var (
	// ErrEmptyID is returned when an ID is empty.
	ErrEmptyID = errors.New("empty ID")

	// ErrEmptyTelegramChatID is returned when an Telegram chat ID is empty.
	ErrEmptyTelegramChatID = errors.New("empty Telegram chat ID")

	// ErrNotFound is returned when a user is not found.
	ErrNotFound = errors.New("not found")
)
