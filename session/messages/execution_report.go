package messages

// ExecutionReport is a execution report.
type ExecutionReport interface {
	New() ExecutionReportBuilder
	Build() ExecutionReportBuilder
}

// ExecutionReportBuilder is an interface providing functionality to a builder of auto-generated ExecutionReport messages.
type ExecutionReportBuilder interface {
	ExecutionReport
	PipelineBuilder
}
