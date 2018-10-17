package fs

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	// ErrNotFound is returned when the record being fetch does not exist.
	ErrNotFound = leveldb.ErrNotFound

	// ErrInValidPath is returned when the given path does not exist in the system
	// or when uploading an empty file to ipfs.
	ErrInValidPath = errors.New("File: path is invalid.")
)
