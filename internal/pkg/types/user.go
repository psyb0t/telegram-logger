package types

// User represents a user in the system.
type User struct {
	// ID is the unique identifier for the user which acts like a token
	ID string
	// TelegramChatID is the telegram chat ID of the user
	TelegramChatID int64
}
