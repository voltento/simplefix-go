package session

import (
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/session/messages"
)

type Unmarshaller interface {
	Unmarshal(msg messages.Builder, d []byte) error
}
