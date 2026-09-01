package messages

import (
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/fix/buffer"
)

// MockMessage is a mock message.
type MockMessage struct {
	Type string
	Data []byte
	Err  error
}

// NewMockMessage ...
func NewMockMessage(tp string, data []byte, err error) *MockMessage {
	return &MockMessage{Type: tp, Data: data, Err: err}
}

// HeaderBuilder header builder.
func (m MockMessage) HeaderBuilder() HeaderBuilder {
	return nil
}

// MsgType msg type.
func (m MockMessage) MsgType() string {
	return m.Type
}

// ToBytes returns the byte representation of mock message.
func (m MockMessage) ToBytes() ([]byte, error) {
	return m.Data, m.Err
}

// ToBytesBuffered writes mock message's bytes into buf and returns the buffered result.
func (m MockMessage) ToBytesBuffered(_ *buffer.MessageByteBuffers) ([]byte, error) {
	return m.Data, m.Err
}
