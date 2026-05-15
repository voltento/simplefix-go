package session

import (
	"gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/session/messages"
)

type Unmarshaller interface {
	Unmarshal(msg messages.Builder, d []byte) error
}
