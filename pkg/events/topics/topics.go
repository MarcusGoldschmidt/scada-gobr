package topics

const (
	DataSeriesInserter    = "dataseriesinserter:"
	RuntimeManagerUpdated = "runtimemanagerupdated"

	// Append queue name
	JobStatusUpdated = "jobstatusupdated:"
)

type Getter[T any] interface {
	Get() T
}
