package messages

import "gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/fix"

// TrailerBuilder is an interface providing functionality to a builder of Trailer messages.
type TrailerBuilder interface {
	New() TrailerBuilder

	AsComponent() *fix.Component
}
