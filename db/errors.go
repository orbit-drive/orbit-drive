package db

import "github.com/syndtr/goleveldb/leveldb"

var (
	// ErrNotFound is returned when the record being fetch does not exist.
	ErrNotFound = leveldb.ErrNotFound
)
