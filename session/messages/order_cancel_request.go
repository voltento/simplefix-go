package messages

// OrderCancelRequest is a order cancel request.
type OrderCancelRequest interface {
	New() OrderCancelRequestBuilder
	Build() OrderCancelRequestBuilder
}

// OrderCancelRequestBuilder is an interface providing functionality to a builder of auto-generated OrderCancelRequest messages.
type OrderCancelRequestBuilder interface {
	OrderCancelRequest
	PipelineBuilder
}
