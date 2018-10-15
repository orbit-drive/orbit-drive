package fs

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	// DB Errors
	ErrNotFound = leveldb.ErrNotFound

	ErrInvalidKey = errors.New("Db: key invalid.")
)
