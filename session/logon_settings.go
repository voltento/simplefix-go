// Package session implements the FIX session state machine: logon/logout, heartbeats, and sequence management.
package session

import "time"

// LogonSettings holds the fields used to build a session's logon request.
// TODO: constructor for acceptor and initiator
type LogonSettings struct {
	TargetCompID    string
	SenderCompID    string
	HeartBtInt      int
	EncryptMethod   string
	Password        string
	Username        string
	LogonTimeout    time.Duration // todo
	HeartBtLimits   *IntLimits
	CloseTimeout    time.Duration
	ResetSeqNumFlag bool
}
