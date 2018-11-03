package fs

import (
	"errors"
)

var (
	// ErrNotADir is returned when accessing the links of a file type vnode.
	ErrNotADir = errors.New("File does not have any links.")

	// ErrNotAFile is returned when accessing the source of a dir type vnode.
	ErrNotAFile = errors.New("Directory does not have a source.")

	// ErrVNodeNotFound is returned when a vnode is missing a Link.
	ErrVNodeNotFound = errors.New("VNode does not exist.")

	// ErrInValidPath is returned when the given path does not exist in the system.
	ErrInvalidPath = errors.New("File path is invalid.")
)
