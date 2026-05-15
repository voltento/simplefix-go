package messages

import (
	"gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/fix"
	"gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/fix/buffer"
)

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

type PipelineBuilder interface {
	HeaderBuilder() HeaderBuilder
	Builder
}
