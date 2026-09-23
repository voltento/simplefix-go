package simplefixgo

import "errors"

// Errors returned by message storage lookups.
var (
	ErrNotEnoughMessages = errors.New("not enough messages in the storage")
	ErrInvalidBoundaries = errors.New("invalid boundaries")
	ErrInvalidSequence   = errors.New("unexpected sequence index")

	// ErrSendTimeout is returned when the outbound channel stays full for the
	// handler's send deadline. The handler context is canceled alongside it,
	// so the connection tears down instead of parking every sender.
	ErrSendTimeout = errors.New("send deadline exceeded")
)
