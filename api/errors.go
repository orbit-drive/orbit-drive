package api

import "errors"

var (
	// ErrInValidPath is returned when the given path does not exist in the system
	// or when uploading an empty file to ipfs.
	ErrInvalidPath = errors.New("File path is invalid.")
)
