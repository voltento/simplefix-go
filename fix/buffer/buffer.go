// Package buffer manages reusable byte buffers for FIX message assembly.
package buffer

import "bytes"

// MessageByteBuffers is a message byte buffers.
type MessageByteBuffers struct {
	msgBuffer    *bytes.Buffer
	bodyBuffer   *bytes.Buffer
	typeBuffer   *bytes.Buffer
	headerBuffer *bytes.Buffer
}

// NewMessageByteBuffers ...
func NewMessageByteBuffers(size int) *MessageByteBuffers {
	return &MessageByteBuffers{
		typeBuffer:   bytes.NewBuffer(make([]byte, 0, 5)), // 35=AA
		msgBuffer:    bytes.NewBuffer(make([]byte, 0, size)),
		bodyBuffer:   bytes.NewBuffer(make([]byte, 0, size)),
		headerBuffer: bytes.NewBuffer(make([]byte, 0, size)),
	}
}

// Reset resets message byte buffers to its zero state.
func (m *MessageByteBuffers) Reset() {
	m.msgBuffer.Reset()
	m.bodyBuffer.Reset()
	m.typeBuffer.Reset()
	m.headerBuffer.Reset()
}

// GetMessageBuffer returns the buffer for the full message.
func (m *MessageByteBuffers) GetMessageBuffer() *bytes.Buffer {
	return m.msgBuffer
}

// GetBodyBuffer returns the buffer for the message body.
func (m *MessageByteBuffers) GetBodyBuffer() *bytes.Buffer {
	return m.bodyBuffer
}

// GetHeaderBuffer returns the buffer for the message header.
func (m *MessageByteBuffers) GetHeaderBuffer() *bytes.Buffer {
	return m.headerBuffer
}

// GetTypeBuffer returns the buffer for the message type field.
func (m *MessageByteBuffers) GetTypeBuffer() *bytes.Buffer {
	return m.typeBuffer
}
