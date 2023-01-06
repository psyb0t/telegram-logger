package storage

// Storage is an interface for interacting with a storage system.
type Storage interface {
	// Open establishes a connection to the storage system using the given DSN.
	Open(dsn string) error

	// Close closes the connection to the storage system.
	Close() error

	// Ping checks if the database is reachable and responding.
	// It returns nil if the database is reachable and responding, or an error otherwise.
	Ping() error

	// GetUserRepositoryReader returns a repository for reading user data from the database
	GetUserRepositoryReader() UserRepositoryReader

	// GetUserRepositoryWriter returns a repository for writing user data from the database
	GetUserRepositoryWriter() UserRepositoryWriter
}
