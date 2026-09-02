package session

import (
	simplefixgo "gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go"
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/fix"
)

// MessageStorage is an interface providing a basic method for storing messages awaiting to be sent.
type MessageStorage interface {
	Save(storageID fix.StorageID, msg simplefixgo.SendingMessage, msgSeqNum int) error
	Messages(storageID fix.StorageID, msgSeqNumFrom, msgSeqNumTo int) ([]simplefixgo.SendingMessage, error)
}

// CounterStorage is a counter storage.
type CounterStorage interface {
	GetNextSeqNum(storageID fix.StorageID) (int, error)
	GetCurrSeqNum(storageID fix.StorageID) (int, error)
	ResetSeqNum(storageID fix.StorageID) error
	SetSeqNum(storageID fix.StorageID, seqNum int) error
}
