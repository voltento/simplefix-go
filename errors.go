package simplefixgo

import "errors"

// Errors returned by message storage lookups.
var (
	ErrNotEnoughMessages = errors.New("not enough messages in the storage")
	ErrInvalidBoundaries = errors.New("invalid boundaries")
	ErrInvalidSequence   = errors.New("unexpected sequence index")
)
