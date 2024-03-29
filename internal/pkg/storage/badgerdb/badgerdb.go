// Package badgerdb provides a BadgerDB-backed implementation of the storage.Storage interface.
package badgerdb

import (
	"context"
	"errors"
	"sync"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/psyb0t/telegram-logger/internal/pkg/storage"
)

const prefixUserKey = "user-"

// filterFunc is a function that accepts a key and its value as parameters
// used for filtering results retrieved from the database
//
// returns a boolean value. true if filter passed, false if not.
type filterFunc func(key, val []byte) bool

// badgerDB is a struct that implements the Storage interface using a BadgerDB database.
type badgerDB struct {
	ctx            context.Context //nolint:containedctx
	cancelFunc     context.CancelFunc
	db             *badger.DB
	wg             sync.WaitGroup
	userRepository struct {
		reader storage.UserRepositoryReader
		writer storage.UserRepositoryWriter
	}
}

// New creates and returns a new badgerDB instance.
func New(parentCtx context.Context) (storage.Storage, error) {
	db := &badgerDB{}
	db.ctx, db.cancelFunc = context.WithCancel(parentCtx)

	db.userRepository.reader = newUserRepositoryReader(db)
	db.userRepository.writer = newUserRepositoryWriter(db)

	return db, nil
}

// Open opens a connection to a BadgerDB database.
// The provided DSN (Data Source Name) specifies the location of the database.
func (db *badgerDB) Open(dsn string) error {
	var err error

	// Open a connection to the database.
	opts := badger.DefaultOptions(dsn)
	opts.Logger = nil

	db.db, err = badger.Open(opts)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the connection to the BadgerDB database.
func (db *badgerDB) Close() error {
	db.wg.Wait()
	db.cancelFunc()

	return db.db.Close()
}

// Ping checks if the database is reachable and responding.
// badgerDB does not have a method for this so it will just return nil.
func (db *badgerDB) Ping() error {
	return nil
}

// GetUserRepositoryReader returns a repository for reading user data from the database.
func (db *badgerDB) GetUserRepositoryReader() storage.UserRepositoryReader {
	return db.userRepository.reader
}

// GetUserRepositoryWriter returns a repository for writing user data from the database.
func (db *badgerDB) GetUserRepositoryWriter() storage.UserRepositoryWriter {
	return db.userRepository.writer
}

// get retrieves a value by key.
func (db *badgerDB) get(key []byte) ([]byte, error) {
	db.wg.Add(1)
	defer db.wg.Done()

	// Start a new transaction.
	tx := db.db.NewTransaction(false)
	defer tx.Discard()

	// Look up the data in the database using the provided key.
	item, err := tx.Get(key)
	if err != nil {
		// If the key is not found, return an error.
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, storage.ErrNotFound
		}

		return nil, err
	}

	// Retrieve the data from the item.
	val, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	return val, err
}

// getAllByPrefix retrieves a slice of values for keys prefixed with prefix.
func (db *badgerDB) getAllByPrefix(prefix []byte) ([][]byte, error) {
	db.wg.Add(1)
	defer db.wg.Done()

	result := [][]byte{}

	// Start a new transaction.
	tx := db.db.NewTransaction(false)
	defer tx.Discard()

	// Create a new iterator to iterate over all prefixed keys in the database.
	opts := badger.DefaultIteratorOptions
	opts.Prefix = prefix

	it := tx.NewIterator(opts)
	defer it.Close()

	// Iterate over all keys in the database.
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		// Retrieve the data for the current item.
		item := it.Item()
		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		// Append the value to the result slice.
		result = append(result, val)
	}

	return result, nil
}

// getByPrefixAndFilterFunc retrieves the first n number of keys and their values from the db
// based on the given key prefix and whose value also passes the filterFunc
// if n < 0 then all matches get returned
//
// the results are returned as [][][]byte
//
// eg:
//
//	results := [][][]byte{
//			[][]byte{[]byte("key1"), []byte("val1")},
//			[][]byte{[]byte("key2"), []byte("val2")},
//		}
func (db *badgerDB) getByPrefixAndFilterFunc(prefix []byte, filterFn filterFunc, n int) ([][][]byte, error) {
	db.wg.Add(1)
	defer db.wg.Done()

	results := [][][]byte{}

	// if desired num results is 0 then just return
	if n == 0 {
		return results, nil
	}

	// Start a new transaction.
	tx := db.db.NewTransaction(false)
	defer tx.Discard()

	// Create a new iterator to iterate over all prefixed keys in the database.
	opts := badger.DefaultIteratorOptions
	opts.Prefix = prefix

	it := tx.NewIterator(opts)
	defer it.Close()

	// Iterate over all keys in the database.
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		// Retrieve the data for the current item.
		item := it.Item()

		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		// Check if the value passes the filter func
		key := item.KeyCopy(nil)
		if filterFn(key, val) {
			// add the result to the list of results
			results = append(results, [][]byte{item.KeyCopy(nil), val})
			// if num results is the desired one return
			if len(results) == n {
				return results, nil
			}
		}
	}

	return results, nil
}

// create stores a value for the given key.
func (db *badgerDB) create(key []byte, val []byte) error {
	db.wg.Add(1)
	defer db.wg.Done()

	// Start a new transaction.
	tx := db.db.NewTransaction(true)
	defer tx.Discard()

	// Set the data in the database using the provided key.
	if err := tx.Set(key, val); err != nil {
		return err
	}

	// Commit the transaction.
	return tx.Commit()
}

// delete removes data from the database by the given key.
func (db *badgerDB) delete(key []byte) error {
	db.wg.Add(1)
	defer db.wg.Done()

	// Start a new transaction.
	tx := db.db.NewTransaction(true)
	defer tx.Discard()

	// Delete the data from the database using the provided key.
	if err := tx.Delete(key); err != nil {
		// If the key is not found, return the ErrNotFound error.
		if errors.Is(err, badger.ErrKeyNotFound) {
			return storage.ErrNotFound
		}

		return err
	}

	// Commit the transaction.
	return tx.Commit()
}

/*
// isCtxDone checks if the context is done and returns a boolean
func (db *badgerDB) isCtxDone() bool {
	select {
	default:
		return false
	case <-db.ctx.Done():
		return true
	}
}
*/
