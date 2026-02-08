package query

import (
	"context"
	"database/sql"
	"time"

	"github.com/vibesql/vibe/internal/postgres"
)

var (
	QueryTimeout = 5 * time.Second
)

type ExecutionResult struct {
	Rows          []map[string]interface{}
	RowCount      int
	ExecutionTime time.Duration
}

type Executor struct {
	db *sql.DB
}

func NewExecutor(db *sql.DB) *Executor {
	return &Executor{db: db}
}

func (e *Executor) Execute(sql string) (*ExecutionResult, error) {
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	rows, err := e.db.QueryContext(ctx, sql)
	if err != nil {
		vibeErr := postgres.TranslateError(err)
		return nil, vibeErr
	}
	defer rows.Close()

	result, err := parseRows(rows)
	if err != nil {
		return nil, err
	}

	executionTime := time.Since(startTime)

	return &ExecutionResult{
		Rows:          result,
		RowCount:      len(result),
		ExecutionTime: executionTime,
	}, nil
}

func parseRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, postgres.TranslateError(err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		if err := CheckRowLimit(len(results)); err != nil {
			return nil, err
		}

		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, postgres.TranslateError(err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.TranslateError(err)
	}

	return results, nil
}
