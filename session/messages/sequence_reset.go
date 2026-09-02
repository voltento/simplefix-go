package messages

// SequenceReset is a sequence reset.
type SequenceReset interface {
	New() SequenceResetBuilder
	Build() SequenceResetBuilder
	NewSeqNo() int
	SetFieldNewSeqNo(int) SequenceResetBuilder
	GapFillFlag() bool
	SetFieldGapFillFlag(bool) SequenceResetBuilder
}

// SequenceResetBuilder is an interface providing functionality to a builder of auto-generated SequenceReset messages.
type SequenceResetBuilder interface {
	SequenceReset
	PipelineBuilder
}
