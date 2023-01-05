# BadgerDB

Package `badgerdb` provides a BadgerDB-backed implementation of the `storage.Storage` interface.

## Usage

```go
import "github.com/psyb0t/telegram-logger/internal/pkg/storage/badgerdb"

// Create a new badgerDB instance.
db, err := badgerdb.New(context.Background())
if err != nil {
	// handle error
}

// Open a connection to the database.
dsn := "/path/to/badgerdb/directory"
if err := db.Open(dsn); err != nil {
	// handle error
}

// Use the UserRepositoryReader and UserRepositoryWriter to interact with the user data in the database.
reader := db.GetUserRepositoryReader()
writer := db.GetUserRepositoryWriter()

// When you're done, close the connection to the database.
if err := db.Close(); err != nil {
	// handle error
}
```

## TODO

- add `update` method and check if key exists on `create`
- write tests
