package badger

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

const (
	// Key prefixes. Each prefix value should be unique.
	projectPrefix   = 0x00
	reqLogPrefix    = 0x01
	resLogPrefix    = 0x02
	senderReqPrefix = 0x03

	// Request log indices.
	reqLogProjectIDIndex = 0x00

	// Sender request indices.
	senderReqProjectIDIndex = 0x00
)

// Database is used to store and retrieve data from an underlying Badger database.
type Database struct {
	badger *badger.DB
}

// OpenDatabase opens a new Badger database.
func OpenDatabase(opts badger.Options) (*Database, error) {
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("badger: failed to open database: %w", err)
	}

	return &Database{badger: db}, nil
}

// Close closes the underlying Badger database.
func (db *Database) Close() error {
	return db.badger.Close()
}

// DatabaseFromBadgerDB returns a Database with `db` set as the underlying
// Badger database.
func DatabaseFromBadgerDB(db *badger.DB) *Database {
	return &Database{badger: db}
}

func entryKey(prefix, index byte, value []byte) []byte {
	// Key consists of: | prefix (byte) | index (byte) | value
	key := make([]byte, 2+len(value))
	key[0] = prefix
	key[1] = index
	copy(key[2:len(value)+2], value)

	return key
}
