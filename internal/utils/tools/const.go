package tools

const (
	DefaultNotAssignInt = -1
	StatusEnable        = 1
)

const (
	ErrorFormatFind         = "mongosql.Find [ Collection: %s, condition: %v, Error: %w ]"
	ErrorFormatCursorDecode = "cursor.Decode [ Collection: %s, condition: %v, Error: %w ]"

	ErrorFormatAggregate             = "mongo.Collection Aggregate [ Collection: %s, pipeline: %v, Error: %w ]"
	ErrorFormatAggregateCursorDecode = "mongo.Collection Aggregate cursor.Decode [ Collection: %s, pipeline: %v, Error: %w ]"
)
