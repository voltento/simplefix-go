package fix

import (
	"sync"

	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/fix/buffer"
)

// MessageByteConverter is a message byte converter.
type MessageByteConverter struct {
	pool sync.Pool
}

// NewMessageByteConverter ...
func NewMessageByteConverter(bufferSize int) *MessageByteConverter {
	b := &MessageByteConverter{
		pool: sync.Pool{
			New: func() interface{} {
				return buffer.NewMessageByteBuffers(bufferSize)
			},
		},
	}
	return b
}

// ConvertableMessage is a convertable message.
type ConvertableMessage interface {
	ToBytesBuffered(buffers *buffer.MessageByteBuffers) ([]byte, error)
}

// ConvertToBytes converts msg to its wire byte representation.
func (m *MessageByteConverter) ConvertToBytes(message ConvertableMessage) ([]byte, error) {
	buffers := m.pool.Get().(*buffer.MessageByteBuffers)
	buffers.Reset()
	defer m.pool.Put(buffers)

	return message.ToBytesBuffered(buffers)
}
