package query

// QueryExecutor defines the interface for executing SQL queries
type QueryExecutor interface {
	Execute(sql string, params ...interface{}) (*ExecutionResult, error)
}

// Ensure Executor implements QueryExecutor
var _ QueryExecutor = (*Executor)(nil)
