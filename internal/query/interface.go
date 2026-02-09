package query

// QueryExecutor defines the interface for executing SQL queries
type QueryExecutor interface {
	Execute(sql string) (*ExecutionResult, error)
}

// Ensure Executor implements QueryExecutor
var _ QueryExecutor = (*Executor)(nil)
