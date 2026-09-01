package messages

// Reject is a reject.
type Reject interface {
	New() RejectBuilder
	Build() RejectBuilder
	RefTagID() int
	SetFieldRefTagID(int) RejectBuilder
	RefSeqNum() int
	SetFieldRefSeqNum(int) RejectBuilder
	SessionRejectReason() string
	SetFieldSessionRejectReason(string) RejectBuilder
}

// RejectBuilder is an interface providing functionality to a builder of auto-generated Reject messages.
type RejectBuilder interface {
	Reject
	PipelineBuilder
}
