// Package messages defines standard FIX control message builders and interfaces.
package messages

import (
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/fix"
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/fix/buffer"
)

// Builder is a builder.
type Builder interface {
	Items() fix.Items
	CalcBodyLength() int
	BodyLength() int
	BytesWithoutChecksum() []byte
	CheckSum() string
	BeginString() *fix.KeyValue
	MsgType() string
	ToBytes() ([]byte, error)
	ToBytesBuffered(buffers *buffer.MessageByteBuffers) ([]byte, error)
	BeginStringTag() string
	BodyLengthTag() string
	CheckSumTag() string
}

// PipelineBuilder is a pipeline builder.
type PipelineBuilder interface {
	HeaderBuilder() HeaderBuilder
	Builder
}
